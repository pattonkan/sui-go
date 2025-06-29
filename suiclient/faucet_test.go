package suiclient_test

import (
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	t.Skip("FIXME we have too many requests to the devnet faucet, so we skip this test for now")
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.DevnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	t.Skip("FIXME we have too many requests to the testnet faucet, so we skip this test for now")
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.TestnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	err := suiclient.RequestFundFromFaucet(suisigner.TEST_ADDRESS, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
}
