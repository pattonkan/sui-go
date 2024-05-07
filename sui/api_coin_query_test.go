package sui_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetAllBalances(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	balances, err := api.GetAllBalances(context.TODO(), sui_signer.TEST_ADDRESS)
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
		address *sui_types.SuiAddress
		cursor  *sui_types.ObjectID
		limit   uint
	}

	tests := []struct {
		name    string
		a       *sui.ImplSuiAPI
		args    args
		want    *models.CoinPage
		wantErr bool
	}{
		{
			name: "test case 1",
			a:    sui.NewSuiClient(conn.DevnetEndpointUrl),
			args: args{
				ctx:     context.TODO(),
				address: sui_signer.TEST_ADDRESS,
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

func TestGetBalance(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	balance, err := api.GetBalance(context.TODO(), sui_signer.TEST_ADDRESS, "")
	require.NoError(t, err)
	t.Logf(
		"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
		balance.CoinType, balance.CoinObjectCount,
		balance.TotalBalance.String(), balance.LockedBalance,
	)
}

func TestGetCoinMetadata(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	metadata, err := api.GetCoinMetadata(context.TODO(), models.SuiCoinType)
	require.NoError(t, err)
	t.Logf("%#v", metadata)
}

func TestGetCoins(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	defaultCoinType := models.SuiCoinType
	coins, err := api.GetCoins(context.TODO(), sui_signer.TEST_ADDRESS, &defaultCoinType, nil, 1)
	require.NoError(t, err)
	t.Logf("%#v", coins)
	require.GreaterOrEqual(t, len(coins.Data), 0)
	require.Equal(t, models.SuiCoinType, coins.Data[0].CoinType)
	require.Greater(t, coins.Data[0].Balance.Int64(), int64(0))
}

func TestGetTotalSupply(t *testing.T) {
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
			api:  sui.NewSuiClient(conn.DevnetEndpointUrl),
			args: args{
				context.TODO(),
				models.SuiCoinType,
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
