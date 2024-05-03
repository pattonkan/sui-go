package sui_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/howjmay/go-sui-sdk/lib"
	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui/conn"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"

	"github.com/stretchr/testify/require"
)

func MainnetClient(t *testing.T) *conn.HttpClient {
	c := conn.Dial(conn.MainnetEndpointUrl)
	return c
}

func TestnetClient(t *testing.T) *conn.HttpClient {
	c := conn.Dial(conn.TestnetEndpointUrl)
	return c
}

func DevnetClient(t *testing.T) *conn.HttpClient {
	c := conn.Dial(conn.DevnetEndpointUrl)
	api := sui.NewSuiClient(c)
	balance, err := api.GetBalance(context.Background(), sui_types.TEST_ADDRESS, types.SUI_COIN_TYPE)
	require.NoError(t, err)
	if balance.TotalBalance.BigInt().Uint64() < sui_types.SUI(0.3).Uint64() {
		_, err = sui.RequestFundFromFaucet(sui_types.TEST_ADDRESS.String(), conn.DevnetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func LocalnetClient(t *testing.T) *conn.HttpClient {
	c := conn.Dial(conn.LocalnetEndpointUrl)
	return c
}

func AddressFromStrMust(str string) *sui_types.SuiAddress {
	s, _ := sui_types.NewAddressFromHex(str)
	return s
}

// @return types.DryRunTransactionBlockResponse
func dryRunTxn(
	t *testing.T,
	api *sui.ImplSuiAPI,
	txBytes lib.Base64Data,
	showJson bool,
) *types.DryRunTransactionBlockResponse {
	simulate, err := api.DryRunTransaction(context.Background(), txBytes)
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
	api *sui.ImplSuiAPI,
	txBytes lib.Base64Data,
	acc *sui_types.Account,
) *types.SuiTransactionBlockResponse {
	// First of all, make sure that there are no problems with simulated trading.
	simulate, err := api.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// sign and send
	signature, err := acc.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	require.NoError(t, err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := api.ExecuteTransactionBlock(
		context.TODO(), txBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	t.Log(resp)
	return resp
}
