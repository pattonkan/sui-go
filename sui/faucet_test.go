package sui_test

import (
	"testing"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui/conn"
	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	res, err := sui.RequestFundFromFaucet(account.TEST_ADDRESS.String(), conn.DevnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	res, err := sui.RequestFundFromFaucet(account.TEST_ADDRESS.String(), conn.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	t.Skip("only run with local node is set up")
	res, err := sui.RequestFundFromFaucet(account.TEST_ADDRESS.String(), conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}
