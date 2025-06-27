package suiclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"

	"github.com/stretchr/testify/require"
)

func TestBatchTransaction(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	getSuiCoins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{Owner: signer.Address})
	require.NoError(t, err)

	amount := 10
	txnBytes, err := client.BatchTransaction(context.Background(), &suiclient.BatchTransactionRequest{
		Signer: signer.Address,
		TxnParams: []suiclient.RPCTransactionRequestParams{
			{
				MoveCallRequestParams: &suiclient.MoveCallParams{
					PackageObjectId: sui.SuiPackageIdSuiFramework,
					Module:          sui.Identifier("pay"),
					Function:        sui.Identifier("split"),
					TypeArguments:   []suiclient.SuiTypeTag{"0x2::sui::SUI"},
					Arguments: []suiclient.SuiJsonValue{
						suiclient.SuiJsonValue(getSuiCoins.Data[2].CoinObjectId.String()),
						suiclient.SuiJsonValue(fmt.Sprintf("%d", amount)),
					},
				},
			},
			{
				TransferObjectRequestParams: &suiclient.TransferObjectParams{
					Recipient: signer.Address,
					ObjectId:  getSuiCoins.Data[3].CoinObjectId,
				},
			},
		},
		GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
	})
	require.NoError(t, err)

	res, err := client.DryRunTransaction(context.Background(), txnBytes.TxBytes)
	require.NoError(t, err)
	require.True(t, res.Effects.Data.IsSuccess())
}

func TestMergeCoins(t *testing.T) {
	t.Skip("TODO use localnet to test this")
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	coins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 10,
	})
	require.NoError(t, err)
	require.Equal(t, len(coins.Data), 5)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := client.MergeCoins(
		context.Background(),
		&suiclient.MergeCoinsRequest{
			Signer:      signer.Address,
			PrimaryCoin: coin1.CoinObjectId,
			CoinToMerge: coin2.CoinObjectId,
			Gas:         coin3.CoinObjectId,
			GasBudget:   coin3.Balance,
		},
	)
	require.NoError(t, err)

	res, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.True(t, res.Effects.Data.IsSuccess())
}

func TestMoveCall(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)

	// directly build (need sui toolchain)
	// modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/sdk_verify/")
	// require.NoError(t, err)
	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/sdk_verify/contract_base64.json")
	require.NoError(t, err)

	var modules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &modules)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageId, err := txnResponse.GetPublishedPackageId()
	require.NoError(t, err)

	// test MoveCall with byte array input
	input := []string{"haha", "gogo"}
	txnBytes, err = client.MoveCall(
		context.Background(),
		&suiclient.MoveCallRequest{
			Signer:    signer.Address,
			PackageId: packageId,
			Module:    "sdk_verify",
			Function:  "read_input_bytes_array",
			TypeArgs:  []string{},
			Arguments: []any{input},
			GasBudget: sui.NewBigInt((suiclient.DefaultGasBudget)),
		},
	)
	require.NoError(t, err)
	txnResponse, err = client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	queryEventsRes, err := client.QueryEvents(
		context.Background(),
		&suiclient.QueryEventsRequest{
			Query: &suiclient.EventFilter{Transaction: &txnResponse.Digest},
		},
	)
	require.NoError(t, err)
	queryEventsResMap := queryEventsRes.Data[0].ParsedJson.(map[string]interface{})
	b, err := json.Marshal(queryEventsResMap["data"])
	require.NoError(t, err)
	var res [][]byte
	err = json.Unmarshal(b, &res)
	require.NoError(t, err)

	require.Equal(t, []byte("haha"), res[0])
	require.Equal(t, []byte("gogo"), res[1])
}

func TestPay(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	coins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 10,
	})
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas
	totalBal := suiclient.Coins(coins.Data).TotalBalance().Uint64()

	amount := uint64(123)
	pickedCoins, err := suiclient.PickupCoins(coins, new(big.Int).SetUint64(amount), suiclient.DefaultGasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := client.Pay(
		context.Background(),
		&suiclient.PayRequest{
			Signer:     signer.Address,
			InputCoins: pickedCoins.CoinIds(),
			Recipients: []*sui.Address{recipient.Address},
			Amount:     []*sui.BigInt{sui.NewBigInt(amount)},
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	require.Len(t, simulate.BalanceChanges, 2)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == recipient.Address {
			require.Equal(t, amount, balChange.Amount)
		} else if balChange.Owner.AddressOwner == signer.Address {
			require.Equal(t, totalBal-amount, balChange.Amount)
		}
	}
}

func TestPayAllSui(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)
	// assume the account holds more than 'limit' amount Sui token objects
	require.Len(t, coinPages.Data, 3)
	totalBal := coins.TotalBalance()

	txn, err := client.PayAllSui(
		context.Background(),
		&suiclient.PayAllSuiRequest{
			Signer:     signer.Address,
			Recipient:  recipient.Address,
			InputCoins: coins.ObjectIds(),
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	require.Len(t, simulate.ObjectChanges, int(limit))
	delObjNum := uint(0)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, *signer.Address, change.Data.Mutated.Sender)
			require.Contains(t, coins.ObjectIdVals(), change.Data.Mutated.ObjectId)
		} else if change.Data.Deleted != nil {
			delObjNum += 1
		}
	}
	// all the input objects are merged into the first input object
	// except the first input object, all the other input objects are deleted
	require.Equal(t, limit-1, delObjNum)

	// one output balance and one input balance
	require.Len(t, simulate.BalanceChanges, 2)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == signer.Address {
			require.Equal(t, totalBal.Neg(totalBal), balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipient.Address {
			require.Equal(t, totalBal, balChange.Amount)
		}
	}
}

