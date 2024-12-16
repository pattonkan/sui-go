package suiclient

// FIXME refactor this to better impl

import (
	"errors"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
)

var (
	ErrNoCoinsFound        = errors.New("no coins found")
	ErrInsufficientBalance = errors.New("insufficient account balance")

	ErrNeedMergeCoin    = errors.New("no coins of such a large amount were found to execute this transaction")
	ErrNeedSplitGasCoin = errors.New("missing an extra coin to use as the transaction fee")

	ErrCoinsNotMatchRequest = errors.New("coins not match request")
	ErrCoinsNeedMoreObject  = errors.New("you should get more SUI coins and try again")
)

const MAX_INPUT_COUNT_MERGE = 256 - 1 // TODO find reference in Sui monorepo repo

type PickedCoins struct {
	Coins        Coins
	TotalAmount  *big.Int
	TargetAmount *big.Int
}

func (p *PickedCoins) Count() int {
	return len(p.Coins)
}

func (p *PickedCoins) CoinIds() []*sui.ObjectId {
	coinIDs := make([]*sui.ObjectId, len(p.Coins))
	for idx, coin := range p.Coins {
		coinIDs[idx] = coin.CoinObjectId
	}
	return coinIDs
}

func (p *PickedCoins) CoinRefs() []*sui.ObjectRef {
	coinRefs := make([]*sui.ObjectRef, len(p.Coins))
	for idx, coin := range p.Coins {
		coinRefs[idx] = coin.Ref()
	}
	return coinRefs
}

func PickupCoins(inputCoins *CoinPage,
	targetAmount *big.Int,
	gasBudget uint64,
	maxCoinNum int,
	minCoinNum int,
) (*PickedCoins, error) {
	coins := inputCoins.Data
	inputCount := len(coins)
	if inputCount <= 0 {
		return nil, ErrNoCoinsFound
	}
	if maxCoinNum <= 0 {
		maxCoinNum = MAX_INPUT_COUNT_MERGE
	}
	if minCoinNum <= 0 {
		minCoinNum = 3
	}
	if minCoinNum > maxCoinNum {
		minCoinNum = maxCoinNum
	}
	totalTarget := new(big.Int).Add(targetAmount, new(big.Int).SetUint64(gasBudget))

	total := big.NewInt(0)
	pickedCoins := []*Coin{}
	for i, coin := range coins {
		total = total.Add(total, new(big.Int).SetUint64(coin.Balance.Uint64()))
		pickedCoins = append(pickedCoins, coin)
		if i+1 > maxCoinNum {
			return nil, ErrNeedMergeCoin
		}
		if i+1 < minCoinNum {
			continue
		}
		if total.Cmp(totalTarget) >= 0 {
			break
		}
	}
	if total.Cmp(totalTarget) < 0 {
		if inputCoins.HasNextPage {
			return nil, ErrNeedMergeCoin
		}
		sub := new(big.Int).Sub(totalTarget, total)
		if sub.Uint64() > gasBudget {
			return nil, ErrInsufficientBalance
		}
	}
	return &PickedCoins{
		Coins:        pickedCoins,
		TargetAmount: targetAmount,
		TotalAmount:  total,
	}, nil
}
