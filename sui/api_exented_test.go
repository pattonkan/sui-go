package sui_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetDynamicFieldObject(t *testing.T) {
	t.Skip("FIXME")
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	parentObjectID, err := sui_types.SuiAddressFromHex("0x1719957d7a2bf9d72459ff0eab8e600cbb1991ef41ddd5b4a8c531035933d256")
	require.NoError(t, err)
	type args struct {
		ctx            context.Context
		parentObjectID *sui_types.ObjectID
		name           *sui_types.DynamicFieldName
	}
	tests := []struct {
		name    string
		args    args
		want    *models.SuiObjectResponse
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: parentObjectID,
				name: &sui_types.DynamicFieldName{
					Type:  "address",
					Value: "0xf9ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.GetDynamicFieldObject(tt.args.ctx, &models.GetDynamicFieldObjectRequest{
					ParentObjectID: tt.args.parentObjectID,
					Name:           tt.args.name,
				})
				if (err != nil) != tt.wantErr {
					t.Errorf("GetDynamicFieldObject() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestGetDynamicFields(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	limit := 5
	type args struct {
		ctx            context.Context
		parentObjectID *sui_types.ObjectID
		cursor         *sui_types.ObjectID
		limit          *uint
	}
	tests := []struct {
		name    string
		args    args
		want    *models.DynamicFieldPage
		wantErr error
	}{
		{
			name: "a deepbook shared object",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: sui_types.MustSuiAddressFromHex("0xa9d09452bba939b3172c0242d022274845cfe4e58648b73dd33b3d5b823dc8ae"),
				cursor:         nil,
				limit:          func() *uint { tmpLimit := uint(limit); return &tmpLimit }(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := client.GetDynamicFields(tt.args.ctx, &models.GetDynamicFieldsRequest{
					ParentObjectID: tt.args.parentObjectID,
					Cursor:         tt.args.cursor,
					Limit:          tt.args.limit,
				})
				require.ErrorIs(t, err, tt.wantErr)
				// object ID is '0x4405b50d791fd3346754e8171aaab6bc2ed26c2c46efdd033c14b30ae507ac33'
				// it has 'internal_nodes' field in type '0x2::table::Table<u64, 0xdee9::critbit::InternalNode'
				require.Len(t, got.Data, limit)
				for _, field := range got.Data {
					require.Equal(t, "u64", field.Name.Type)
					require.Equal(t, "0xdee9::critbit::InternalNode", field.ObjectType)
				}
			},
		)
	}
}

func TestGetOwnedObjects(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.NewSignerByIndex(sui_signer.TEST_SEED, sui_signer.KeySchemeFlagEd25519, 0)
	query := models.SuiObjectResponseQuery{
		Filter: &models.SuiObjectDataFilter{
			StructType: "0x2::coin::Coin<0x2::sui::SUI>",
		},
		Options: &models.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
		},
	}
	limit := uint(10)
	objs, err := api.GetOwnedObjects(context.Background(), &models.GetOwnedObjectsRequest{
		Address: signer.Address,
		Query:   &query,
		Cursor:  nil,
		Limit:   &limit,
	})
	require.NoError(t, err)
	require.Equal(t, len(objs.Data), int(limit))
	require.NoError(t, err)
	var fields models.CoinFields
	err = json.Unmarshal(objs.Data[9].Data.Content.Data.MoveObject.Fields, &fields)
	require.NoError(t, err)
	require.Equal(t, "1000000000", fields.Balance.String())
}

func TestQueryEvents(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	limit := 10

	type args struct {
		ctx             context.Context
		query           *models.EventFilter
		cursor          *models.EventId
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *models.EventPage
		wantErr error
	}{
		{
			name: "event in deepbook.batch_cancel_order()",
			args: args{
				ctx: context.TODO(),
				query: &models.EventFilter{
					Sender: sui_types.MustSuiAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
				},
				cursor:          nil,
				limit:           func() *uint { tmpLimit := uint(limit); return &tmpLimit }(),
				descendingOrder: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.QueryEvents(
					tt.args.ctx,
					&models.QueryEventsRequest{
						Query:           tt.args.query,
						Cursor:          tt.args.cursor,
						Limit:           tt.args.limit,
						DescendingOrder: tt.args.descendingOrder,
					},
				)
				require.ErrorIs(t, err, tt.wantErr)
				require.Len(t, got.Data, int(limit))

				for _, event := range got.Data {
					// FIXME we should change other filter to, so we can verify each fields of event more detailed.
					require.Equal(
						t,
						sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
						event.PackageId,
					)
					require.Equal(t, "clob_v2", event.TransactionModule)
					require.Equal(t, tt.args.query.Sender, event.Sender)
				}
			},
		)
	}
}

func TestQueryTransactionBlocks(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           *models.SuiTransactionBlockResponseQuery
		cursor          *sui_types.TransactionDigest
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *models.TransactionBlocksPage
		wantErr bool
	}{
		{
			name: "test for queryTransactionBlocks",
			args: args{
				ctx: context.TODO(),
				query: &models.SuiTransactionBlockResponseQuery{
					Filter: &models.TransactionFilter{
						FromAddress: sui_signer.TEST_ADDRESS,
					},
					Options: &models.SuiTransactionBlockResponseOptions{
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
					&models.QueryTransactionBlocksRequest{
						Query:           tt.args.query,
						Cursor:          tt.args.cursor,
						Limit:           tt.args.limit,
						DescendingOrder: tt.args.descendingOrder,
					},
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

func TestResolveNameServiceAddress(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	addr, err := api.ResolveNameServiceAddress(context.Background(), "2222.sui")
	require.NoError(t, err)
	require.Equal(t, "0x6174c5bd8ab9bf492e159a64e102de66429cfcde4fa883466db7b03af28b3ce9", addr.String())

	_, err = api.ResolveNameServiceAddress(context.Background(), "2222.suijjzzww")
	require.ErrorContains(t, err, "not found")
}

func TestResolveNameServiceNames(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	owner := sui_types.MustSuiAddressFromHex("0x57188743983628b3474648d8aa4a9ee8abebe8f6816243773d7e8ed4fd833a28")
	namePage, err := api.ResolveNameServiceNames(context.Background(), &models.ResolveNameServiceNamesRequest{
		Owner: owner,
	})
	require.NoError(t, err)
	require.NotEmpty(t, namePage.Data)
	t.Log(namePage.Data)

	owner = sui_types.MustSuiAddressFromHex("0x57188743983628b3474648d8aa4a9ee8abebe8f681")
	namePage, err = api.ResolveNameServiceNames(context.Background(), &models.ResolveNameServiceNamesRequest{
		Owner: owner,
	})
	require.NoError(t, err)
	require.Empty(t, namePage.Data)
}

func TestSubscribeEvent(t *testing.T) {
	t.Skip("passed at local side, but returned error on GitHub")
	api := sui.NewSuiWebsocketClient(conn.MainnetWebsocketEndpointUrl)

	type args struct {
		ctx      context.Context
		filter   *models.EventFilter
		resultCh chan models.SuiEvent
	}
	tests := []struct {
		name    string
		args    args
		want    *models.EventPage
		wantErr bool
	}{
		{
			name: "test for filter events",
			args: args{
				ctx: context.TODO(),
				filter: &models.EventFilter{
					Package: sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
				},
				resultCh: make(chan models.SuiEvent),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := api.SubscribeEvent(
					tt.args.ctx,
					tt.args.filter,
					tt.args.resultCh,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("SubscribeEvent() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				cnt := 0
				for results := range tt.args.resultCh {
					fmt.Println("results: ", results)
					// FIXME we need to check finite number request in details
					cnt++
					if cnt > 3 {
						break
					}
				}
			},
		)
	}
}
