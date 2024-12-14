package pkg

import (
	"context"
	"fmt"

	"github.com/fardream/go-bcs/bcs"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/suiptb"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suisigner"
)

func SwapSui(
	suiClient *suiclient.ClientImpl,
	swapper *suisigner.Signer,
	swapPackageId *sui.PackageId,
	testcoinId *sui.ObjectId,
	poolObjectId *sui.ObjectId,
	suiCoins []*suiclient.Coin,
) {
	poolGetObjectRes, err := suiClient.GetObject(context.Background(), &suiclient.GetObjectRequest{
		ObjectId: poolObjectId,
		Options: &suiclient.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// swap sui to testcoin
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	arg0 := ptb.MustObj(suiptb.ObjectArg{SharedObject: &suiptb.SharedObjectArg{
		Id:                   poolObjectId,
		InitialSharedVersion: poolGetObjectRes.Data.Ref().Version,
		Mutable:              true,
	}})
	arg1 := ptb.MustObj(suiptb.ObjectArg{ImmOrOwnedObject: suiCoins[0].Ref()})

	retCoinArg := ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:  swapPackageId,
			Module:   "swap",
			Function: "swap_sui",
			TypeArguments: []sui.TypeTag{{Struct: &sui.StructTag{
				Address: testcoinId,
				Module:  "testcoin",
				Name:    "TESTCOIN",
			}}},
			Arguments: []suiptb.Argument{arg0, arg1},
		}},
	)
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{retCoinArg},
			Address: ptb.MustPure(swapper.Address),
		},
	})
	pt := ptb.Finish()
	txData := suiptb.NewTransactionData(
		swapper.Address,
		pt,
		[]*sui.ObjectRef{suiCoins[1].Ref()},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		panic(err)
	}

	resp, err := suiClient.SignAndExecuteTransaction(
		context.Background(),
		swapper,
		txBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowObjectChanges: true,
			ShowEffects:       true,
		},
	)
	if err != nil || !resp.Effects.Data.IsSuccess() {
		panic(err)
	}

	for _, change := range resp.ObjectChanges {
		if change.Data.Created != nil {
			fmt.Println("change.Data.Created.ObjectId: ", change.Data.Created.ObjectId)
			fmt.Println("change.Data.Created.ObjectType: ", change.Data.Created.ObjectType)
			fmt.Println("change.Data.Created.Owner.AddressOwner: ", change.Data.Created.Owner.AddressOwner)
		}
		if change.Data.Mutated != nil {
			fmt.Println("change.Data.Mutated.ObjectId: ", change.Data.Mutated.ObjectId)
			fmt.Println("change.Data.Mutated.ObjectType: ", change.Data.Mutated.ObjectType)
			fmt.Println("change.Data.Mutated.Owner.AddressOwner: ", change.Data.Mutated.Owner.AddressOwner)
		}
	}
}
