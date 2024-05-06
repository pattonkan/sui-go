package sui_types

import (
	"github.com/howjmay/sui-go/lib"
	"github.com/howjmay/sui-go/move_types"
)

type SuiAddress = move_types.AccountAddress
type PackageID = move_types.AccountAddress
type ObjectID = move_types.AccountAddress
type SequenceNumber = uint64

func AddressFromHex(str string) (*SuiAddress, error) {
	return move_types.NewAccountAddressHex(str)
}

func MustAddressFromHex(str string) *SuiAddress {
	addr, err := move_types.NewAccountAddressHex(str)
	if err != nil {
		panic(err)
	}
	return addr
}

func PackageIDFromHex(str string) (*PackageID, error) {
	return move_types.NewAccountAddressHex(str)
}

func MustPackageIDFromHex(str string) *PackageID {
	packageID, err := move_types.NewAccountAddressHex(str)
	if err != nil {
		panic(err)
	}
	return packageID
}

func NewObjectIDFromHex(str string) (*ObjectID, error) {
	return move_types.NewAccountAddressHex(str)
}

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectID ObjectID       `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   ObjectDigest   `json:"digest"`
}

type MoveObjectType struct {
	Other     *move_types.StructTag
	GasCoin   *lib.EmptyEnum
	StakedSui *lib.EmptyEnum
	Coin      *move_types.TypeTag
}

func (o MoveObjectType) IsBcsEnum() {}
