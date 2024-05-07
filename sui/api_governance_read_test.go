package sui_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetLatestSuiSystemState(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	state, err := api.GetLatestSuiSystemState(context.Background())
	require.NoError(t, err)
	t.Logf("system state: %v", state)
}

func TestGetReferenceGasPrice(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	gasPrice, err := api.GetReferenceGasPrice(context.Background())
	require.NoError(t, err)
	t.Logf("current gas price: %v", gasPrice)
}

// TODO Add TestGetStakes

func TestGetStakesByIds(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	owner, err := sui_types.SuiAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	stakes, err := api.GetStakes(context.Background(), owner)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(stakes), 1)

	stake1 := stakes[0].Stakes[0].Data
	stakeId := stake1.StakedSuiId
	stakesFromId, err := api.GetStakesByIds(context.Background(), []sui_types.ObjectID{stakeId})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(stakesFromId), 1)

	queriedStake := stakesFromId[0].Stakes[0].Data
	require.Equal(t, stake1, queriedStake)
	t.Log(stakesFromId)
}

func TestGetValidatorsApy(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	apys, err := api.GetValidatorsApy(context.Background())
	require.NoError(t, err)
	t.Logf("current epoch %v", apys.Epoch)
	apyMap := apys.ApyMap()
	for _, apy := range apys.Apys {
		key := apy.Address
		t.Logf("%v apy: %v", key, apyMap[key])
	}
}
