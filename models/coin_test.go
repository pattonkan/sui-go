package models_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func balanceObject(val uint64) models.SafeSuiBigInt[uint64] {
	return models.NewSafeSuiBigInt(val)
}

func TestCoins_PickSUICoinsWithGas(t *testing.T) {
	// coins 1,2,3,4,5
	testCoins := models.Coins{
		{Balance: balanceObject(3)},
		{Balance: balanceObject(5)},
		{Balance: balanceObject(1)},
		{Balance: balanceObject(4)},
		{Balance: balanceObject(2)},
	}
	type args struct {
		amount     *big.Int
		gasAmount  uint64
		pickMethod int
	}
	tests := []struct {
		name    string
		cs      models.Coins
		args    args
		want    models.Coins
		want1   *models.Coin
		wantErr bool
	}{
		{
			name: "case success 1",
			cs:   testCoins,
			args: args{
				amount:     new(big.Int),
				gasAmount:  0,
				pickMethod: models.PickSmaller,
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "case success 2",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(1),
				gasAmount:  2,
				pickMethod: models.PickSmaller,
			},
			want:    models.Coins{{Balance: balanceObject(1)}},
			want1:   &models.Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case success 3",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(4),
				gasAmount:  2,
				pickMethod: models.PickSmaller,
			},
			want:    models.Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(3)}},
			want1:   &models.Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case success 4",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(6),
				gasAmount:  2,
				pickMethod: models.PickSmaller,
			},
			want:    models.Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(3)}, {Balance: balanceObject(4)}},
			want1:   &models.Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case error 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(6),
				gasAmount:  6,
				pickMethod: models.PickSmaller,
			},
			want:    models.Coins{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "case error 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(100),
				gasAmount:  3,
				pickMethod: models.PickSmaller,
			},
			want:    models.Coins{},
			want1:   &models.Coin{Balance: balanceObject(3)},
			wantErr: true,
		},
		{
			name: "case bigger 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(3),
				gasAmount:  3,
				pickMethod: models.PickBigger,
			},
			want:    models.Coins{{Balance: balanceObject(5)}},
			want1:   &models.Coin{Balance: balanceObject(3)},
			wantErr: false,
		},
		{
			name: "case order 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(3),
				gasAmount:  3,
				pickMethod: models.PickByOrder,
			},
			want:    models.Coins{{Balance: balanceObject(5)}},
			want1:   &models.Coin{Balance: balanceObject(3)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, got1, err := tt.cs.PickSUICoinsWithGas(tt.args.amount, tt.args.gasAmount, tt.args.pickMethod)
				if (err != nil) != tt.wantErr {
					t.Errorf("Coins.PickSUICoinsWithGas() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != 0 && len(tt.want) != 0 {
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("Coins.PickSUICoinsWithGas() got: %v, want %v", got, tt.want)
					}
				}
				if !reflect.DeepEqual(got1, tt.want1) {
					t.Errorf("Coins.PickSUICoinsWithGas() got1: %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}

func TestCoins_PickCoins(t *testing.T) {
	// coins 1,2,3,4,5
	testCoins := models.Coins{
		{Balance: balanceObject(3)},
		{Balance: balanceObject(5)},
		{Balance: balanceObject(1)},
		{Balance: balanceObject(4)},
		{Balance: balanceObject(2)},
	}
	type args struct {
		amount     *big.Int
		pickMethod int
	}
	tests := []struct {
		name    string
		cs      models.Coins
		args    args
		want    models.Coins
		wantErr bool
	}{
		{
			name:    "smaller 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(2), pickMethod: models.PickSmaller},
			want:    models.Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(2)}},
			wantErr: false,
		},
		{
			name:    "smaller 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(4), pickMethod: models.PickSmaller},
			want:    models.Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(2)}, {Balance: balanceObject(3)}},
			wantErr: false,
		},
		{
			name:    "bigger 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(2), pickMethod: models.PickBigger},
			want:    models.Coins{{Balance: balanceObject(5)}},
			wantErr: false,
		},
		{
			name:    "bigger 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: models.PickBigger},
			want:    models.Coins{{Balance: balanceObject(5)}, {Balance: balanceObject(4)}},
			wantErr: false,
		},
		{
			name:    "pick by order 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: models.PickByOrder},
			want:    models.Coins{{Balance: balanceObject(3)}, {Balance: balanceObject(5)}},
			wantErr: false,
		},
		{
			name:    "pick by order 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(15), pickMethod: models.PickByOrder},
			want:    testCoins,
			wantErr: false,
		},
		{
			name:    "pick error",
			cs:      testCoins,
			args:    args{amount: big.NewInt(16), pickMethod: models.PickByOrder},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.cs.PickCoins(tt.args.amount, tt.args.pickMethod)
				if (err != nil) != tt.wantErr {
					t.Errorf("Coins.PickCoins() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Coins.PickCoins(): %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPickupCoins(t *testing.T) {
	coin := func(n uint64) models.Coin {
		return models.Coin{Balance: balanceObject(n), CoinType: models.SuiCoinType}
	}

	type args struct {
		inputCoins   *models.CoinPage
		targetAmount big.Int
		gasBudget    uint64
		limit        int
		moreCount    int
	}
	tests := []struct {
		name    string
		args    args
		want    *models.PickedCoins
		wantErr error
	}{
		{
			name: "moreCount = 3",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{
						coin(1e3), coin(1e5), coin(1e2), coin(1e4),
					},
				},
				targetAmount: *big.NewInt(1e3),
				moreCount:    3,
			},
			want: &models.PickedCoins{
				Coins: []models.Coin{
					coin(1e3), coin(1e5), coin(1e2),
				},
				TotalAmount:  *big.NewInt(1e3 + 1e5 + 1e2),
				TargetAmount: *big.NewInt(1e3),
			},
		},
		{
			name: "large gas",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{
						coin(1e3), coin(1e5), coin(1e2), coin(1e4),
					},
				},
				targetAmount: *big.NewInt(1e3),
				gasBudget:    1e9,
				moreCount:    3,
			},
			want: &models.PickedCoins{
				Coins: []models.Coin{
					coin(1e3), coin(1e5), coin(1e2), coin(1e4),
				},
				TotalAmount:  *big.NewInt(1e3 + 1e5 + 1e2 + 1e4),
				TargetAmount: *big.NewInt(1e3),
			},
		},
		{
			name: "ErrNoCoinsFound",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{},
				},
				targetAmount: *big.NewInt(101000),
			},
			wantErr: models.ErrNoCoinsFound,
		},
		{
			name: "ErrInsufficientBalance",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{
						coin(1e5), coin(1e6), coin(1e4),
					},
				},
				targetAmount: *big.NewInt(1e9),
			},
			wantErr: models.ErrInsufficientBalance,
		},
		{
			name: "ErrNeedMergeCoin 1",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{
						coin(1e5), coin(1e6), coin(1e4),
					},
					HasNextPage: true,
				},
				targetAmount: *big.NewInt(1e9),
			},
			wantErr: models.ErrNeedMergeCoin,
		},
		{
			name: "ErrNeedMergeCoin 2",
			args: args{
				inputCoins: &models.Page[models.Coin, sui_types.ObjectID]{
					Data: []models.Coin{
						coin(1e5), coin(1e6), coin(1e4), coin(1e5),
					},
					HasNextPage: false,
				},
				targetAmount: *big.NewInt(1e6 + 1e5*2 + 1e3),
				limit:        3,
			},
			wantErr: models.ErrNeedMergeCoin,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := models.PickupCoins(
					tt.args.inputCoins,
					tt.args.targetAmount,
					tt.args.gasBudget,
					tt.args.limit,
					tt.args.moreCount,
				)
				if tt.wantErr != nil {
					require.Equal(t, err, tt.wantErr)
				} else {
					require.Equal(t, got, tt.want)
					t.Log("suggest max gas budget: ", tt.want.SuggestMaxGasBudget())
				}
			},
		)
	}
}
