package suiptb

import (
	"github.com/pattonkan/sui-go/sui"
)

// https://sdk.mystenlabs.com/typescript/transaction-building/basics#object-references
// https://docs.sui.io/concepts/transactions/prog-txn-blocks
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
	// The gas coin. The gas coin can only be used by-ref, except for with
	// `TransferObjects`, which can use it by-value.
	GasCoin *sui.EmptyEnum
	// One of the input objects or primitive values (from `ProgrammableTransaction` inputs)
	Input *uint16
	// The result of another transaction (from `ProgrammableTransaction` transactions)
	Result *uint16
	// Like a `Result` but it accesses a nested result. Currently, the only usage of this is to access a
	// value from a Move call with multiple return values.
	NestedResult *NestedResult
}
type NestedResult struct {
	Cmd    uint16 // command index
	Result uint16 // result index
}

func (a Argument) IsBcsEnum() {}

type ProgrammableMoveCall struct {
	Package       *sui.PackageId
	Module        sui.Identifier
	Function      sui.Identifier
	TypeArguments []sui.TypeTag
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
	Dependencies []*sui.ObjectId
}

type ProgrammableMakeMoveVec struct {
	Type    *sui.TypeTag `bcs:"optional"`
	Objects []Argument
}

type ProgrammableUpgrade struct {
	Modules      [][]byte
	Dependencies []*sui.ObjectId
	PackageId    *sui.PackageId
	Ticket       Argument
}
