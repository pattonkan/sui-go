package suiclient_test

import (
	"context"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/stretchr/testify/require"
)

func TestGetChainIdentifier(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	chainId, err := client.GetChainIdentifier(context.Background())
	require.NoError(t, err)
	require.Equal(t, conn.ChainIdentifierSuiMainnet, chainId)
}

func TestGetCheckpoint(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	checkpoint, err := client.GetCheckpoint(context.Background(), sui.NewBigInt(1000))
	require.NoError(t, err)
	targetCheckpoint := &suiclient.Checkpoint{
		Epoch:                    sui.NewBigInt(0),
		SequenceNumber:           sui.NewBigInt(1000),
		Digest:                   *sui.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
		NetworkTotalTransactions: sui.NewBigInt(1001),
		PreviousDigest:           sui.MustNewDigest("41nPNZWHvvajmBQjX3GbppsgGZDEB6DhN4UxPkjSYRRj"),
		EpochRollingGasCostSummary: suiclient.GasCostSummary{
			ComputationCost:         sui.NewBigInt(0),
			StorageCost:             sui.NewBigInt(0),
			StorageRebate:           sui.NewBigInt(0),
			NonRefundableStorageFee: sui.NewBigInt(0),
		},
		TimestampMs:           sui.NewBigInt(1681393657483),
		Transactions:          []*sui.Digest{sui.MustNewDigest("9NnjyPG8V2TPCSbNE391KDyge42AwV3vUD7aNtQQ9eqS")},
		CheckpointCommitments: []suiclient.CheckpointCommitment{},
		ValidatorSignature:    *sui.MustNewBase64("r8/5+Rm7niIlndcnvjSJ/vZLPrH3xY/ePGYTvrVbTascoQSpS+wsGlC+bQBpzIwA"),
	}
	require.Equal(t, targetCheckpoint, checkpoint)
}

func TestGetCheckpoints(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	cursor := sui.NewBigInt(999)
	limit := uint64(2)
	checkpointPage, err := client.GetCheckpoints(context.Background(), &suiclient.GetCheckpointsRequest{
		Cursor: cursor,
		Limit:  &limit,
	})
	require.NoError(t, err)
	targetCheckpoints := []*suiclient.Checkpoint{
		{
			Epoch:                    sui.NewBigInt(0),
			SequenceNumber:           sui.NewBigInt(1000),
			Digest:                   *sui.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
			NetworkTotalTransactions: sui.NewBigInt(1001),
			PreviousDigest:           sui.MustNewDigest("41nPNZWHvvajmBQjX3GbppsgGZDEB6DhN4UxPkjSYRRj"),
			EpochRollingGasCostSummary: suiclient.GasCostSummary{
				ComputationCost:         sui.NewBigInt(0),
				StorageCost:             sui.NewBigInt(0),
				StorageRebate:           sui.NewBigInt(0),
				NonRefundableStorageFee: sui.NewBigInt(0),
			},
			TimestampMs:           sui.NewBigInt(1681393657483),
			Transactions:          []*sui.Digest{sui.MustNewDigest("9NnjyPG8V2TPCSbNE391KDyge42AwV3vUD7aNtQQ9eqS")},
			CheckpointCommitments: []suiclient.CheckpointCommitment{},
			ValidatorSignature:    *sui.MustNewBase64("r8/5+Rm7niIlndcnvjSJ/vZLPrH3xY/ePGYTvrVbTascoQSpS+wsGlC+bQBpzIwA"),
		},
		{
			Epoch:                    sui.NewBigInt(0),
			SequenceNumber:           sui.NewBigInt(1001),
			Digest:                   *sui.MustNewDigest("8umKe5Ae2TAH5ySw2zeEua8cTeeTFZV8F3GfFViZ5cq3"),
			NetworkTotalTransactions: sui.NewBigInt(1002),
			PreviousDigest:           sui.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
			EpochRollingGasCostSummary: suiclient.GasCostSummary{
				ComputationCost:         sui.NewBigInt(0),
				StorageCost:             sui.NewBigInt(0),
				StorageRebate:           sui.NewBigInt(0),
				NonRefundableStorageFee: sui.NewBigInt(0),
			},
			TimestampMs:           sui.NewBigInt(1681393661034),
			Transactions:          []*sui.Digest{sui.MustNewDigest("9muLz7ZTocpBTdSo5Ak7ZxzEpfzywr6Y12Hj3AdT8dvV")},
			CheckpointCommitments: []suiclient.CheckpointCommitment{},
			ValidatorSignature:    *sui.MustNewBase64("jG5ViKThziBpnJnOw9dVdjIrv2IHhCrn8ZhvI1gUS2X1t90aRqhnLF6+WbS1S2WT"),
		},
	}
	require.Len(t, checkpointPage.Data, 2)
	require.Equal(t, checkpointPage.Data, targetCheckpoints)
	require.Equal(t, true, checkpointPage.HasNextPage)
	require.Equal(t, sui.NewBigInt(1001), checkpointPage.NextCursor)
}

