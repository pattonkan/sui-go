package sui_types

import (
	"github.com/howjmay/sui-go/lib"
)

type PackageID = SuiAddress
type ObjectID = SuiAddress
type SequenceNumber = uint64

func PackageIDFromHex(str string) (*PackageID, error) {
	return SuiAddressFromHex(str)
}

func MustPackageIDFromHex(str string) *PackageID {
	packageID, err := SuiAddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return packageID
}

func ObjectIDFromHex(str string) (*ObjectID, error) {
	return SuiAddressFromHex(str)
}

func MustObjectIDFromHex(str string) *ObjectID {
	objectID, err := SuiAddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return objectID
}

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectID ObjectID       `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   ObjectDigest   `json:"digest"`
}

type MoveObjectType struct {
	Other     *StructTag
	GasCoin   *lib.EmptyEnum
	StakedSui *lib.EmptyEnum
	Coin      *TypeTag
}

func (o MoveObjectType) IsBcsEnum() {}
