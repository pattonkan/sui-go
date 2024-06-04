package sui_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetCommitteeInfo(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	epochId := models.NewSafeSuiBigInt(uint64(400))
	committeeInfo, err := client.GetCommitteeInfo(context.Background(), &epochId)
	require.NoError(t, err)
	require.Equal(t, epochId, committeeInfo.EpochId)
	// just use a arbitrary big number to ensure there are enough validator
	require.Greater(t, len(committeeInfo.Validators), 20)
}

func TestGetLatestSuiSystemState(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	state, err := client.GetLatestSuiSystemState(context.Background())
	require.NoError(t, err)
	require.NotNil(t, state)
}

func TestGetReferenceGasPrice(t *testing.T) {
	client := sui.NewSuiClient(conn.DevnetEndpointUrl)
	gasPrice, err := client.GetReferenceGasPrice(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, gasPrice.Int64(), int64(1000))
}

func TestGetStakes(t *testing.T) {
	// FIXME change the valid staking sui address
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	address, err := sui_types.SuiAddressFromHex("0x8ecaf4b95b3c82c712d3ddb22e7da88d2286c4653f3753a86b6f7a216a3ca518")
	require.NoError(t, err)
	stakes, err := client.GetStakes(context.Background(), address)
	require.NoError(t, err)
	require.Greater(t, len(stakes), 0)
	for _, validator := range stakes {
		require.Equal(t, address, &validator.ValidatorAddress)
		for _, stake := range validator.Stakes {
			if stake.Data.StakeStatus.Data.Active != nil {
				t.Logf(
					"earned amount %10v at %v",
					stake.Data.StakeStatus.Data.Active.EstimatedReward.Uint64(),
					validator.ValidatorAddress,
				)
			}
		}
	}
}

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
