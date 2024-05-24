package sui_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	digest, err := sui_types.NewDigest("3vVi8XZgNpzQ34PFgwJTQqWtPMU84njcBX1EUxUHhyDk")
	require.NoError(t, err)
	events, err := client.GetEvents(context.Background(), digest)
	require.NoError(t, err)
	require.Len(t, events, 1)
	for _, event := range events {
		require.Equal(t, digest, &event.Id.TxDigest)
		require.Equal(
			t,
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			&event.PackageId,
		)
		require.Equal(t, "clob_v2", event.TransactionModule)
		require.Equal(
			t,
			sui_types.MustSuiAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
			&event.Sender,
		)
		require.Equal(
			t,
			"0xdee9::clob_v2::OrderPlaced<0x2::sui::SUI, 0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN>",
			event.Type,
		)
		// TODO check ParsedJson map
	}
}

func TestGetLatestCheckpointSequenceNumber(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	sequenceNumber, err := client.GetLatestCheckpointSequenceNumber(context.Background())
	require.NoError(t, err)
	num, err := strconv.Atoi(sequenceNumber)
	require.NoError(t, err)
	require.Greater(t, num, 34317507)
}

func TestGetObject(t *testing.T) {
	type args struct {
		ctx   context.Context
		objID *sui_types.ObjectID
	}
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	coins, err := api.GetCoins(context.TODO(), sui_signer.TEST_ADDRESS, nil, nil, 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		api     *sui.ImplSuiAPI
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet",
			api:  api,
			args: args{
				ctx:   context.TODO(),
				objID: coins.Data[0].CoinObjectID,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.api.GetObject(
					tt.args.ctx, tt.args.objID, &models.SuiObjectDataOptions{
						ShowType:                true,
						ShowOwner:               true,
						ShowContent:             true,
						ShowDisplay:             true,
						ShowBcs:                 true,
						ShowPreviousTransaction: true,
						ShowStorageRebate:       true,
					},
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetObject() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%+v", got)
			},
		)
	}
}

func TestGetTotalTransactionBlocks(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	res, err := api.GetTotalTransactionBlocks(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

func TestGetTransactionBlock(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	digest, err := sui_types.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
	require.NoError(t, err)
	resp, err := client.GetTransactionBlock(
		context.Background(), digest, &models.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
	)
	require.NoError(t, err)

	gasCostSummary := models.GasCostSummary{
		ComputationCost:         models.NewSafeSuiBigInt(uint64(750000)),
		StorageCost:             models.NewSafeSuiBigInt(uint64(32383600)),
		StorageRebate:           models.NewSafeSuiBigInt(uint64(21955032)),
		NonRefundableStorageFee: models.NewSafeSuiBigInt(uint64(221768)),
	}

	require.Equal(t, digest, &resp.Digest)

	require.True(t, resp.Effects.Data.IsSuccess())
	require.Equal(t, int64(183), resp.Effects.Data.V1.ExecutedEpoch.Int64())
	require.Equal(t, gasCostSummary, resp.Effects.Data.V1.GasUsed)
	require.Equal(t, int64(11178568), resp.Effects.Data.GasFee())
	// TODO check all the fields
}

func TestMultiGetObjects(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	coins, err := api.GetCoins(context.TODO(), sui_signer.TEST_ADDRESS, nil, nil, 1)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectID
	objs := []*sui_types.ObjectID{obj, obj}
	resp, err := api.MultiGetObjects(
		context.Background(), objs, &models.SuiObjectDataOptions{
			ShowType:                true,
			ShowOwner:               true,
			ShowContent:             true,
			ShowDisplay:             true,
			ShowBcs:                 true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, len(objs), len(resp))
	require.Equal(t, resp[0], resp[1])
}

func TestTryGetPastObject(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	objId, err := sui_types.SuiAddressFromHex("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d")
	require.NoError(t, err)
	data, err := api.TryGetPastObject(context.Background(), objId, 903, nil)
	require.NoError(t, err)
	t.Log(data)
}
