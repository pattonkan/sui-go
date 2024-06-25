package models

import "github.com/howjmay/sui-go/sui_types"

type Balance struct {
	CoinType        sui_types.ObjectType                `json:"coinType"`
	CoinObjectCount uint64                              `json:"coinObjectCount"`
	TotalBalance    SuiBigInt                           `json:"totalBalance"`
	LockedBalance   map[SafeSuiBigInt[uint64]]SuiBigInt `json:"lockedBalance"`
}
