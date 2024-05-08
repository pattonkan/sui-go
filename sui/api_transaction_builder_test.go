package sui_test

import (
	"context"
	"fmt"
	"math/big"
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
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Log("sui_signer: ", signer.Address)
	digest, err := sui.RequestFundFromFaucet(signer.Address, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)

	packageID, err := sui_types.SuiAddressFromHex("0x2")
	require.NoError(t, err)

	txnBytes, err := api.MoveCall(
		context.Background(),
		signer.Address,
		packageID,
		"address",
		"length",
		[]string{},
		[]any{},
		nil,
		models.NewSafeSuiBigInt(uint64(10000000)),
	)
	require.NoError(t, err)

	signature, err := signer.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_signer.DefaultIntent())
	require.NoError(t, err)
	txnResponse, err := api.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes.Data(), []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	require.NoError(t, err)
	t.Log(txnResponse)

	// try dry-run
	dryRunTxn(t, api, txnBytes.TxBytes, true)
}

func TestPay(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := api.Pay(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amount),
		},
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPayAllSui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := api.PayAllSui(
		context.Background(), signer, recipient,
		pickedCoins.CoinIds(),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPaySui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := api.PaySui(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amount),
		},
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPublish(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)

	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin")
	require.NoError(t, err)

	dependencies := make([]*sui_types.ObjectID, len(modules.Dependencies))
	for i, v := range modules.Dependencies {
		dependencies[i] = sui_types.MustObjectIDFromHex(v)
	}

	coins, err := client.GetCoins(context.Background(), sui_signer.TEST_ADDRESS, nil, nil, 10)
	require.NoError(t, err)
	gasBudget := uint64(1000000)
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(100000), gasBudget, 10, 10)
	require.NoError(t, err)

	txnBytes, err := client.Publish(context.Background(), sui_signer.TEST_ADDRESS, modules.Modules, dependencies, pickedCoins.CoinIds()[0], models.NewSafeSuiBigInt(gasBudget))
	require.NoError(t, err)

	fmt.Printf("txnBytes: %#v\n", txnBytes)
}

func TestSplitCoin(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)
	splitCoins := []models.SafeSuiBigInt[uint64]{models.NewSafeSuiBigInt(amount / 2)}

	txn, err := api.SplitCoin(
		context.Background(), signer,
		pickedCoins.Coins[0].CoinObjectID,
		splitCoins,
		nil, models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, false)
}

func TestSplitCoinEqual(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)

	txn, err := api.SplitCoinEqual(
		context.Background(), signer,
		pickedCoins.Coins[0].CoinObjectID,
		models.NewSafeSuiBigInt(uint64(2)),
		nil, models.NewSafeSuiBigInt(gasBudget),
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
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 1, 0)
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
