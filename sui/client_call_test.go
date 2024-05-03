package sui_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/howjmay/go-sui-sdk/lib"
	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui_types"

	"github.com/howjmay/go-sui-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_TagJson_Owner(t *testing.T) {
	test := func(str string) lib.TagJson[sui_types.Owner] {
		var s lib.TagJson[sui_types.Owner]
		data := []byte(str)
		err := json.Unmarshal(data, &s)
		require.NoError(t, err)
		return s
	}
	{
		v := test(`"Immutable"`).Data
		require.Nil(t, v.AddressOwner)
		require.Nil(t, v.ObjectOwner)
		require.Nil(t, v.Shared)
		require.NotNil(t, v.Immutable)
	}
	{
		v := test(`{"AddressOwner": "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e"}`).Data
		require.NotNil(t, v.AddressOwner)
		require.Nil(t, v.ObjectOwner)
		require.Nil(t, v.Shared)
		require.Nil(t, v.Immutable)
	}
}

func TestClient_DryRunTransaction(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	signer := sui_types.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)
	tx, err := api.PayAllSui(
		context.Background(), signer, signer,
		pickedCoins.CoinIds(),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	resp, err := api.DryRunTransaction(context.Background(), tx.TxBytes)
	require.NoError(t, err)
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Empty(t, resp.Effects.Data.V1.Status.Error)
}

// TestClient_ExecuteTransactionSerializedSig
// This test case will affect the real coin in the test case of account
// temporary disabled
//func TestClient_ExecuteTransactionSerializedSig(t *testing.T) {
//	api := sui.NewSuiClient(DevnetClient(t))
//	coins, err := api.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := api.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectID, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := api.ExecuteTransactionSerializedSig(context.TODO(), *signedTx, types.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}

//func TestClient_ExecuteTransaction(t *testing.T) {
//	api := sui.NewSuiClient(DevnetClient(t))
//	coins, err := api.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := api.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectID, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := api.ExecuteTransaction(context.TODO(), *signedTx, types.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}

func TestClient_BatchGetObjectsOwnedByAddress(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))

	options := types.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
	}
	coinType := fmt.Sprintf("0x2::coin::Coin<%v>", types.SuiCoinType)
	filterObject, err := api.BatchGetObjectsOwnedByAddress(context.TODO(), sui_types.TEST_ADDRESS, &options, coinType)
	require.NoError(t, err)
	t.Log(filterObject)
}

func TestClient_GetCoinMetadata(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	metadata, err := api.GetCoinMetadata(context.TODO(), types.SuiCoinType)
	require.NoError(t, err)
	t.Logf("%#v", metadata)
}

func TestClient_GetAllBalances(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	balances, err := api.GetAllBalances(context.TODO(), sui_types.TEST_ADDRESS)
	require.NoError(t, err)
	for _, balance := range balances {
		t.Logf(
			"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
			balance.CoinType, balance.CoinObjectCount,
			balance.TotalBalance.String(), balance.LockedBalance,
		)
	}
}

func TestClient_GetBalance(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	balance, err := api.GetBalance(context.TODO(), sui_types.TEST_ADDRESS, "")
	require.NoError(t, err)
	t.Logf(
		"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
		balance.CoinType, balance.CoinObjectCount,
		balance.TotalBalance.String(), balance.LockedBalance,
	)
}

func TestClient_GetCoins(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	defaultCoinType := types.SuiCoinType
	coins, err := api.GetCoins(context.TODO(), sui_types.TEST_ADDRESS, &defaultCoinType, nil, 1)
	require.NoError(t, err)
	t.Logf("%#v", coins)
}

