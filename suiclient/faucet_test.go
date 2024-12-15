package suiclient_test

import (
	"testing"

	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suiclient/conn"
	"github.com/howjmay/sui-go/suisigner"

	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.DevnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.TestnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	t.Skip("only run with local node is set up")
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
}
