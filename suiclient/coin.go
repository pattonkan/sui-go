package suiclient

import (
	"errors"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
)

type Coin struct {
	CoinType     sui.ObjectType    `json:"coinType"`
	CoinObjectId *sui.ObjectId     `json:"coinObjectId"`
	Version      *sui.BigInt       `json:"version"`
	Digest       *sui.ObjectDigest `json:"digest"`
	Balance      *sui.BigInt       `json:"balance"`

	LockedUntilEpoch    *sui.BigInt           `json:"lockedUntilEpoch,omitempty"`
	PreviousTransaction sui.TransactionDigest `json:"previousTransaction"`
}

type CoinPage = Page[*Coin, sui.ObjectId]

func (c *Coin) Ref() *sui.ObjectRef {
	return &sui.ObjectRef{
		Digest:   c.Digest,
		Version:  c.Version.Uint64(),
		ObjectId: c.CoinObjectId,
	}
}

func (c *Coin) IsSUI() bool {
	return c.CoinType == sui.SuiCoinType
}

type CoinFields struct {
	Balance *sui.BigInt
	Id      struct {
		Id *sui.ObjectId
	}
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

func (cs Coins) CoinRefs() []*sui.ObjectRef {
	coinRefs := make([]*sui.ObjectRef, len(cs))
	for idx, coin := range cs {
		coinRefs[idx] = coin.Ref()
	}
	return coinRefs
}

func (cs Coins) ObjectIds() []*sui.ObjectId {
	coinIds := make([]*sui.ObjectId, len(cs))
	for idx, coin := range cs {
		coinIds[idx] = coin.CoinObjectId
	}
	return coinIds
}

func (cs Coins) ObjectIdVals() []sui.ObjectId {
	coinIds := make([]sui.ObjectId, len(cs))
	for idx, coin := range cs {
		coinIds[idx] = *coin.CoinObjectId
	}
	return coinIds
}

func (cs Coins) CoinIds() []*sui.Address {
	panic("TODO")
}
