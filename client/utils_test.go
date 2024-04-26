package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/howjmay/go-sui-sdk/client"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/stretchr/testify/require"
)

var (
	TestAddress, _ = sui_types.NewAddressFromHex("0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5")
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

	balance, err := c.GetBalance(context.Background(), *account.TEST_ADDRESS, types.SUI_COIN_TYPE)
	require.NoError(t, err)
	if balance.TotalBalance.BigInt().Uint64() < sui_types.SUI(0.3).Uint64() {
		_, err = client.FaucetRequestFund(account.TEST_ADDRESS.String(), client.DevNetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func ChainClient(t *testing.T) *client.Client {
	suiEnv := os.Getenv("SUI_NETWORK")
	switch suiEnv {
	case "testnet":
		return TestnetClient(t)
	case "devnet":
		return DevnetClient(t)
	case "":
		fallthrough
	default:
		return TestnetClient(t)
	}
}

func M1Account(t *testing.T) *account.Account {
	a, err := account.NewAccountWithMnemonic(account.TEST_MNEMONIC)
	require.NoError(t, err)
	return a
}

func M1Address(t *testing.T) *sui_types.SuiAddress {
	return account.TEST_ADDRESS
}

func Signer(t *testing.T) *account.Account {
	return M1Account(t)
}

func AddressFromStrMust(str string) *sui_types.SuiAddress {
	s, _ := sui_types.NewAddressFromHex(str)
	return s
}
