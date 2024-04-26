package client_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/howjmay/go-sui-sdk/client"
	"github.com/howjmay/go-sui-sdk/move_types"
	"github.com/howjmay/go-sui-sdk/sui_types"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/howjmay/go-sui-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClient_TransferObject(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 2)
	coin := coins.Data[0]

	txn, err := cli.TransferObject(
		context.Background(), signer, recipient,
		&coin.CoinObjectID, nil, types.NewSafeSuiBigInt(sui_types.SUI(0.01).Uint64()),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_TransferSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.0001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 1, 0)
	require.NoError(t, err)

	txn, err := cli.TransferSui(
		context.Background(), signer, recipient,
		&pickedCoins.Coins[0].CoinObjectID,
		types.NewSafeSuiBigInt(amount),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_PayAllSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := cli.PayAllSui(
		context.Background(), signer, recipient,
		pickedCoins.CoinIds(),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_Pay(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	recipient := account.TEST_ADDRESS
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := cli.Pay(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		nil,
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_PaySui(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	recipient := account.TEST_ADDRESS
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := cli.PaySui(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_SplitCoin(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)
	splitCoins := []types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amount / 2)}

	txn, err := cli.SplitCoin(
		context.Background(), signer,
		&pickedCoins.Coins[0].CoinObjectID,
		splitCoins,
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, false)
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)

	txn, err := cli.SplitCoinEqual(
		context.Background(), signer,
		&pickedCoins.Coins[0].CoinObjectID,
		types.NewSafeSuiBigInt(uint64(2)),
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_MergeCoins(t *testing.T) {
	cli := DevnetClient(t)
	signer := account.TEST_ADDRESS
	coins, err := cli.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := cli.MergeCoins(
		context.Background(), signer,
		&coin1.CoinObjectID, &coin2.CoinObjectID,
		&coin3.CoinObjectID, coin3.Balance,
	)
	require.NoError(t, err)

	dryRunTxn(t, cli, txn.TxBytes, true)
}

func TestClient_Publish(t *testing.T) {
	t.Log("TestClient_Publish TODO")
	// cli := DevnetClient(t)

	// txnBytes, err := cli.Publish(context.Background(), signer, *coin1, *coin2, nil, 10000)
	// require.NoError(t, err)
	// dryRunTxn(t, cli, txnBytes, M1Account(t))
}

func TestClient_MoveCall(t *testing.T) {
	cli := TestnetClient(t)
	account, err := account.NewAccountWithMnemonic(account.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Log("signer: ", account.Address)
	digest, err := client.RequestFundFromFaucet(account.Address, client.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)

	packageID, err := move_types.NewAccountAddressHex("0x2")
	require.NoError(t, err)
	txnBytes, err := cli.MoveCall(
		context.Background(),
		account.AccountAddress(),
		packageID,
		"address",
		"length",
		[]string{},
		[]any{},
		nil,
		types.NewSafeSuiBigInt(uint64(10000000)),
	)
	require.NoError(t, err)

	signature, err := account.SignSecureWithoutEncode(txnBytes.TxBytes.Data(), sui_types.DefaultIntent())
	require.NoError(t, err)
	txnResponse, err := cli.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes.Data(), []any{signature}, &types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, types.TxnRequestTypeWaitForLocalExecution)
	require.NoError(t, err)
	t.Log(txnResponse)

	// try dry-run
	dryRunTxn(t, cli, txnBytes.TxBytes, true)
}

func TestClient_BatchTransaction(t *testing.T) {
	t.Log("TestClient_BatchTransaction TODO")
	// cli := DevnetClient(t)

	// txnBytes, err := cli.BatchTransaction(context.Background(), signer, *coin1, *coin2, nil, 10000)
	// require.NoError(t, err)
	// dryRunTxn(t, cli, txnBytes, M1Account(t))
}
