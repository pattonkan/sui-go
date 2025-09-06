package suiclient

import (
	"encoding/json"
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

type CoinPage = Page[*Coin, string]

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

func (c Coin) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(b)
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

func (cs Coins) String() string {
	b, err := json.Marshal(cs)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (cs Coins) CoinIds() []*sui.Address {
	panic("TODO")
}
