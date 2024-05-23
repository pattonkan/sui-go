package sui_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestDevInspectTransactionBlock(t *testing.T) {
	t.Skip("FIXME")
	// 	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	// 	signer := Address
	// 	price, err := api.GetReferenceGasPrice(context.TODO())
	// 	require.NoError(t, err)
	// 	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	// 	require.NoError(t, err)

	// 	amount := SUI(0.01).Int64()
	// 	gasBudget := SUI(0.01).Uint64()
	// 	pickedCoins, err := models.PickupCoins(coins, *big.NewInt(amount * 2), 0, false)
	// 	require.NoError(t, err)
	// 	tx, err := api.PayAllSui(context.Background(),
	// 		signer, signer,
	// 		pickedCoins.CoinIds(),
	// 		models.NewSafeSuiBigInt(gasBudget))
	// 	require.NoError(t, err)

	// resp, err := api.DevInspectTransactionBlock(context.Background(), signer, tx.TxBytes, price, nil)
	// require.NoError(t, err)
	// t.Log(resp)
}

func TestDryRunTransaction(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)
	tx, err := api.PayAllSui(
		context.Background(), signer, signer,
		pickedCoins.CoinIds(),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	resp, err := api.DryRunTransaction(context.Background(), tx.TxBytes)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Empty(t, resp.Effects.Data.V1.Status.Error)
}
