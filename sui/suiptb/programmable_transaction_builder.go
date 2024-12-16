package suiptb

import (
	"errors"
	"fmt"

	"github.com/fardream/go-bcs/bcs"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/utils/indexmap"
)

// ProgrammableTransactionBuilder calls 'Finish()' to be built
// into 'ProgrammableTransaction' for following encoding in BCS format
type ProgrammableTransactionBuilder struct {
	Inputs   *indexmap.IndexMap[BuilderArg, CallArg] //maybe has hash clash
	Commands []Command
}

// ProgrammableTransaction is the packed immediate transaction type which will be encoded
// into BCS format by 'NewTransactionData()' or 'NewTransactionDataAllowSponsor()'
type ProgrammableTransaction struct {
	Inputs   []CallArg
	Commands []Command
}

type BuilderArg struct {
	Object              *sui.ObjectId
	Pure                *[]uint8
	ForcedNonUniquePure *uint
}

func NewTransactionDataTransactionBuilder() *ProgrammableTransactionBuilder {
	return &ProgrammableTransactionBuilder{
		Inputs: indexmap.NewIndexMap[BuilderArg, CallArg](),
	}
}

func (p *ProgrammableTransactionBuilder) Finish() ProgrammableTransaction {
	var inputs []CallArg
	p.Inputs.ForEach(func(k BuilderArg, v CallArg) {
		inputs = append(inputs, v)
	})
	return ProgrammableTransaction{
		Inputs:   inputs,
		Commands: p.Commands,
	}
}

func (p *ProgrammableTransactionBuilder) Pure(value any) (Argument, error) {
	pureData, err := bcs.Marshal(value)
	if err != nil {
		return Argument{}, err
	}
	return p.pureBytes(pureData, false), nil
}

func (p *ProgrammableTransactionBuilder) MustPure(value any) Argument {
	pureData, err := bcs.Marshal(value)
	if err != nil {
		panic(err)
	}
	return p.pureBytes(pureData, false)
}

// refer crates/sui-types/src/programmable_transaction_builder.rs
func (p *ProgrammableTransactionBuilder) Obj(objArg ObjectArg) (Argument, error) {
	id := objArg.id()
	var oj ObjectArg
	if oldValue, ok := p.Inputs.Get(BuilderArg{Object: id}); ok {
		var oldObjArg ObjectArg
		switch {
		case oldValue.Pure != nil:
			return Argument{}, errors.New("invariant violation! object has Pure argument")
		case oldValue.Object != nil:
			oldObjArg = *oldValue.Object
		}

		switch {
		case oldObjArg.SharedObject.InitialSharedVersion == objArg.SharedObject.InitialSharedVersion:
			if oldObjArg.id() != objArg.id() {
				return Argument{}, errors.New("invariant violation! object has id does not match call arg")
			}
			oj = ObjectArg{
				SharedObject: &SharedObjectArg{
					Id:                   id,
					InitialSharedVersion: objArg.SharedObject.InitialSharedVersion,
					Mutable:              oldObjArg.SharedObject.Mutable || objArg.SharedObject.Mutable,
				},
			}
		default:
			if oldObjArg != objArg {
				return Argument{}, fmt.Errorf(
					"mismatched Object argument kind for object %s. "+
						"%v is not compatible with %v", id.String(), oldValue, objArg,
				)
			}
			oj = objArg
		}
	} else {
		oj = objArg
	}
	i := uint16(p.Inputs.InsertFull(
		BuilderArg{Object: id},
		CallArg{Object: &oj},
	))
	return Argument{Input: &i}, nil
}

func (p *ProgrammableTransactionBuilder) MustObj(objArg ObjectArg) Argument {
	arg, err := p.Obj(objArg)
	if err != nil {
		panic(err)
	}
	return arg
}

func (p *ProgrammableTransactionBuilder) ForceSeparatePure(value any) (Argument, error) {
	pureData, err := bcs.Marshal(value)
	if err != nil {
		return Argument{}, err
	}
	return p.pureBytes(pureData, true), nil
}

func (p *ProgrammableTransactionBuilder) MustForceSeparatePure(value any) Argument {
	arg, err := p.ForceSeparatePure(value)
	if err != nil {
		panic(err)
	}
	return arg
}

func (p *ProgrammableTransactionBuilder) pureBytes(bytes []byte, forceSeparate bool) Argument {
	var arg BuilderArg
	if forceSeparate {
		length := uint(p.Inputs.Len())
		arg = BuilderArg{
			ForcedNonUniquePure: &length,
		}
	} else {
		arg = BuilderArg{
			Pure: &bytes,
		}
	}
	i := uint16(p.Inputs.InsertFull(
		arg,
		CallArg{Pure: &bytes},
	))
	return Argument{
		Input: &i,
	}
}