func TestGetEvents(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	digest, err := sui.NewDigest("3vVi8XZgNpzQ34PFgwJTQqWtPMU84njcBX1EUxUHhyDk")
	require.NoError(t, err)
	events, err := client.GetEvents(context.Background(), digest)
	require.NoError(t, err)
	require.Len(t, events, 1)
	for _, event := range events {
		require.Equal(t, digest, &event.Id.TxDigest)
		require.Equal(
			t,
			sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			event.PackageId,
		)
		require.Equal(t, "clob_v2", event.TransactionModule)
		require.Equal(
			t,
			sui.MustAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
			event.Sender,
		)
		targetStructTag := sui.StructTag{
			Address: sui.MustAddressFromHex("0xdee9"),
			Module:  sui.Identifier("clob_v2"),
			Name:    sui.Identifier("OrderPlaced"),
			TypeParams: []sui.TypeTag{
				{Struct: &sui.StructTag{
					Address: sui.MustAddressFromHex("0x2"),
					Module:  sui.Identifier("sui"),
					Name:    sui.Identifier("SUI"),
				}},
				{Struct: &sui.StructTag{
					Address: sui.MustAddressFromHex("0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf"),
					Module:  sui.Identifier("coin"),
					Name:    sui.Identifier("COIN"),
				}},
			},
		}
		require.Equal(t, targetStructTag.Address, event.Type.Address)
		require.Equal(t, targetStructTag.Module, event.Type.Module)
		require.Equal(t, targetStructTag.Name, event.Type.Name)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Address, event.Type.TypeParams[0].Struct.Address)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Module, event.Type.TypeParams[0].Struct.Module)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Name, event.Type.TypeParams[0].Struct.Name)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.TypeParams, event.Type.TypeParams[0].Struct.TypeParams)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Address, event.Type.TypeParams[1].Struct.Address)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Module, event.Type.TypeParams[1].Struct.Module)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Name, event.Type.TypeParams[1].Struct.Name)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.TypeParams, event.Type.TypeParams[1].Struct.TypeParams)
		targetBcsBase85 := base58.Decode("yNS5iDS3Gvdo3DhXdtFpuTS12RrSiNkrvjcm2rejntCuqWjF1DdwnHgjowdczAkR18LQHcBqbX2tWL76rys9rTCzG6vm7Tg34yqUkpFSMqNkcS6cfWbN8SdVsxn5g4ZEQotdBgEFn8yN7hVZ7P1MKvMwWf")
		require.Equal(t, targetBcsBase85, event.Bcs.Data())
		// TODO check ParsedJson map
	}
}

func TestGetLatestCheckpointSequenceNumber(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	sequenceNumber, err := client.GetLatestCheckpointSequenceNumber(context.Background())
	require.NoError(t, err)
	num, err := strconv.Atoi(sequenceNumber)
	require.NoError(t, err)
	require.Greater(t, num, 34317507)
}

