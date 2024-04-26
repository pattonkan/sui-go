package client_test

import (
	"testing"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/howjmay/go-sui-sdk/client"
	"github.com/stretchr/testify/require"
)

func TestFaucetRequestFund_Devnet(t *testing.T) {
	res, err := client.FaucetRequestFund(account.TEST_ADDRESS.String(), client.DevnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestFaucetRequestFund_Testnet(t *testing.T) {
	res, err := client.FaucetRequestFund(account.TEST_ADDRESS.String(), client.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestFaucetRequestFund_Localnet(t *testing.T) {
	res, err := client.FaucetRequestFund(account.TEST_ADDRESS.String(), client.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}
