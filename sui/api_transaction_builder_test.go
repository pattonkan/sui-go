package sui_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"

	"github.com/stretchr/testify/require"
)

func TestBatchTransaction(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.NewSignerByIndex(sui_signer.TEST_SEED, sui_signer.KeySchemeFlagDefault, 0)

	getSuiCoins, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{Owner: signer.Address})
	require.NoError(t, err)

	amount := 10
	txnBytes, err := client.BatchTransaction(context.Background(), &sui.BatchTransactionRequest{
		Signer: signer.Address,
		TxnParams: []models.RPCTransactionRequestParams{
			{
				MoveCallRequestParams: &models.MoveCallParams{
					PackageObjectId: sui_types.SuiPackageIdSuiFramework,
					Module:          sui_types.Identifier("pay"),
					Function:        sui_types.Identifier("split"),
					TypeArguments:   []models.SuiTypeTag{"0x2::sui::SUI"},
					Arguments: []models.SuiJsonValue{
						models.SuiJsonValue(getSuiCoins.Data[2].CoinObjectID.String()),
						models.SuiJsonValue(fmt.Sprintf("%d", amount)),
					},
				},
			},
			{
				TransferObjectRequestParams: &models.TransferObjectParams{
					Recipient: signer.Address,
					ObjectId:  getSuiCoins.Data[3].CoinObjectID,
				},
			},
		},
		GasBudget: models.NewBigInt(sui.DefaultGasBudget),
	})
	require.NoError(t, err)

	res, err := client.DryRunTransaction(context.Background(), txnBytes.TxBytes)
	require.NoError(t, err)
	require.True(t, res.Effects.Data.IsSuccess())
}

func TestMergeCoins(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer,
		Limit: 10,
	})
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := client.MergeCoins(
		context.Background(),
		&sui.MergeCoinsRequest{
			Signer:      signer,
			PrimaryCoin: coin1.CoinObjectID,
			CoinToMerge: coin2.CoinObjectID,
			Gas:         coin3.CoinObjectID,
			GasBudget:   coin3.Balance,
		},
	)
	require.NoError(t, err)

	res, err := client.DryRunTransaction(context.Background(), txn.TxBytes)
	require.NoError(t, err)
	require.True(t, res.Effects.Data.IsSuccess())
}

func TestMoveCall(t *testing.T) {
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)

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
		&sui.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID, err := txnResponse.GetPublishedPackageID()
	require.NoError(t, err)

	// test MoveCall with byte array input
	input := []string{"haha", "gogo"}
	txnBytes, err = client.MoveCall(
		context.Background(),
		&sui.MoveCallRequest{
			Signer:    signer.Address,
			PackageID: packageID,
			Module:    "sdk_verify",
			Function:  "read_input_bytes_array",
			TypeArgs:  []string{},
			Arguments: []any{input},
			GasBudget: models.NewBigInt((sui.DefaultGasBudget)),
		},
	)
	require.NoError(t, err)
	txnResponse, err = client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	queryEventsRes, err := client.QueryEvents(
		context.Background(),
		&sui.QueryEventsRequest{
			Query: &models.EventFilter{Transaction: &txnResponse.Digest},
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	coins, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 10,
	})
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas
	totalBal := models.Coins(coins.Data).TotalBalance().Uint64()

	amount := uint64(123)
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), sui.DefaultGasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := client.Pay(
		context.Background(),
		&sui.PayRequest{
			Signer:     signer.Address,
			InputCoins: pickedCoins.CoinIds(),
			Recipients: []*sui_types.SuiAddress{recipient.Address},
			Amount:     []*models.BigInt{models.NewBigInt(amount)},
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)
	// assume the account holds more than 'limit' amount Sui token objects
	require.Len(t, coinPages.Data, 3)
	totalBal := coins.TotalBalance()

	txn, err := client.PayAllSui(
		context.Background(),
		&sui.PayAllSuiRequest{
			Signer:     signer.Address,
			Recipient:  recipient.Address,
			InputCoins: coins.ObjectIDs(),
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
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
			require.Contains(t, coins.ObjectIDVals(), change.Data.Mutated.ObjectID)
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient1 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	_, recipient2 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 2)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	sentAmounts := []uint64{123, 456, 789}
	txn, err := client.PaySui(
		context.Background(),
		&sui.PaySuiRequest{
			Signer:     signer.Address,
			InputCoins: coins.ObjectIDs(),
			Recipients: []*sui_types.SuiAddress{
				recipient1.Address,
				recipient2.Address,
				recipient2.Address,
			},
			Amount: []*models.BigInt{
				models.NewBigInt(sentAmounts[0]), // to recipient1
				models.NewBigInt(sentAmounts[1]), // to recipient2
				models.NewBigInt(sentAmounts[2]), // to recipient2
			},
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
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
			require.Contains(t, coins.ObjectIDVals(), change.Data.Mutated.ObjectID)
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)

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
		&sui.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       models.NewBigInt(sui.DefaultGasBudget * 5),
		},
	)
	require.NoError(t, err)

	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())
}

func TestSplitCoin(t *testing.T) {
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	txn, err := client.SplitCoin(
		context.Background(),
		&sui.SplitCoinRequest{
			Signer: signer.Address,
			Coin:   coins[1].CoinObjectID,
			SplitAmounts: []*models.BigInt{
				// assume coins[0] has more than the sum of the following splitAmounts
				models.NewBigInt(2222),
				models.NewBigInt(1111),
			},
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	splitShares := uint64(3)
	txn, err := client.SplitCoinEqual(
		context.Background(),
		&sui.SplitCoinEqualRequest{
			Signer:     signer.Address,
			Coin:       coins[0].CoinObjectID,
			SplitCount: models.NewBigInt(splitShares),
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	transferCoin := coinPages.Data[0]

	txn, err := client.TransferObject(
		context.Background(),
		&sui.TransferObjectRequest{
			Signer:    signer.Address,
			Recipient: recipient.Address,
			ObjectID:  transferCoin.CoinObjectID,
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	transferCoin := coinPages.Data[0]

	txn, err := client.TransferSui(
		context.Background(),
		&sui.TransferSuiRequest{
			Signer:    signer.Address,
			Recipient: recipient.Address,
			ObjectID:  transferCoin.CoinObjectID,
			Amount:    models.NewBigInt(3),
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
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
			require.Equal(t, *transferCoin.CoinObjectID, change.Data.Mutated.ObjectID)
			require.Equal(t, signer.Address, change.Data.Mutated.Owner.AddressOwner)

		} else if change.Data.Created != nil {
			require.Equal(t, recipient.Address, change.Data.Created.Owner.AddressOwner)
		}
	}

	require.Len(t, simulate.BalanceChanges, 2)
}
