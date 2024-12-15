package movebcs

import (
	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui"
)

const (
	// A type passed to create_supply is not a one-time witness.
	ECoinBadWitness = 0
	// Invalid arguments are passed to a function.
	ECoinInvalidArg = 1
	// Trying to split a coin more times than its balance allows.
	ECoinNotEnough = 2
	// #[error]
	// const EGlobalPauseNotAllowed: vector<u8> =
	//    b"Kill switch was not allowed at the creation of the DenyCapV2";
	ECoinGlobalPauseNotAllowed = 3
)

type MoveCoin struct {
	Id      *sui.ObjectId
	Balance uint64
}

type MoveCoinMetadata struct {
	Id          *sui.ObjectId
	Decimals    uint8
	Name        string
	Symbol      string
	Description string
	IconUrl     bcs.Option[string]
}

type MoveRegulatedCoinMetadata struct {
	Id                 *sui.ObjectId
	CoinMetadataObject *sui.ObjectId
	DenyCapObject      *sui.ObjectId
}

type MoveTreasuryCap struct {
	Id          *sui.ObjectId
	TotalSupply *MoveSupply
}

type MoveDenyCapV2 struct {
	Id               *sui.ObjectId
	AllowGlobalPause bool
}
