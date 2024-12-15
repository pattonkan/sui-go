package movebcs_test

import (
	"context"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/movebcs"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suiclient/conn"
	"github.com/howjmay/sui-go/suisigner"
	"github.com/stretchr/testify/require"
)

func TestCoinDecode(t *testing.T) {
	client := suiclient.NewClient(conn.TestnetEndpointUrl)
	signer := suisigner.NewSigner(suisigner.TEST_SEED, suisigner.KeySchemeFlagEd25519)
	resGetCoins, err := client.GetCoins(context.TODO(), &suiclient.GetCoinsRequest{
		Owner:    signer.Address,
		CoinType: &sui.SuiCoinType,
	})
	require.NoError(t, err)
	// we should have requested a lot of coin objects from faucet. So coin in index 5, should have been used
	targetCoinRef := resGetCoins.Data[5].Ref()
	resGetObject, err := client.GetObject(context.TODO(), &suiclient.GetObjectRequest{
		ObjectId: targetCoinRef.ObjectId,
		Options: &suiclient.SuiObjectDataOptions{
			ShowBcs: true,
		},
	})
	require.NoError(t, err)
	var moveCoin movebcs.MoveCoin
	_, err = bcs.Unmarshal(resGetObject.GetMoveObjectInBcs(), &moveCoin)
	require.NoError(t, err)
	require.Equal(t, uint64(1000000000), moveCoin.Balance)
}
