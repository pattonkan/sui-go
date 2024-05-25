package sui_test

import (
	"context"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/btcsuite/btcutil/base58"
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
		targetBcsBase85 := base58.Decode("yNS5iDS3Gvdo3DhXdtFpuTS12RrSiNkrvjcm2rejntCuqWjF1DdwnHgjowdczAkR18LQHcBqbX2tWL76rys9rTCzG6vm7Tg34yqUkpFSMqNkcS6cfWbN8SdVsxn5g4ZEQotdBgEFn8yN7hVZ7P1MKvMwWf")
		require.Equal(t, targetBcsBase85, event.Bcs.Data())
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
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowRawEffects:     true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	targetGasCostSummary := models.GasCostSummary{
		ComputationCost:         new(models.BigInt).SetUint64(uint64(750000)),
		StorageCost:             new(models.BigInt).SetUint64(uint64(32383600)),
		StorageRebate:           new(models.BigInt).SetUint64(uint64(21955032)),
		NonRefundableStorageFee: new(models.BigInt).SetUint64(uint64(221768)),
	}
	require.Equal(t, digest, &resp.Digest)
	targetRawTxBase64, err := base64.StdEncoding.DecodeString("AQAAAAAACgEBpqVCwrKBCI6PELxQWossTD9mgGbIy8W++ipS7CWatqOAVmEAAAAAAAEBAG85p+0UjVUsc5qkxhWSZ/qr2vghuqeSNiZr1gQzhCIAV3XJAQAAAAAgKEbgAIwWMBRZ1grRBFQ6qrSWLHa/AfKG8ubjmkxM/zoAIEnHBYEE/EtGK3r1lzrUU9QPAiTHLBd2+R8GS7k042UqAQF/3Yg8C3Qn8YzbSYxMh6SnnWvsR4PLPyGqOBa7xkzo7wDr5AEAAAAAAQEBbg3e/ArZiInAS6uWOeUSwhdmxeY2b4nmlpVtm+aVKHENAAAAAAAAAAEAERAyMjIyMjIyMjIyMjIuc3VpAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgEAAAAAAAAAAAAgVxiHQ5g2KLNHRkjYqkqe6Kvr6PaBYkN3PX6O1P2DOigAERAyMjIyMjIyMjIyMjIuc3VpACBXGIdDmDYos0dGSNiqSp7oq+vo9oFiQ3c9fo7U/YM6KAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIFa2lvc2sKYm9ycm93X3ZhbAEH7klqDMBNBqNFmCumaXyQxhkCDenidECMeBn3h/9m4aEIc3VpZnJlbnMHU3VpRnJlbgEHiJT6AvxvNsvEha6RRdBfJHp44iCBT7hBmrJhvYHwjzIJYnVsbHNoYXJrCUJ1bGxzaGFyawADAQAAAQEAAQIAAGpuoUDgld3YL3x0WQUFSzIDEp3QSgnQN1QWwxFhky0tC2ZyZWVfY2xhaW1zCmZyZWVfY2xhaW0BB+5JagzATQajRZgrpml8kMYZAg3p4nRAjHgZ94f/ZuGhCHN1aWZyZW5zB1N1aUZyZW4BB4iU+gL8bzbLxIWukUXQXyR6eOIggU+4QZqyYb2B8I8yCWJ1bGxzaGFyawlCdWxsc2hhcmsABQEDAAEEAAMAAAAAAQUAAQYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACBWtpb3NrCnJldHVybl92YWwBB+5JagzATQajRZgrpml8kMYZAg3p4nRAjHgZ94f/ZuGhCHN1aWZyZW5zB1N1aUZyZW4BB4iU+gL8bzbLxIWukUXQXyR6eOIggU+4QZqyYb2B8I8yCWJ1bGxzaGFyawlCdWxsc2hhcmsAAwEAAAMAAAAAAwAAAQAA2sImUutAC+sfXiEmRZyuju3BFrc7itYLcePo1/2zF+IMZGlyZWN0X3NldHVwEnNldF90YXJnZXRfYWRkcmVzcwAEAQQAAgEAAQcAAQYAANrCJlLrQAvrH14hJkWcro7twRa3O4rWC3Hj6Nf9sxfiDGRpcmVjdF9zZXR1cBJzZXRfcmV2ZXJzZV9sb29rdXAAAgEEAAEIAAEBAgEAAQkAVxiHQ5g2KLNHRkjYqkqe6Kvr6PaBYkN3PX6O1P2DOigBAIV+3vABgFUzNcciYyljcM6zXwvwuD9FeVw6JU3rDUD/YO8BAAAAACBmxGapu4poDXYNHxLCokFFdgFBwBhoQW8vcK8+XuklpFcYh0OYNiizR0ZI2KpKnuir6+j2gWJDdz1+jtT9gzoo7gIAAAAAAADA8MQAAAAAAAABYQBao7U4xuiDfVJM+YnHs7cBOs9VJJVriNBdHr7neIyT+M9tzPcRbANj2P9q2s21wtgIiNtayH6IAAhgFEhKsEANMFE7Y3jZzVZy0dJdgxaL8YB9JBE0745Io7/8t/XlJ3w=")
	require.NoError(t, err)
	require.Equal(t, targetRawTxBase64, resp.RawTransaction.Data())
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Equal(t, int64(183), resp.Effects.Data.V1.ExecutedEpoch.Int64())
	require.Equal(t, targetGasCostSummary, resp.Effects.Data.V1.GasUsed)
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
