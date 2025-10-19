package suiptb

import (
	"fmt"
	"strings"

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

func (c Command) String() string {
	if c.MoveCall != nil {
		return c.MoveCall.String()
	}
	if c.TransferObjects != nil {
		return c.TransferObjects.String()
	}
	if c.SplitCoins != nil {
		return c.SplitCoins.String()
	}
	if c.MergeCoins != nil {
		return c.MergeCoins.String()
	}
	if c.Publish != nil {
		return c.Publish.String()
	}
	if c.MakeMoveVec != nil {
		return c.MakeMoveVec.String()
	}
	if c.Upgrade != nil {
		return c.Upgrade.String()
	}
	panic("invalid command")
}

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

func (a Argument) String() string {
	if a.GasCoin != nil {
		return "GasCoin"
	}
	if a.Input != nil {
		return fmt.Sprintf("Input(%d)", *a.Input)
	}
	if a.Result != nil {
		return fmt.Sprintf("Result(%d)", *a.Result)
	}
	if a.NestedResult != nil {
		return fmt.Sprintf("NestedResult(%d, %d)", a.NestedResult.Cmd, a.NestedResult.Result)
	}
	panic("invalid argument")
}

func GetArgumentGasCoin() Argument {
	return Argument{GasCoin: &sui.EmptyEnum{}}
}

type ProgrammableMoveCall struct {
	Package       *sui.PackageId
	Module        sui.Identifier
	Function      sui.Identifier
	TypeArguments []sui.TypeTag
	Arguments     []Argument
}

func (p *ProgrammableMoveCall) String() string {
	typeArgs := make([]string, len(p.TypeArguments))
	for i, t := range p.TypeArguments {
		typeArgs[i] = t.String()
	}

	args := make([]string, len(p.Arguments))
	for i, a := range p.Arguments {
		args[i] = a.String()
	}

	return fmt.Sprintf(
		"MoveCall: %s::%s<%s>(%s)",
		p.Module,
		p.Function,
		strings.Join(typeArgs, ", "),
		strings.Join(args, ", "),
	)
}

type ProgrammableTransferObjects struct {
	Objects []Argument
	Address Argument
}

func (p *ProgrammableTransferObjects) String() string {
	objects := make([]string, len(p.Objects))
	for i, obj := range p.Objects {
		objects[i] = obj.String()
	}
	return fmt.Sprintf("TransferObjects([%s], %s)", strings.Join(objects, ", "), p.Address.String())
}

type ProgrammableSplitCoins struct {
	Coin    Argument
	Amounts []Argument
}

func (p *ProgrammableSplitCoins) String() string {
	amounts := make([]string, len(p.Amounts))
	for i, amt := range p.Amounts {
		amounts[i] = amt.String()
	}
	return fmt.Sprintf("SplitCoins(%s, [%s])", p.Coin.String(), strings.Join(amounts, ", "))
}

type ProgrammableMergeCoins struct {
	Destination Argument
	Sources     []Argument
}

func (p *ProgrammableMergeCoins) String() string {
	sources := make([]string, len(p.Sources))
	for i, src := range p.Sources {
		sources[i] = src.String()
	}
	return fmt.Sprintf("MergeCoins(%s, [%s])", p.Destination.String(), strings.Join(sources, ", "))
}

type ProgrammablePublish struct {
	Modules      [][]byte
	Dependencies []*sui.ObjectId
}

func (p *ProgrammablePublish) String() string {
	deps := make([]string, len(p.Dependencies))
	for i, dep := range p.Dependencies {
		deps[i] = dep.String()
	}
	return fmt.Sprintf("Publish(%d modules, deps: [%s])", len(p.Modules), strings.Join(deps, ", "))
}

type ProgrammableMakeMoveVec struct {
	Type    *sui.TypeTag `bcs:"optional"`
	Objects []Argument
}

func (p *ProgrammableMakeMoveVec) String() string {
	objects := make([]string, len(p.Objects))
	for i, obj := range p.Objects {
		objects[i] = obj.String()
	}
	typeStr := "None"
	if p.Type != nil {
		typeStr = (*p.Type).String()
	}
	return fmt.Sprintf("MakeMoveVec<%s>([%s])", typeStr, strings.Join(objects, ", "))
}

type ProgrammableUpgrade struct {
	Modules      [][]byte
	Dependencies []*sui.ObjectId
	PackageId    *sui.PackageId
	Ticket       Argument
}

func (p *ProgrammableUpgrade) String() string {
	deps := make([]string, len(p.Dependencies))
	for i, dep := range p.Dependencies {
		deps[i] = dep.String()
	}
	return fmt.Sprintf("Upgrade(package: %s, %d modules, deps: [%s], ticket: %s)",
		p.PackageId.String(), len(p.Modules), strings.Join(deps, ", "), p.Ticket.String())
}
