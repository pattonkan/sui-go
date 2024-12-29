package movebcs_test

import (
	"context"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/movebcs"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/stretchr/testify/require"
)

func TestCoinDecode(t *testing.T) {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 10)
	resGetCoins, err := client.GetCoins(context.TODO(), &suiclient.GetCoinsRequest{
		Owner:    signer.Address,
		CoinType: &sui.SuiCoinType,
	})
	require.NoError(t, err)
	// we should have requested a lot of coin objects from faucet. So coin in index 5, should have been used
	targetCoinRef := resGetCoins.Data[len(resGetCoins.Data)-1].Ref()
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
	require.Equal(t, uint64(200000000000), moveCoin.Balance)
}
