package suiclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/suisigner/suicrypto"

	"github.com/stretchr/testify/require"
)

func TestGetDynamicFieldObject(t *testing.T) {
	t.Skip("FIXME")
	api := suiclient.NewClient(conn.TestnetEndpointUrl)
	parentObjectId, err := sui.AddressFromHex("0x1719957d7a2bf9d72459ff0eab8e600cbb1991ef41ddd5b4a8c531035933d256")
	require.NoError(t, err)
	type args struct {
		ctx            context.Context
		parentObjectId *sui.ObjectId
		name           *suiclient.DynamicFieldName
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.SuiObjectResponse
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				ctx:            context.TODO(),
				parentObjectId: parentObjectId,
				name: &suiclient.DynamicFieldName{
					Type:  "address",
					Value: "0xf9ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.GetDynamicFieldObject(tt.args.ctx, &suiclient.GetDynamicFieldObjectRequest{
					ParentObjectId: tt.args.parentObjectId,
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
	client := suiclient.NewClient(conn.MainnetEndpointUrl)
	limit := 5
	type args struct {
		ctx            context.Context
		parentObjectId *sui.ObjectId
		cursor         *sui.ObjectId
		limit          *uint
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.DynamicFieldPage
		wantErr error
	}{
		{
			name: "a deepbook shared object",
			args: args{
				ctx:            context.TODO(),
				parentObjectId: sui.MustAddressFromHex("0xa9d09452bba939b3172c0242d022274845cfe4e58648b73dd33b3d5b823dc8ae"),
				cursor:         nil,
				limit:          func() *uint { tmpLimit := uint(limit); return &tmpLimit }(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := client.GetDynamicFields(tt.args.ctx, &suiclient.GetDynamicFieldsRequest{
					ParentObjectId: tt.args.parentObjectId,
					Cursor:         tt.args.cursor,
					Limit:          tt.args.limit,
				})
				require.ErrorIs(t, err, tt.wantErr)
				// object Id is '0x4405b50d791fd3346754e8171aaab6bc2ed26c2c46efdd033c14b30ae507ac33'
				// it has 'internal_nodes' field in type '0x2::table::Table<u64, 0xdee9::critbit::InternalNode'
				require.Len(t, got.Data, limit)
				for _, field := range got.Data {
					require.Equal(t, "u64", field.Name.Type)
					require.Equal(t, sui.MustNewTypeTag("0xdee9::critbit::InternalNode").String(), sui.MustNewTypeTag(field.ObjectType).String())
				}
			},
		)
	}
}

func TestGetOwnedObjects(t *testing.T) {
	t.Run("struct tag", func(t *testing.T) {
		client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, suicrypto.KeySchemeFlagDefault, 1)
		client.WithSignerAndFund(suisigner.TEST_SEED, suicrypto.KeySchemeFlagDefault, 1)
		structTag, err := sui.StructTagFromString("0x2::coin::Coin<0x2::sui::SUI>")
		require.NoError(t, err)
		query := suiclient.SuiObjectResponseQuery{
			Filter: &suiclient.SuiObjectDataFilter{
				StructType: structTag,
			},
			Options: &suiclient.SuiObjectDataOptions{
				ShowType:    true,
				ShowContent: true,
			},
		}
		limit := uint(5)
		objs, err := client.GetOwnedObjects(context.Background(), &suiclient.GetOwnedObjectsRequest{
			Address: signer.Address,
			Query:   &query,
			Limit:   &limit,
		})
		require.NoError(t, err)
		require.Equal(t, int(limit), len(objs.Data))
		require.NoError(t, err)
		var fields suiclient.CoinFields
		err = json.Unmarshal(objs.Data[len(objs.Data)-1].Data.Content.Data.MoveObject.Fields, &fields)
		require.NoError(t, err)
		require.Equal(t, "200000000000", fields.Balance.String())
	})

	t.Run("move module", func(t *testing.T) {
		client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, suicrypto.KeySchemeFlagDefault, 1)
		client.WithSignerAndFund(suisigner.TEST_SEED, suicrypto.KeySchemeFlagDefault, 1)
		query := suiclient.SuiObjectResponseQuery{
			Filter: &suiclient.SuiObjectDataFilter{
				AddressOwner: signer.Address,
			},
			Options: &suiclient.SuiObjectDataOptions{
				ShowType:    true,
				ShowContent: true,
			},
		}
		limit := uint(5)
		objs, err := client.GetOwnedObjects(context.Background(), &suiclient.GetOwnedObjectsRequest{
			Address: signer.Address,
			Query:   &query,
			Limit:   &limit,
		})
		require.NoError(t, err)
		require.Equal(t, int(limit), len(objs.Data))
		require.NoError(t, err)
		var fields suiclient.CoinFields
		err = json.Unmarshal(objs.Data[len(objs.Data)-1].Data.Content.Data.MoveObject.Fields, &fields)
		require.NoError(t, err)
		require.Equal(t, "200000000000", fields.Balance.String())
	})
}

func TestQueryEvents(t *testing.T) {
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	limit := 10

	type args struct {
		ctx             context.Context
		query           *suiclient.EventFilter
		cursor          *suiclient.EventId
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.EventPage
		wantErr error
	}{
		{
			name: "event in deepbook.batch_cancel_order()",
			args: args{
				ctx: context.TODO(),
				query: &suiclient.EventFilter{
					Sender: sui.MustAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
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
					&suiclient.QueryEventsRequest{
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
						sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
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
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           *suiclient.TransactionBlockResponseQuery
		cursor          *sui.TransactionDigest
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.TransactionBlocksPage
		wantErr bool
	}{
		{
			name: "test for queryTransactionBlocks",
			args: args{
				ctx: context.TODO(),
				query: &suiclient.TransactionBlockResponseQuery{
					Filter: &suiclient.TransactionFilter{
						FromAddress: suisigner.TEST_ADDRESS,
					},
					Options: &suiclient.SuiTransactionBlockResponseOptions{
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
					&suiclient.QueryTransactionBlocksRequest{
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
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	addr, err := api.ResolveNameServiceAddress(context.Background(), "example.sui")
	require.NoError(t, err)
	require.Equal(t, "0x214a4199264348df2364acd683a3971a9927a5252747f4e0776f0506922f9db0", addr.String())

	addr, err = api.ResolveNameServiceAddress(context.Background(), "2222.suijjzzww")
	require.Equal(t, "0x0", addr.ShortString())
	require.NoError(t, err)
}

func TestResolveNameServiceNames(t *testing.T) {
	api := suiclient.NewClient(conn.MainnetEndpointUrl)
	owner := sui.MustAddressFromHex("0x214a4199264348df2364acd683a3971a9927a5252747f4e0776f0506922f9db0")
	namePage, err := api.ResolveNameServiceNames(context.Background(), &suiclient.ResolveNameServiceNamesRequest{
		Owner: owner,
	})
	require.NoError(t, err)
	require.NotEmpty(t, namePage.Data)
	t.Log(namePage.Data)

	owner = sui.MustAddressFromHex("0x57188743983628b3474648d8aa4a9ee8abebe8f681")
	namePage, err = api.ResolveNameServiceNames(context.Background(), &suiclient.ResolveNameServiceNamesRequest{
		Owner: owner,
	})
	require.NoError(t, err)
	require.Empty(t, namePage.Data)
}

func TestSubscribeEvent(t *testing.T) {
	t.Skip("fixme: change to another endpoint")
	api := suiclient.NewSuiWebsocketClient(context.TODO(), "wss://sui-mainnet.public.blastapi.io")

	type args struct {
		ctx      context.Context
		filter   *suiclient.EventFilter
		resultCh chan suiclient.Event
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.EventPage
		wantErr bool
	}{
		{
			name: "test for filter events",
			args: args{
				ctx: context.TODO(),
				filter: &suiclient.EventFilter{
					Package: sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
				},
				resultCh: make(chan suiclient.Event),
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

func TestSubscribeTransaction(t *testing.T) {
	t.Skip("fixme: change to another endpoint")
	api := suiclient.NewSuiWebsocketClient(context.TODO(), "wss://sui-mainnet.public.blastapi.io")

	type args struct {
		ctx      context.Context
		filter   *suiclient.TransactionFilter
		resultCh chan suiclient.WrapperTaggedJson[suiclient.SuiTransactionBlockEffects]
	}
	tests := []struct {
		name    string
		args    args
		want    *suiclient.SuiTransactionBlockEffects
		wantErr bool
	}{
		{
			name: "test for filter transaction",
			args: args{
				ctx: context.TODO(),
				filter: &suiclient.TransactionFilter{
					MoveFunction: &suiclient.TransactionFilterMoveFunction{
						Package: sui.MustPackageIdFromHex("0x2c68443db9e8c813b194010c11040a3ce59f47e4eb97a2ec805371505dad7459"),
					},
				},
				resultCh: make(chan suiclient.WrapperTaggedJson[suiclient.SuiTransactionBlockEffects]),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := api.SubscribeTransaction(
					tt.args.ctx,
					tt.args.filter,
					tt.args.resultCh,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("SubscribeTransaction() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				cnt := 0
				for results := range tt.args.resultCh {
					fmt.Println("results: ", results.Data.V1)
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
