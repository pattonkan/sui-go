package sui_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
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
	t.Log("TestBatchTransaction TODO")
	// api := sui.NewSuiClient(conn.DevnetEndpointUrl)

	// txnBytes, err := api.BatchTransaction(context.Background(), signer, *coin1, *coin2, nil, 10000)
	// require.NoError(t, err)
	// dryRunTxn(t, api, txnBytes, M1Account(t))
}

func TestMergeCoins(t *testing.T) {
	t.Skip("FIXME create an account has at least two coin objects on chain")
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := api.MergeCoins(
		context.Background(), signer,
		coin1.CoinObjectID, coin2.CoinObjectID,
		coin3.CoinObjectID, coin3.Balance,
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestMoveCall(t *testing.T) {
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.TestnetEndpointUrl, sui_signer.TEST_MNEMONIC)

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
		signer.Address,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
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

	packageID := txnResponse.GetPublishedPackageID()

	// test MoveCall with byte array input
	input := []string{"haha", "gogo"}
	txnBytes, err = client.MoveCall(
		context.Background(),
		signer.Address,
		packageID,
		"sdk_verify",
		"read_input_bytes_array",
		[]string{},
		[]any{input},
		nil,
		models.NewSafeSuiBigInt(uint64(sui.DefaultGasBudget)),
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
		&models.EventFilter{
			Transaction: &txnResponse.Digest,
		},
		nil,
		nil,
		false,
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
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipient := sui_signer.NewRandomSigner(sui_signer.KeySchemeFlagDefault)
	coinType := models.SuiCoinType
	coins, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, 10)
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas
	totalBal := models.Coins(coins.Data).TotalBalance().Uint64()

	amount := uint64(123)
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), sui.DefaultGasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := client.Pay(
		context.Background(),
		signer.Address,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient.Address},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amount),
		},
		nil,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.TestnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipient := sui_signer.NewRandomSigner(sui_signer.KeySchemeFlagDefault)
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)
	// assume the account holds more than 'limit' amount Sui token objects
	require.Len(t, coinPages.Data, 3)
	totalBal := coins.TotalBalance()

	txn, err := client.PayAllSui(
		context.Background(),
		signer.Address,
		recipient.Address,
		coins.ObjectIDs(),
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
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
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.TestnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipients := []*sui_signer.Signer{
		sui_signer.NewRandomSigner(sui_signer.KeySchemeFlagDefault),
		sui_signer.NewRandomSigner(sui_signer.KeySchemeFlagDefault),
	}

	coinType := models.SuiCoinType
	limit := uint(4)
	coinPages, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	sentAmounts := []uint64{123, 456, 789}
	txn, err := client.PaySui(
		context.Background(),
		signer.Address,
		coins.ObjectIDs(),
		[]*sui_types.SuiAddress{
			recipients[0].Address,
			recipients[1].Address,
			recipients[1].Address,
		},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(sentAmounts[0]), // to recipients[0]
			models.NewSafeSuiBigInt(sentAmounts[1]), // to recipients[1]
			models.NewSafeSuiBigInt(sentAmounts[2]), // to recipients[1]
		},
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
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
	// 1 for recipients[0], and 2 for recipients[1]
	require.Equal(t, amountNum, createdObjNum)

	// one output balance and one input balance for recipients[0] and one input balance for recipients[1]
	require.Len(t, simulate.BalanceChanges, 3)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == signer.Address {
			require.Equal(t, coins.TotalBalance().Neg(coins.TotalBalance()), balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipients[0].Address {
			require.Equal(t, sentAmounts[0], balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipients[1].Address {
			require.Equal(t, sentAmounts[1]+sentAmounts[2], balChange.Amount)
		}
	}
}

func TestPublish(t *testing.T) {
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.TestnetEndpointUrl, sui_signer.TEST_MNEMONIC)

	// If local side has installed Sui-cli then the user can use the following func to build move contracts
	// modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin")
	// require.NoError(t, err)

	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/testcoin/contract_base64.json")
	require.NoError(t, err)

	var modules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &modules)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		signer.Address,
		modules.Modules,
		modules.Dependencies,
		nil, // 'unsafe_publish' API can automatically assign gas object
		models.NewSafeSuiBigInt(sui.DefaultGasBudget*5),
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
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.TestnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	coins, err := client.GetCoins(context.Background(), signer.Address, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)
	splitCoins := []models.SafeSuiBigInt[uint64]{models.NewSafeSuiBigInt(amount / 2)}

	txn, err := client.SplitCoin(
		context.Background(),
		signer.Address,
		pickedCoins.Coins[0].CoinObjectID,
		splitCoins,
		nil,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, client, txn.TxBytes, false)
}

func TestSplitCoinEqual(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)

	txn, err := api.SplitCoinEqual(
		context.Background(),
		signer,
		pickedCoins.Coins[0].CoinObjectID,
		models.NewSafeSuiBigInt(uint64(2)),
		nil,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestTransferObject(t *testing.T) {
	t.Skip("FIXME create an account has at least two coin objects on chain")
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 2)
	coin := coins.Data[0]

	txn, err := api.TransferObject(
		context.Background(),
		signer,
		recipient,
		coin.CoinObjectID,
		nil,
		models.NewSafeSuiBigInt(sui_types.SUI(0.01).Uint64()),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestTransferSui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.0001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), gasBudget, 1, 0)
	require.NoError(t, err)

	txn, err := api.TransferSui(
		context.Background(), signer, recipient,
		pickedCoins.Coins[0].CoinObjectID,
		models.NewSafeSuiBigInt(amount),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}