// developers should only use `Pure()`, `MustPure()` and `Obj()` to create PTB Arguments
// `Input()` is a function for internal usage
// TODO add explanation for `Input()`
func (p *ProgrammableTransactionBuilder) Input(callArg CallArg) (Argument, error) {
	switch {
	case callArg.Pure != nil:
		return p.pureBytes(*callArg.Pure, false), nil
	case callArg.Object != nil:
		return p.Obj(*callArg.Object)
	default:
		return Argument{}, errors.New("this callArg is nil")
	}
}

// Add command to `ProgrammableTransactionBuilder.Commands`, and return the result in `Argument` type
func (p *ProgrammableTransactionBuilder) Command(command Command) Argument {
	p.Commands = append(p.Commands, command)
	i := uint16(len(p.Commands)) - 1
	return Argument{
		Result: &i,
	}
}

//// ProgrammableTransactionBuilder fast API calls ////

func (p *ProgrammableTransactionBuilder) MakeObjVec(objs []ObjectArg) (Argument, error) {
	var objArgs []Argument
	for _, v := range objs {
		objArg, err := p.Obj(v)
		if err != nil {
			return Argument{}, err
		}
		objArgs = append(objArgs, objArg)
	}
	arg := p.Command(Command{
		MakeMoveVec: &ProgrammableMakeMoveVec{Type: nil, Objects: objArgs},
	})
	return arg, nil
}

// construct `move_call` with argument `CallArg` type
func (p *ProgrammableTransactionBuilder) MoveCall(
	packageId *sui.PackageId,
	module sui.Identifier,
	function sui.Identifier,
	typeArguments []sui.TypeTag,
	callArgs []CallArg,
) error {
	var arguments []Argument
	for _, v := range callArgs {
		argument, err := p.Input(v)
		if err != nil {
			return err
		}
		arguments = append(arguments, argument)
	}
	p.Command(Command{
		MoveCall: &ProgrammableMoveCall{
			Package:       packageId,
			Module:        module,
			Function:      function,
			TypeArguments: typeArguments,
			Arguments:     arguments,
		}},
	)
	return nil
}

// construct `move_call` with argument `Argument` type, and return `Argument`
func (p *ProgrammableTransactionBuilder) ProgrammableMoveCall(
	packageId *sui.PackageId,
	module sui.Identifier,
	function sui.Identifier,
	typeArguments []sui.TypeTag,
	arguments []Argument,
) Argument {
	return p.Command(Command{
		MoveCall: &ProgrammableMoveCall{
			Package:       packageId,
			Module:        module,
			Function:      function,
			TypeArguments: typeArguments,
			Arguments:     arguments,
		}},
	)
}

func (p *ProgrammableTransactionBuilder) PublishUpgradeable(
	modules [][]byte,
	dependencies []*sui.ObjectId,
) Argument {
	return p.Command(Command{
		Publish: &ProgrammablePublish{
			Modules:      modules,
			Dependencies: dependencies,
		}},
	)
}

func (p *ProgrammableTransactionBuilder) PublishImmutable(
	modules [][]byte,
	dependencies []*sui.ObjectId,
) Argument {
	return p.Command(Command{
		MoveCall: &ProgrammableMoveCall{
			Package:       sui.SuiPackageIdSuiFramework,
			Module:        sui.SuiSystemModuleName,
			Function:      "make_immutable",
			TypeArguments: nil,
			Arguments:     []Argument{p.PublishUpgradeable(modules, dependencies)},
		}},
	)
}

func (p *ProgrammableTransactionBuilder) Upgrade(
	currentPackageObjectId *sui.ObjectId,
	upgradeTicket Argument,
	transitiveDeps []*sui.ObjectId,
	modules [][]byte,
) Argument {
	return p.Command(Command{
		Upgrade: &ProgrammableUpgrade{
			Modules:      modules,
			Dependencies: transitiveDeps,
			PackageId:    currentPackageObjectId,
			Ticket:       upgradeTicket,
		}},
	)
}

func (p *ProgrammableTransactionBuilder) TransferArg(recipient *sui.Address, arg Argument) {
	p.TransferArgs(recipient, []Argument{arg})
}

func (p *ProgrammableTransactionBuilder) TransferArgs(recipient *sui.Address, args []Argument) {
	p.Command(Command{
		TransferObjects: &ProgrammableTransferObjects{
			Objects: args,
			Address: p.MustPure(recipient),
		}},
	)
}

func (p *ProgrammableTransactionBuilder) TransferObject(recipient *sui.Address, objectRef *sui.ObjectRef) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return fmt.Errorf("can't add recipient as arg: %w", err)
	}
	objArg, err := p.Obj(ObjectArg{ImmOrOwnedObject: objectRef})
	if err != nil {
		return err
	}
	p.Command(Command{
		TransferObjects: &ProgrammableTransferObjects{
			Objects: []Argument{objArg},
			Address: recArg,
		}},
	)
	return nil
}

