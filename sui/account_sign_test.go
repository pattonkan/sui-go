package sui_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"

	"github.com/stretchr/testify/require"
)

func TestAccountSignAndSend(t *testing.T) {
	account, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)
	t.Log(account.Address)

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := account.Address
	coins, err := api.GetSuiCoinsOwnedByAddress(context.Background(), signer)
	require.NoError(t, err)
	require.Greater(t, coins.TotalBalance().Int64(), sui_types.SUI(0.01).Int64(), "insufficient balance")

	coinIDs := make([]sui_types.ObjectID, len(coins))
	for i, c := range coins {
		coinIDs[i] = c.CoinObjectID
	}
	gasBudget := models.NewSafeSuiBigInt(uint64(10000000))
	txn, err := api.PayAllSui(context.Background(), signer, signer, coinIDs, gasBudget)
	require.NoError(t, err)

	resp := executeTxn(t, api, txn.TxBytes, account)
	t.Log("txn digest: ", resp.Digest)
}