func TestGetObject(t *testing.T) {
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	object, err := api.GetObject(context.TODO(), &suiclient.GetObjectRequest{
		ObjectId: sui.MustObjectIdFromHex("0x6"),
		Options:  &suiclient.SuiObjectDataOptions{ShowBcs: true},
	})
	require.NoError(t, err)
	require.NotNil(t, object)
}

func TestGetProtocolConfig(t *testing.T) {
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	version := sui.NewBigInt(47)
	protocolConfig, err := api.GetProtocolConfig(context.Background(), version)
	require.NoError(t, err)
	require.Equal(t, uint64(47), protocolConfig.ProtocolVersion.Uint64())
}

func TestGetTotalTransactionBlocks(t *testing.T) {
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	res, err := api.GetTotalTransactionBlocks(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

func TestGetTransactionBlock(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	digest, err := sui.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
	require.NoError(t, err)
	resp, err := client.GetTransactionBlock(
		context.Background(), &suiclient.GetTransactionBlockRequest{
			Digest: digest,
			Options: &suiclient.SuiTransactionBlockResponseOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowRawEffects:     true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
				ShowEvents:         true,
			},
		},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	targetGasCostSummary := suiclient.GasCostSummary{
		ComputationCost:         sui.NewBigInt(750000),
		StorageCost:             sui.NewBigInt(32383600),
		StorageRebate:           sui.NewBigInt(21955032),
		NonRefundableStorageFee: sui.NewBigInt(221768),
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
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	coins, err := api.GetCoins(context.TODO(), &suiclient.GetCoinsRequest{
		Owner: suisigner.TEST_ADDRESS,
		Limit: 1,
	})
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectId
	objs := []*sui.ObjectId{obj, obj}
	resp, err := api.MultiGetObjects(
		context.Background(), &suiclient.MultiGetObjectsRequest{
			ObjectIds: objs,
			Options: &suiclient.SuiObjectDataOptions{
				ShowType:                true,
				ShowOwner:               true,
				ShowContent:             true,
				ShowDisplay:             true,
				ShowBcs:                 true,
				ShowPreviousTransaction: true,
				ShowStorageRebate:       true,
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, len(objs), len(resp))
	require.Equal(t, resp[0], resp[1])
}

func TestMultiGetTransactionBlocks(t *testing.T) {
	client := suiclient.NewClient(conn.TestnetEndpointUrl)

	resp, err := client.MultiGetTransactionBlocks(
		context.Background(),
		&suiclient.MultiGetTransactionBlocksRequest{
			Digests: []*sui.Digest{
				sui.MustNewDigest("6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr"),
				sui.MustNewDigest("8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5"),
			},
			Options: &suiclient.SuiTransactionBlockResponseOptions{
				ShowEffects: true,
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, resp, 2)
	require.Equal(t, "6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr", resp[0].Digest.String())
	require.Equal(t, "8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5", resp[1].Digest.String())
}

func TestTryGetPastObject(t *testing.T) {
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	// there is no software-level guarantee/SLA that objects with past versions can be retrieved by this API
	resp, err := api.TryGetPastObject(context.Background(), &suiclient.TryGetPastObjectRequest{
		ObjectId: sui.MustObjectIdFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
		Version:  187584506,
		Options: &suiclient.SuiObjectDataOptions{
			ShowType:  true,
			ShowOwner: true,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Data.VersionNotFound)
}

func TestTryMultiGetPastObjects(t *testing.T) {
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	req := []*suiclient.SuiGetPastObjectRequest{
		{
			ObjectId: sui.MustObjectIdFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  sui.NewBigInt(187584506),
		},
		{
			ObjectId: sui.MustObjectIdFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  sui.NewBigInt(187584500),
		},
	}
	// there is no software-level guarantee/SLA that objects with past versions can be retrieved by this API
	resp, err := api.TryMultiGetPastObjects(context.Background(), &suiclient.TryMultiGetPastObjectsRequest{
		PastObjects: req,
		Options: &suiclient.SuiObjectDataOptions{
			ShowType:  true,
			ShowOwner: true,
		},
	})
	require.NoError(t, err)
	for _, data := range resp {
		require.NotNil(t, data.Data.VersionNotFound)
	}
}
