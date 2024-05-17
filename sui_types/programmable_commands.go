package sui_types

import "github.com/howjmay/sui-go/sui_types/serialization"

type Command struct {
	MoveCall        *ProgrammableMoveCall
	TransferObjects *ProgrammableTransferObjects
	SplitCoins      *ProgrammableSplitCoins
	MergeCoins      *ProgrammableMergeCoins
	// `Publish` publishes a Move package. Returns the upgrade capability object.
	Publish *ProgrammablePublish
	// `MakeMoveVec` constructs a vector of objects that can be passed into a moveCall.
	// This is required as there’s no way to define a vector as an input.
	MakeMoveVec *ProgrammableMakeMoveVec
	// upgrades a Move package
	Upgrade *ProgrammableUpgrade
}

func (c Command) IsBcsEnum() {}

type Argument struct {
	GasCoin      *serialization.EmptyEnum
	Input        *uint16
	Result       *uint16
	NestedResult *struct {
		Result1 uint16
		Result2 uint16
	}
}

func (a Argument) IsBcsEnum() {}

type ProgrammableMoveCall struct {
	Package       *ObjectID
	Module        string
	Function      string
	TypeArguments []TypeTag
	Arguments     []Argument
}

type ProgrammableTransferObjects struct {
	Objects []Argument
	Address Argument
}

type ProgrammableSplitCoins struct {
	Coin    Argument
	Amounts []Argument
}

type ProgrammableMergeCoins struct {
	Destination Argument
	Sources     []Argument
}

type ProgrammablePublish struct {
	Modules      [][]byte
	Dependencies []ObjectID
}

type ProgrammableMakeMoveVec struct {
	Type    *TypeTag `bcs:"optional"`
	Objects []Argument
}

type ProgrammableUpgrade struct {
	Modules      [][]byte
	Dependencies []ObjectID
	PackageId    ObjectID
	Ticket       Argument
}
