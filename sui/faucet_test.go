package sui_test

import (
	"testing"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"

	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	res, err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.DevnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	res, err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	t.Skip("only run with local node is set up")
	res, err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("txn digest: ", res)
}
