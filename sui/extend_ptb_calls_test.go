package sui_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"

	"github.com/stretchr/testify/require"
)

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)

	options := models.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
	}
	coinType := fmt.Sprintf("0x2::coin::Coin<%v>", models.SuiCoinType)
	filterObject, err := api.BatchGetObjectsOwnedByAddress(context.TODO(), sui_signer.TEST_ADDRESS, &options, coinType)
	require.NoError(t, err)
	t.Log(filterObject)
}

func getCoins(t *testing.T, api *sui.ImplSuiAPI, sender *sui_types.SuiAddress, needCoinObjNum int) []*models.Coin {
	coins, err := api.GetCoins(context.Background(), sender, nil, nil, uint(needCoinObjNum))
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= needCoinObjNum)
	return coins.Data
}
