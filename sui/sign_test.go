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
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC, sui_signer.KeySchemeFlagDefault)
	require.NoError(t, err)
	t.Log(signer.Address)

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	coins, err := api.GetSuiCoinsOwnedByAddress(context.Background(), signer.Address)
	require.NoError(t, err)
	require.Greater(t, coins.TotalBalance().Int64(), int64(sui_types.UnitSui)/100, "insufficient balance")

	coinIDs := make([]*sui_types.ObjectID, len(coins))
	for i, c := range coins {
		coinIDs[i] = c.CoinObjectID
	}
	gasBudget := new(models.BigInt).SetUint64(sui.DefaultGasBudget)
	txn, err := api.PayAllSui(context.Background(), signer.Address, signer.Address, coinIDs, gasBudget)
	require.NoError(t, err)

	resp := executeTxn(t, api, txn.TxBytes, signer)
	t.Log("txn digest: ", resp.Digest)
}
