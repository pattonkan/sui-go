package suiclient_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
)

func TestDevInspectTransactionBlock(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: limit,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	ptb.PayAllSui(sender.Address)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		coins.CoinRefs(),
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx.V1.Kind)
	require.NoError(t, err)

	resp, err := client.DevInspectTransactionBlock(
		context.Background(),
		&suiclient.DevInspectTransactionBlockRequest{
			SenderAddress: sender.Address,
			TxKindBytes:   txBytes,
		},
	)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
}

func TestDryRunTransaction(t *testing.T) {
	api := suiclient.NewClient(conn.TestnetEndpointUrl)
	signer := suisigner.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer,
		Limit: 10,
	})
	require.NoError(t, err)
	pickedCoins, err := suiclient.PickupCoins(coins, big.NewInt(100), suiclient.DefaultGasBudget, 0, 0)
	require.NoError(t, err)
	pt, err := api.PayAllSui(
		context.Background(),
		&suiclient.PayAllSuiRequest{
			Signer:     signer,
			Recipient:  signer,
			InputCoins: pickedCoins.CoinIds(),
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)

	resp, err := api.DryRunTransaction(context.Background(), pt.TxBytes)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Empty(t, resp.Effects.Data.V1.Status.Error)
}
