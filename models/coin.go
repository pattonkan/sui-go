package models

import (
	"errors"
	"math/big"
	"sort"

	"github.com/howjmay/sui-go/sui_types"
)

type Coin struct {
	CoinType     string                  `json:"coinType"`
	CoinObjectID *sui_types.ObjectID     `json:"coinObjectID"`
	Version      SafeSuiBigInt[uint64]   `json:"version"`
	Digest       *sui_types.ObjectDigest `json:"digest"`
	Balance      SafeSuiBigInt[uint64]   `json:"balance"`

	LockedUntilEpoch    *SafeSuiBigInt[uint64]      `json:"lockedUntilEpoch,omitempty"`
	PreviousTransaction sui_types.TransactionDigest `json:"previousTransaction"`
}

type CoinPage = Page[*Coin, sui_types.ObjectID]

func (c *Coin) Ref() *sui_types.ObjectRef {
	return &sui_types.ObjectRef{
		Digest:   c.Digest,
		Version:  c.Version.data,
		ObjectID: c.CoinObjectID,
	}
}

func (c *Coin) IsSUI() bool {
	return c.CoinType == SuiCoinType
}

type Coins []*Coin

func (cs Coins) TotalBalance() *big.Int {
	total := new(big.Int)
	for _, coin := range cs {
		total = total.Add(total, new(big.Int).SetUint64(coin.Balance.Uint64()))
	}
	return total
}

func (cs Coins) PickCoinNoLess(amount uint64) (*Coin, error) {
	for i, coin := range cs {
		if coin.Balance.Uint64() >= amount {
			cs = append(cs[:i], cs[i+1:]...)
			return coin, nil
		}
	}
	if len(cs) <= 3 {
		return nil, errors.New("insufficient balance")
	}
	return nil, errors.New("no coin is enough to cover the gas")
}

func (cs Coins) CoinRefs() []*sui_types.ObjectRef {
	coinRefs := make([]*sui_types.ObjectRef, len(cs))
	for idx, coin := range cs {
		coinRefs[idx] = coin.Ref()
	}
	return coinRefs
}

func (cs Coins) ObjectIDs() []*sui_types.ObjectID {
	coinIDs := make([]*sui_types.ObjectID, len(cs))
	for idx, coin := range cs {
		coinIDs[idx] = coin.CoinObjectID
	}
	return coinIDs
}

func (cs Coins) ObjectIDVals() []sui_types.ObjectID {
	coinIDs := make([]sui_types.ObjectID, len(cs))
	for idx, coin := range cs {
		coinIDs[idx] = *coin.CoinObjectID
	}
	return coinIDs
}

const (
	PickMethodSmaller = iota // pick smaller coins to match amount
	PickMethodBigger         // pick bigger coins to match amount
	PickMethodByOrder        // pick coins by coins order to match amount
)

// PickSUICoinsWithGas pick coins, which sum >= amount, and pick a gas coin >= gasAmount which not in coins
// if not satisfied amount/gasAmount, an ErrCoinsNotMatchRequest/ErrCoinsNeedMoreObject error will return
// if gasAmount == 0, a nil gasCoin will return
// pickMethod, see PickMethodSmaller|PickMethodBigger|PickMethodByOrder
func (cs Coins) PickSUICoinsWithGas(amount *big.Int, gasAmount uint64, pickMethod int) (Coins, *Coin, error) {
	if gasAmount == 0 {
		res, err := cs.PickCoins(amount, pickMethod)
		return res, nil, err
	}

	if amount.Cmp(new(big.Int)) == 0 && gasAmount == 0 {
		return make(Coins, 0), nil, nil
	} else if len(cs) == 0 {
		return cs, nil, ErrCoinsNeedMoreObject
	}

	// find smallest to match gasAmount
	var gasCoin *Coin
	var selectIndex int
	for i := range cs {
		if cs[i].Balance.Uint64() < gasAmount {
			continue
		}

		if nil == gasCoin || gasCoin.Balance.Uint64() > cs[i].Balance.Uint64() {
			gasCoin = cs[i]
			selectIndex = i
		}
	}
	if nil == gasCoin {
		return nil, nil, ErrCoinsNotMatchRequest
	}

	lastCoins := make(Coins, 0, len(cs)-1)
	lastCoins = append(lastCoins, cs[0:selectIndex]...)
	lastCoins = append(lastCoins, cs[selectIndex+1:]...)
	pickCoins, err := lastCoins.PickCoins(amount, pickMethod)
	return pickCoins, gasCoin, err
}

// PickCoins pick coins, which sum >= amount,
// pickMethod, see PickMethodSmaller|PickMethodBigger|PickMethodByOrder
// if not satisfied amount, an ErrCoinsNeedMoreObject error will return
func (cs Coins) PickCoins(amount *big.Int, pickMethod int) (Coins, error) {
	var sortedCoins Coins
	if pickMethod == PickMethodByOrder {
		sortedCoins = cs
	} else {
		sortedCoins = make(Coins, len(cs))
		copy(sortedCoins, cs)
		sort.Slice(
			sortedCoins, func(i, j int) bool {
				if pickMethod == PickMethodSmaller {
					return sortedCoins[i].Balance.Uint64() < sortedCoins[j].Balance.Uint64()
				} else {
					return sortedCoins[i].Balance.Uint64() >= sortedCoins[j].Balance.Uint64()
				}
			},
		)
	}

	result := make(Coins, 0)
	total := new(big.Int)
	for _, coin := range sortedCoins {
		result = append(result, coin)
		total = new(big.Int).Add(total, new(big.Int).SetUint64(coin.Balance.Uint64()))
		if total.Cmp(amount) >= 0 {
			return result, nil
		}
	}

	return nil, ErrCoinsNeedMoreObject
}
