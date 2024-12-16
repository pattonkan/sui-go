package suiclient_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/stretchr/testify/require"
)

const (
	ComingChatValidatorAddress = "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
)

func TestRequestAddDelegation(t *testing.T) {
	client, signer := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	coins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 10,
	})
	require.NoError(t, err)

	amount := uint64(sui.UnitSui)
	pickedCoins, err := suiclient.PickupCoins(coins, new(big.Int).SetUint64(amount), 0, 0, 0)
	require.NoError(t, err)

	validatorAddress := ComingChatValidatorAddress
	validator, err := sui.AddressFromHex(validatorAddress)
	require.NoError(t, err)

	txBytes, err := suiclient.BCS_RequestAddStake(
		signer.Address,
		pickedCoins.CoinRefs(),
		sui.NewBigInt(amount),
		validator,
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Equal(t, "", simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
}

func TestRequestWithdrawDelegation(t *testing.T) {
	client := suiclient.NewClient(conn.TestnetEndpointUrl)

	signer, err := sui.AddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	stakes, err := client.GetStakes(context.Background(), signer)
	require.NoError(t, err)
	require.True(t, len(stakes) > 0)
	require.True(t, len(stakes[0].Stakes) > 0)

	coins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer,
		Limit: 10,
	})
	require.NoError(t, err)
	pickedCoins, err := suiclient.PickupCoins(coins, new(big.Int), suiclient.DefaultGasBudget, 0, 0)
	require.NoError(t, err)

	detail, err := client.GetObject(context.Background(), &suiclient.GetObjectRequest{
		ObjectId: &stakes[0].Stakes[0].Data.StakedSuiId,
	})
	require.NoError(t, err)
	txBytes, err := suiclient.BCS_RequestWithdrawStake(signer, detail.Data.Ref(), pickedCoins.CoinRefs(), suiclient.DefaultGasBudget, 1000)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Equal(t, "", simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
}
