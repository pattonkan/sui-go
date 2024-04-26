package client_test

import (
	"testing"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/howjmay/go-sui-sdk/client"
	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	res, err := client.RequestFundFromFaucet(account.TEST_ADDRESS.String(), client.DevnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	res, err := client.RequestFundFromFaucet(account.TEST_ADDRESS.String(), client.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	t.Skip("only run with local node is set up")
	res, err := client.RequestFundFromFaucet(account.TEST_ADDRESS.String(), client.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}
