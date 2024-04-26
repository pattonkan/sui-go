package client_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/howjmay/go-sui-sdk/client"
	"github.com/howjmay/go-sui-sdk/lib"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/stretchr/testify/require"
)

func MainnetClient(t *testing.T) *client.Client {
	c := client.Dial(client.MainnetEndpointUrl)
	return c
}

func TestnetClient(t *testing.T) *client.Client {
	c := client.Dial(client.TestnetEndpointUrl)
	return c
}

func DevnetClient(t *testing.T) *client.Client {
	c := client.Dial(client.DevnetEndpointUrl)

	balance, err := c.GetBalance(context.Background(), account.TEST_ADDRESS, types.SUI_COIN_TYPE)
	require.NoError(t, err)
	if balance.TotalBalance.BigInt().Uint64() < sui_types.SUI(0.3).Uint64() {
		_, err = client.RequestFundFromFaucet(account.TEST_ADDRESS.String(), client.DevnetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func LocalnetClient(t *testing.T) *client.Client {
	c := client.Dial(client.LocalnetEndpointUrl)
	return c
}

func AddressFromStrMust(str string) *sui_types.SuiAddress {
	s, _ := sui_types.NewAddressFromHex(str)
	return s
}

// @return types.DryRunTransactionBlockResponse
func dryRunTxn(
	t *testing.T,
	cli *client.Client,
	txBytes lib.Base64Data,
	showJson bool,
) *types.DryRunTransactionBlockResponse {
	simulate, err := cli.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Equal(t, "", simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	if showJson {
		data, err := json.Marshal(simulate)
		require.NoError(t, err)
		t.Log(string(data))
		t.Log("gasFee: ", simulate.Effects.Data.GasFee())
	}
	return simulate
}

func executeTxn(
	t *testing.T,
	cli *client.Client,
	txBytes lib.Base64Data,
	acc *account.Account,
) *types.SuiTransactionBlockResponse {
	// First of all, make sure that there are no problems with simulated trading.
	simulate, err := cli.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// sign and send
	signature, err := acc.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	require.NoError(t, err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := cli.ExecuteTransactionBlock(
		context.TODO(), txBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	t.Log(resp)
	return resp
}
