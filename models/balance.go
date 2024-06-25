package models

import (
	"github.com/howjmay/sui-go/sui_types"
)

type Balance struct {
	CoinType        sui_types.ObjectType `json:"coinType"`
	CoinObjectCount uint64               `json:"coinObjectCount"`
	TotalBalance    *BigInt              `json:"totalBalance"`
	LockedBalance   map[BigInt]*BigInt   `json:"lockedBalance"`
}