func TestPaySui(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient1 := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	_, recipient2 := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 2)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)

	sentAmounts := []uint64{123, 456, 789}
	txn, err := client.PaySui(
		context.Background(),
		&suiclient.PaySuiRequest{
			Signer:     signer.Address,
			InputCoins: coins.ObjectIds(),
			Recipients: []*sui.Address{
				recipient1.Address,
				recipient2.Address,
				recipient2.Address,
			},
			Amount: []*sui.BigInt{
				sui.NewBigInt(sentAmounts[0]), // to recipient1
				sui.NewBigInt(sentAmounts[1]), // to recipient2
				sui.NewBigInt(sentAmounts[2]), // to recipient2
			},
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// 3 stands for the three amounts (3 crated SUI objects) in unsafe_paySui API
	amountNum := uint(3)
	require.Len(t, simulate.ObjectChanges, int(limit)+int(amountNum))
	delObjNum := uint(0)
	createdObjNum := uint(0)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, *signer.Address, change.Data.Mutated.Sender)
			require.Contains(t, coins.ObjectIdVals(), change.Data.Mutated.ObjectId)
		} else if change.Data.Created != nil {
			createdObjNum += 1
			require.Equal(t, *signer.Address, change.Data.Created.Sender)
		} else if change.Data.Deleted != nil {
			delObjNum += 1
		}
	}

	// all the input objects are merged into the first input object
	// except the first input object, all the other input objects are deleted
	require.Equal(t, limit-1, delObjNum)
	// 1 for recipient1, and 2 for recipient2
	require.Equal(t, amountNum, createdObjNum)

	// one output balance and one input balance for recipient1 and one input balance for recipient2
	require.Len(t, simulate.BalanceChanges, 3)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == signer.Address {
			require.Equal(t, coins.TotalBalance().Neg(coins.TotalBalance()), balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipient1.Address {
			require.Equal(t, sentAmounts[0], balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipient2.Address {
			require.Equal(t, sentAmounts[1]+sentAmounts[2], balChange.Amount)
		}
	}
}

func TestPublish(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)

	// If local side has installed Sui-cli then the user can use the following func to build move contracts
	// modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin/")
	// require.NoError(t, err)

	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/testcoin/contract_base64.json")
	require.NoError(t, err)

	var modules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &modules)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(suiclient.DefaultGasBudget * 5),
		},
	)
	require.NoError(t, err)

	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())
}

func TestSplitCoin(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)

	txn, err := client.SplitCoin(
		context.Background(),
		&suiclient.SplitCoinRequest{
			Signer: signer.Address,
			Coin:   coins[1].CoinObjectId,
			SplitAmounts: []*sui.BigInt{
				// assume coins[0] has more than the sum of the following splitAmounts
				sui.NewBigInt(2222),
				sui.NewBigInt(1111),
			},
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// 2 mutated and 2 created (split coins)
	require.Len(t, simulate.ObjectChanges, 4)
	// TODO check each element ObjectChanges
	require.Len(t, simulate.BalanceChanges, 1)
	amt, _ := strconv.ParseInt(simulate.BalanceChanges[0].Amount, 10, 64)
	require.Equal(t, amt, -simulate.Effects.Data.GasFee())
}

func TestSplitCoinEqual(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)

	splitShares := uint64(3)
	txn, err := client.SplitCoinEqual(
		context.Background(),
		&suiclient.SplitCoinEqualRequest{
			Signer:     signer.Address,
			Coin:       coins[0].CoinObjectId,
			SplitCount: sui.NewBigInt(splitShares),
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// 1 mutated and 3 created (split coins)
	require.Len(t, simulate.ObjectChanges, 1+int(splitShares))
	// TODO check each element ObjectChanges
	require.Len(t, simulate.BalanceChanges, 1)
	amt, _ := strconv.ParseInt(simulate.BalanceChanges[0].Amount, 10, 64)
	require.Equal(t, amt, -simulate.Effects.Data.GasFee())
}

func TestTransferObject(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	transferCoin := coinPages.Data[0]

	txn, err := client.TransferObject(
		context.Background(),
		&suiclient.TransferObjectRequest{
			Signer:    signer.Address,
			Recipient: recipient.Address,
			ObjectId:  transferCoin.CoinObjectId,
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// one is transferred object, one is the gas object
	require.Len(t, simulate.ObjectChanges, 2)

	require.Len(t, simulate.BalanceChanges, 2)
}

func TestTransferSui(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	transferCoin := coinPages.Data[0]

	txn, err := client.TransferSui(
		context.Background(),
		&suiclient.TransferSuiRequest{
			Signer:    signer.Address,
			Recipient: recipient.Address,
			ObjectId:  transferCoin.CoinObjectId,
			Amount:    sui.NewBigInt(3),
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// one is transferred object, one is the gas object
	require.Len(t, simulate.ObjectChanges, 2)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, *transferCoin.CoinObjectId, change.Data.Mutated.ObjectId)
			require.Equal(t, signer.Address, change.Data.Mutated.Owner.AddressOwner)

		} else if change.Data.Created != nil {
			require.Equal(t, recipient.Address, change.Data.Created.Owner.AddressOwner)
		}
	}

	require.Len(t, simulate.BalanceChanges, 2)
}
