package suiclient_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suiclient/conn"
	"github.com/howjmay/sui-go/suisigner"
	"github.com/stretchr/testify/require"
)

func TestGetAllBalances(t *testing.T) {
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	balances, err := api.GetAllBalances(context.TODO(), suisigner.TEST_ADDRESS)
	require.NoError(t, err)
	for _, balance := range balances {
		t.Logf(
			"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
			balance.CoinType, balance.CoinObjectCount,
			balance.TotalBalance.String(), balance.LockedBalance,
		)
	}
}

func TestGetAllCoins(t *testing.T) {
	type args struct {
		ctx     context.Context
		address *sui.Address
		cursor  *sui.ObjectId
		limit   uint
	}

	tests := []struct {
		name    string
		a       *suiclient.ClientImpl
		args    args
		want    *suiclient.CoinPage
		wantErr bool
	}{
		{
			name: "successful with limit",
			a:    suiclient.NewClient(conn.TestnetEndpointUrl),
			args: args{
				ctx:     context.TODO(),
				address: suisigner.TEST_ADDRESS,
				cursor:  nil,
				limit:   3,
			},
			wantErr: false,
		},
		{
			name: "successful without limit",
			a:    suiclient.NewClient(conn.TestnetEndpointUrl),
			args: args{
				ctx:     context.TODO(),
				address: suisigner.TEST_ADDRESS,
				cursor:  nil,
				limit:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.a.GetAllCoins(tt.args.ctx, &suiclient.GetAllCoinsRequest{
					Owner:  tt.args.address,
					Cursor: tt.args.cursor,
					Limit:  tt.args.limit,
				})
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAllCoins() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				// we have called multiple times RequestFundFromFaucet() on testnet, so the account have several SUI objects.
				require.GreaterOrEqual(t, len(got.Data), int(tt.args.limit))
				require.NotNil(t, got.NextCursor)
			},
		)
	}
}

func TestGetBalance(t *testing.T) {
	api := suiclient.NewClient(conn.DevnetEndpointUrl)
	balance, err := api.GetBalance(context.TODO(), &suiclient.GetBalanceRequest{Owner: suisigner.TEST_ADDRESS})
	require.NoError(t, err)
	t.Logf(
		"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
		balance.CoinType, balance.CoinObjectCount,
		balance.TotalBalance.String(), balance.LockedBalance,
	)
}

func TestGetCoinMetadata(t *testing.T) {
	api := suiclient.NewClient(conn.TestnetEndpointUrl)
	metadata, err := api.GetCoinMetadata(context.TODO(), sui.SuiCoinType)
	require.NoError(t, err)
	testSuiMetadata := &suiclient.CoinMetadata{
		Decimals:    9,
		Description: "",
		IconUrl:     "",
		Id:          sui.MustObjectIdFromHex("0x587c29de216efd4219573e08a1f6964d4fa7cb714518c2c8a0f29abfa264327d"),
		Name:        "Sui",
		Symbol:      "SUI",
	}
	require.Equal(t, testSuiMetadata, metadata)
}

func TestGetCoins(t *testing.T) {
	api := suiclient.NewClient(conn.TestnetEndpointUrl)
	defaultCoinType := sui.SuiCoinType
	coins, err := api.GetCoins(context.TODO(), &suiclient.GetCoinsRequest{
		Owner:    suisigner.TEST_ADDRESS,
		CoinType: &defaultCoinType,
		Limit:    3,
	})
	require.NoError(t, err)

	require.Len(t, coins.Data, 3)
	for _, data := range coins.Data {
		require.Equal(t, sui.SuiCoinType, data.CoinType)
		require.Greater(t, data.Balance.Int64(), int64(0))
	}
}

func TestGetTotalSupply(t *testing.T) {
	type args struct {
		ctx      context.Context
		coinType string
	}

	tests := []struct {
		name    string
		api     *suiclient.ClientImpl
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "get Sui supply",
			api:  suiclient.NewClient(conn.DevnetEndpointUrl),
			args: args{
				context.TODO(),
				sui.SuiCoinType,
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
				targetSupply := &suiclient.Supply{Value: sui.NewBigInt(10000000000000000000)}
				require.Equal(t, targetSupply, got)
			},
		)
	}
}
