package sui

type EpochId = uint64

type PackageId = Address
type ObjectId = Address

type SequenceNumber = uint64

type Identifier = string
type ObjectType = string

func PackageIdFromHex(str string) (*PackageId, error) {
	return AddressFromHex(str)
}

func MustPackageIdFromHex(str string) *PackageId {
	packageId, err := AddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return packageId
}

func ObjectIdFromHex(str string) (*ObjectId, error) {
	return AddressFromHex(str)
}

func MustObjectIdFromHex(str string) *ObjectId {
	objectId, err := AddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return objectId
}

func ObjectTypeFromString(str string) ObjectType {
	return str
}

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectId *ObjectId      `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   *ObjectDigest  `json:"digest"`
}

type MoveObjectType struct {
	Other     *StructTag
	GasCoin   *EmptyEnum
	StakedSui *EmptyEnum
	Coin      *TypeTag
}

func (o MoveObjectType) IsBcsEnum() {}
