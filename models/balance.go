package models

type Balance struct {
	CoinType        string             `json:"coinType"`
	CoinObjectCount uint64             `json:"coinObjectCount"`
	TotalBalance    *BigInt            `json:"totalBalance"`
	LockedBalance   map[uint64]*BigInt `json:"lockedBalance"`
}
