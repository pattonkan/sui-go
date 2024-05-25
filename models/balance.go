package models

type Balance struct {
	CoinType        string                              `json:"coinType"`
	CoinObjectCount uint64                              `json:"coinObjectCount"`
	TotalBalance    SuiBigInt                           `json:"totalBalance"`
	LockedBalance   map[SafeSuiBigInt[uint64]]SuiBigInt `json:"lockedBalance"`
}
