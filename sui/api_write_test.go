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

	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
)

func TestDevInspectTransactionBlock(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	ptb := sui_types.NewProgrammableTransactionBuilder()
	ptb.PayAllSui(sender.Address)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		coins.CoinRefs(),
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx.V1.Kind)
	require.NoError(t, err)

	resp, err := client.DevInspectTransactionBlock(
		context.Background(),
		sender.Address,
		txBytes,
		nil,
		nil,
	)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
}

func TestDryRunTransaction(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(100), sui.DefaultGasBudget, 0, 0)
	require.NoError(t, err)
	tx, err := api.PayAllSui(
		context.Background(),
		signer,
		signer,
		pickedCoins.CoinIds(),
		models.NewBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)

	resp, err := api.DryRunTransaction(context.Background(), tx.TxBytes)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Empty(t, resp.Effects.Data.V1.Status.Error)
}
