package client_test

import (
	"context"
	"testing"

	"github.com/howjmay/go-sui-sdk/account"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAccountSignAndSend(t *testing.T) {
	// ManualTest_AccountSignAndSend(t)
}

func ManualTest_AccountSignAndSend(t *testing.T) {
	unsafeMnemonic := account.TEST_MNEMONIC

	account, err := account.NewAccountWithMnemonic(unsafeMnemonic)
	require.Nil(t, err)
	t.Log(account.Address)

	cli := TestnetClient(t)
	signer := AddressFromStrMust(account.Address)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	require.Nil(t, err)
	require.Greater(t, coins.TotalBalance().Int64(), sui_types.SUI(0.01).Int64(), "insufficient balance")

	coinIds := make([]sui_types.ObjectID, len(coins))
	for i, c := range coins {
		coinIds[i] = c.CoinObjectId
	}
	gasBudget := types.NewSafeSuiBigInt(uint64(10000000))
	txn, err := cli.PayAllSui(context.Background(), *signer, *signer, coinIds, gasBudget)
	require.Nil(t, err)

	resp := executeTxn(t, cli, txn.TxBytes, account)
	t.Log("txn digest: ", resp.Digest)
}