func TestClient_GetAllCoins(t *testing.T) {
	type args struct {
		ctx     context.Context
		address *sui_types.SuiAddress
		cursor  *sui_types.ObjectID
		limit   uint
	}
	tests := []struct {
		name    string
		a       *sui.ImplSuiAPI
		args    args
		want    *types.CoinPage
		wantErr bool
	}{
		{
			name: "test case 1",
			a:    sui.NewSuiClient(DevnetClient(t)),
			args: args{
				ctx:     context.TODO(),
				address: sui_types.TEST_ADDRESS,
				cursor:  nil,
				limit:   3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.a.GetAllCoins(tt.args.ctx, tt.args.address, tt.args.cursor, tt.args.limit)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAllCoins() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestClient_GetTransaction(t *testing.T) {
	api := sui.NewSuiClient(MainnetClient(t))
	digest := "D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL"
	d, err := sui_types.NewDigest(digest)
	require.NoError(t, err)
	resp, err := api.GetTransactionBlock(
		context.Background(), *d, types.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
	)
	require.NoError(t, err)
	t.Logf("%#v", resp)

	require.Equal(t, int64(11178568), resp.Effects.Data.GasFee())
}

func TestBatchCall_GetObject(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	if false {
		// get specified object
		idstr := "0x4ad2f0a918a241d6a19573212aeb56947bb9255a14e921a7ec78b262536826f0"
		objId, err := sui_types.NewAddressFromHex(idstr)
		require.NoError(t, err)
		obj, err := api.GetObject(
			context.Background(), objId, &types.SuiObjectDataOptions{
				ShowType:    true,
				ShowContent: true,
			},
		)
		require.NoError(t, err)
		t.Log(obj.Data)
	}

	coins, err := api.GetCoins(context.TODO(), sui_types.TEST_ADDRESS, nil, nil, 3)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		return
	}
	objId := coins.Data[0].CoinObjectID
	obj, err := api.GetObject(context.Background(), &objId, nil)
	require.NoError(t, err)
	t.Log(obj.Data)
}

func TestClient_GetObject(t *testing.T) {
	type args struct {
		ctx   context.Context
		objID *sui_types.ObjectID
	}
	api := sui.NewSuiClient(DevnetClient(t))
	coins, err := api.GetCoins(context.TODO(), sui_types.TEST_ADDRESS, nil, nil, 1)
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
				objID: &coins.Data[0].CoinObjectID,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.api.GetObject(
					tt.args.ctx, tt.args.objID, &types.SuiObjectDataOptions{
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

func TestClient_MultiGetObjects(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	coins, err := api.GetCoins(context.TODO(), sui_types.TEST_ADDRESS, nil, nil, 1)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectID
	objs := []sui_types.ObjectID{obj, obj}
	resp, err := api.MultiGetObjects(
		context.Background(), objs, &types.SuiObjectDataOptions{
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

func TestClient_GetOwnedObjects(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	obj, err := sui_types.NewAddressFromHex("0x2")
	require.NoError(t, err)
	query := types.SuiObjectResponseQuery{
		Filter: &types.SuiObjectDataFilter{
			Package: obj,
			// StructType: "0x2::coin::Coin<0x2::sui::SUI>",
		},
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	limit := uint(1)
	objs, err := api.GetOwnedObjects(context.Background(), sui_types.TEST_ADDRESS, &query, nil, &limit)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(objs.Data), int(limit))
}

func TestClient_GetTotalSupply(t *testing.T) {
	type args struct {
		ctx      context.Context
		coinType string
	}
	tests := []struct {
		name    string
		api     *sui.ImplSuiAPI
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "test 1",
			api:  sui.NewSuiClient(DevnetClient(t)),
			args: args{
				context.TODO(),
				types.SuiCoinType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.api.GetTotalSupply(tt.args.ctx, tt.args.coinType)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTotalSupply() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%d", got)
			},
		)
	}
}
func TestClient_GetTotalTransactionBlocks(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	res, err := api.GetTotalTransactionBlocks(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

func TestClient_GetLatestCheckpointSequenceNumber(t *testing.T) {
	api := sui.NewSuiClient(MainnetClient(t))
	res, err := api.GetLatestCheckpointSequenceNumber(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

//func TestClient_Publish(t *testing.T) {
//	api := sui.NewSuiClient(DevnetClient(t))
//	dmens, err := types.NewBase64Data(DmensDmensB64)
//	require.NoError(t, err)
//	profile, err := types.NewBase64Data(DmensProfileB64)
//	require.NoError(t, err)
//	coins, err := api.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(30000)
//	require.NoError(t, err)
//	type args struct {
//		ctx             context.Context
//		address         types.Address
//		compiledModules []*types.Base64Data
//		gas             types.ObjectID
//		gasBudget       uint
//	}
//	tests := []struct {
//		name    string
//		client  *client.Client
//		args    args
//		want    *types.TransactionBytes
//		wantErr bool
//	}{
//		{
//			name:   "test for dmens publish",
//			client: chain,
//			args: args{
//				ctx:             context.TODO(),
//				address:         *Address,
//				compiledModules: []*types.Base64Data{dmens, profile},
//				gas:             coin.CoinObjectID,
//				gasBudget:       30000,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := tt.client.Publish(tt.args.ctx, tt.args.address, tt.args.compiledModules, tt.args.gas, tt.args.gasBudget)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Publish() error: %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			t.Logf("%#v", got)
//
//			txResult, err := tt.client.DryRunTransaction(context.TODO(), got)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Publish() error: %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			t.Logf("%#v", txResult)
//		})
//	}
//}

func TestClient_TryGetPastObject(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	objId, err := sui_types.NewAddressFromHex("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d")
	require.NoError(t, err)
	data, err := api.TryGetPastObject(context.Background(), objId, 903, nil)
	require.NoError(t, err)
	t.Log(data)
}

func TestClient_GetEvents(t *testing.T) {
	api := sui.NewSuiClient(MainnetClient(t))
	digest := "D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL"
	d, err := sui_types.NewDigest(digest)
	require.NoError(t, err)
	res, err := api.GetEvents(context.Background(), *d)
	require.NoError(t, err)
	t.Log(res)
}

func TestClient_GetReferenceGasPrice(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	gasPrice, err := api.GetReferenceGasPrice(context.Background())
	require.NoError(t, err)
	t.Logf("current gas price: %v", gasPrice)
}

// func TestClient_DevInspectTransactionBlock(t *testing.T) {
// 	api := sui.NewSuiClient(DevnetClient(t))
// 	signer := Address
// 	price, err := api.GetReferenceGasPrice(context.TODO())
// 	require.NoError(t, err)
// 	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
// 	require.NoError(t, err)

// 	amount := SUI(0.01).Int64()
// 	gasBudget := SUI(0.01).Uint64()
// 	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(amount * 2), 0, false)
// 	require.NoError(t, err)
// 	tx, err := api.PayAllSui(context.Background(),
// 		signer, signer,
// 		pickedCoins.CoinIds(),
// 		types.NewSafeSuiBigInt(gasBudget))
// 	require.NoError(t, err)

// 	resp, err := api.DevInspectTransactionBlock(context.Background(), signer, tx.TxBytes, price, nil)
// 	require.NoError(t, err)
// 	t.Log(resp)
// }

func TestClient_QueryTransactionBlocks(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           types.SuiTransactionBlockResponseQuery
		cursor          *sui_types.TransactionDigest
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.TransactionBlocksPage
		wantErr bool
	}{
		{
			name: "test for queryTransactionBlocks",
			args: args{
				ctx: context.TODO(),
				query: types.SuiTransactionBlockResponseQuery{
					Filter: &types.TransactionFilter{
						FromAddress: sui_types.TEST_ADDRESS,
					},
					Options: &types.SuiTransactionBlockResponseOptions{
						ShowInput:   true,
						ShowEffects: true,
					},
				},
				cursor:          nil,
				limit:           &limit,
				descendingOrder: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.QueryTransactionBlocks(
					tt.args.ctx,
					tt.args.query,
					tt.args.cursor,
					tt.args.limit,
					tt.args.descendingOrder,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryTransactionBlocks() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestClient_ResolveNameServiceAddress(t *testing.T) {
	api := sui.NewSuiClient(MainnetClient(t))
	addr, err := api.ResolveNameServiceAddress(context.Background(), "2222.sui")
	require.NoError(t, err)
	require.Equal(t, "0x6174c5bd8ab9bf492e159a64e102de66429cfcde4fa883466db7b03af28b3ce9", addr.String())

	_, err = api.ResolveNameServiceAddress(context.Background(), "2222.suijjzzww")
	require.ErrorContains(t, err, "not found")
}

func TestClient_ResolveNameServiceNames(t *testing.T) {
	api := sui.NewSuiClient(MainnetClient(t))
	owner := AddressFromStrMust("0x57188743983628b3474648d8aa4a9ee8abebe8f6816243773d7e8ed4fd833a28")
	namePage, err := api.ResolveNameServiceNames(context.Background(), owner, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, namePage.Data)
	t.Log(namePage.Data)

	owner = AddressFromStrMust("0x57188743983628b3474648d8aa4a9ee8abebe8f681")
	namePage, err = api.ResolveNameServiceNames(context.Background(), owner, nil, nil)
	require.NoError(t, err)
	require.Empty(t, namePage.Data)
}

func TestClient_QueryEvents(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           types.EventFilter
		cursor          *types.EventId
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.EventPage
		wantErr bool
	}{
		{
			name: "test for query events",
			args: args{
				ctx: context.TODO(),
				query: types.EventFilter{
					Sender: sui_types.TEST_ADDRESS,
				},
				cursor:          nil,
				limit:           &limit,
				descendingOrder: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.QueryEvents(
					tt.args.ctx,
					tt.args.query,
					tt.args.cursor,
					tt.args.limit,
					tt.args.descendingOrder,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryEvents() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Log(got)
			},
		)
	}
}

func TestClient_GetDynamicFields(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	parentObjectID, err := sui_types.NewAddressFromHex("0x1719957d7a2bf9d72459ff0eab8e600cbb1991ef41ddd5b4a8c531035933d256")
	require.NoError(t, err)
	limit := uint(5)
	type args struct {
		ctx            context.Context
		parentObjectID *sui_types.ObjectID
		cursor         *sui_types.ObjectID
		limit          *uint
	}
	tests := []struct {
		name    string
		args    args
		want    *types.DynamicFieldPage
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: parentObjectID,
				cursor:         nil,
				limit:          &limit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.GetDynamicFields(tt.args.ctx, tt.args.parentObjectID, tt.args.cursor, tt.args.limit)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetDynamicFields() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Log(got)
			},
		)
	}
}

func TestClient_GetDynamicFieldObject(t *testing.T) {
	api := sui.NewSuiClient(DevnetClient(t))
	parentObjectID, err := sui_types.NewAddressFromHex("0x1719957d7a2bf9d72459ff0eab8e600cbb1991ef41ddd5b4a8c531035933d256")
	require.NoError(t, err)
	type args struct {
		ctx            context.Context
		parentObjectID *sui_types.ObjectID
		name           sui_types.DynamicFieldName
	}
	tests := []struct {
		name    string
		args    args
		want    *types.SuiObjectResponse
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: parentObjectID,
				name: sui_types.DynamicFieldName{
					Type:  "address",
					Value: "0xf9ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.GetDynamicFieldObject(tt.args.ctx, tt.args.parentObjectID, tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetDynamicFieldObject() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}
