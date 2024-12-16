package suiclient

import "github.com/pattonkan/sui-go/sui"

type Balance struct {
	CoinType        sui.ObjectType             `json:"coinType"`
	CoinObjectCount uint64                     `json:"coinObjectCount"`
	TotalBalance    *sui.BigInt                `json:"totalBalance"`
	LockedBalance   map[sui.BigInt]*sui.BigInt `json:"lockedBalance"`
}

type CoinMetadata struct {
	Decimals    uint8         `json:"decimals"`
	Name        string        `json:"name"`
	Symbol      string        `json:"symbol"`
	Description string        `json:"description"`
	IconUrl     string        `json:"iconUrl,omitempty"`
	Id          *sui.ObjectId `json:"id"`
}

type Supply struct {
	Value *sui.BigInt `json:"value"`
}