func (p *ProgrammableTransactionBuilder) TransferSui(recipient *sui.Address, amount *uint64) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return fmt.Errorf("can't add recipient as arg: %w", err)
	}
	var coinArg Argument
	if amount != nil {
		amtArg := p.MustPure(amount)
		coinArg = p.Command(Command{
			SplitCoins: &ProgrammableSplitCoins{
				Coin:    Argument{GasCoin: &sui.EmptyEnum{}},
				Amounts: []Argument{amtArg},
			}},
		)
	} else {
		coinArg = Argument{GasCoin: &sui.EmptyEnum{}}
	}
	p.Command(Command{
		TransferObjects: &ProgrammableTransferObjects{
			Objects: []Argument{coinArg},
			Address: recArg,
		}},
	)
	return nil
}

// the gas coin is consumed as the coin to be paid
func (p *ProgrammableTransactionBuilder) PayAllSui(recipient *sui.Address) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return fmt.Errorf("can't add recipient as arg: %w", err)
	}
	p.Command(Command{
		TransferObjects: &ProgrammableTransferObjects{
			Objects: []Argument{{GasCoin: &sui.EmptyEnum{}}},
			Address: recArg,
		}},
	)
	return nil
}

// the gas coin is consumed as the coin to be paid
func (p *ProgrammableTransactionBuilder) PaySui(recipients []*sui.Address, amounts []uint64) error {
	return p.payImpl(recipients, amounts, Argument{GasCoin: &sui.EmptyEnum{}})
}

// merge all given coins into the 1st coin, and pay it
// with the corresponding amounts to the corresponding recipients
func (p *ProgrammableTransactionBuilder) Pay(
	coins []*sui.ObjectRef,
	recipients []*sui.Address,
	amounts []uint64,
) error {
	if len(coins) == 0 {
		return errors.New("coins vector is empty")
	}
	coinArg, err := p.Obj(ObjectArg{ImmOrOwnedObject: coins[0]})
	if err != nil {
		return err
	}
	coins = coins[1:]

	var mergeArgs []Argument
	for _, v := range coins {
		mergeCoin, err := p.Obj(ObjectArg{ImmOrOwnedObject: v})
		if err != nil {
			return err
		}
		mergeArgs = append(mergeArgs, mergeCoin)
	}
	if len(mergeArgs) != 0 {
		p.Command(
			Command{
				MergeCoins: &ProgrammableMergeCoins{
					Destination: coinArg,
					Sources:     mergeArgs,
				},
			},
		)
	}
	return p.payImpl(recipients, amounts, coinArg)
}

// And the commands to pay a coin object to multiple recipients
// golang implementation of pay_impl() in `sui/crates/sui-types/src/programmable_transaction_builder.rs`
func (p *ProgrammableTransactionBuilder) payImpl(
	recipients []*sui.Address,
	amounts []uint64,
	coin Argument, // the coin to be consumed
) error {
	if len(recipients) != len(amounts) {
		return fmt.Errorf(
			"recipients and amounts mismatch. Got %d recipients but %d amounts",
			len(recipients),
			len(amounts),
		)
	}
	if len(amounts) == 0 {
		return nil
	}

	var amtArgs []Argument
	// map[<recipients accounts>]<index in input amounts array>. The `[]int` array is `split_secondaries` in rust-sdk
	var recipientMap = make(map[*sui.Address][]int)
	// this allows us to traverse the `recipientMap` with order (like indexmap)
	var recipientMapKeyIndex []*sui.Address

	for i := 0; i < len(amounts); i++ {
		amt, err := p.Pure(amounts[i])
		if err != nil {
			return err
		}
		recipientMap[recipients[i]] = append(recipientMap[recipients[i]], i)
		if len(recipientMap[recipients[i]]) == 1 {
			recipientMapKeyIndex = append(recipientMapKeyIndex, recipients[i])
		}
		amtArgs = append(amtArgs, amt)
	}
	splitCoinResult := p.Command(
		Command{
			SplitCoins: &ProgrammableSplitCoins{
				Coin:    coin,
				Amounts: amtArgs,
			},
		},
	)
	if splitCoinResult.Result == nil {
		return errors.New("self.command should always give a Argument::Result")
	}
	for _, v := range recipientMapKeyIndex {
		recArg, err := p.Pure(v)
		if err != nil {
			return err
		}
		var coins []Argument
		for _, j := range recipientMap[v] {
			// the portions of the coins that slipt from the given coin, which are going to pay for recipients
			coinTransfer := Argument{
				NestedResult: &NestedResult{Cmd: *splitCoinResult.Result, Result: uint16(j)},
			}
			coins = append(coins, coinTransfer)
		}
		p.Command(
			Command{
				TransferObjects: &ProgrammableTransferObjects{
					Objects: coins,
					Address: recArg,
				},
			},
		)
	}
	return nil
}
