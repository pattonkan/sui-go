package pkg

import (
	"context"
	"fmt"

	"github.com/fardream/go-bcs/bcs"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

func SwapSui(
	suiClient *sui.ImplSuiAPI,
	swapper *sui_signer.Signer,
	swapPackageID *sui_types.PackageID,
	testcoinID *sui_types.ObjectID,
	poolObjectID *sui_types.ObjectID,
	suiCoins []*models.Coin,
) {
	poolGetObjectRes, err := suiClient.GetObject(context.Background(), &sui.GetObjectRequest{
		ObjectID: poolObjectID,
		Options: &models.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// swap sui to testcoin
	ptb := sui_types.NewProgrammableTransactionBuilder()

	arg0 := ptb.MustObj(sui_types.ObjectArg{SharedObject: &sui_types.SharedObjectArg{
		Id:                   poolObjectID,
		InitialSharedVersion: poolGetObjectRes.Data.Ref().Version,
		Mutable:              true,
	}})
	arg1 := ptb.MustObj(sui_types.ObjectArg{ImmOrOwnedObject: suiCoins[0].Ref()})

	retCoinArg := ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:  swapPackageID,
			Module:   "swap",
			Function: "swap_sui",
			TypeArguments: []sui_types.TypeTag{{Struct: &sui_types.StructTag{
				Address: testcoinID,
				Module:  "testcoin",
				Name:    "TESTCOIN",
			}}},
			Arguments: []sui_types.Argument{arg0, arg1},
		}},
	)
	ptb.Command(sui_types.Command{
		TransferObjects: &sui_types.ProgrammableTransferObjects{
			Objects: []sui_types.Argument{retCoinArg},
			Address: ptb.MustPure(swapper.Address),
		},
	})
	pt := ptb.Finish()
	txData := sui_types.NewProgrammable(
		swapper.Address,
		pt,
		[]*sui_types.ObjectRef{suiCoins[1].Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		panic(err)
	}

	resp, err := suiClient.SignAndExecuteTransaction(
		context.Background(),
		swapper,
		txBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowObjectChanges: true,
			ShowEffects:       true,
		},
	)
	if err != nil || !resp.Effects.Data.IsSuccess() {
		panic(err)
	}

	for _, change := range resp.ObjectChanges {
		if change.Data.Created != nil {
			fmt.Println("change.Data.Created.ObjectID: ", change.Data.Created.ObjectID)
			fmt.Println("change.Data.Created.ObjectType: ", change.Data.Created.ObjectType)
			fmt.Println("change.Data.Created.Owner.AddressOwner: ", change.Data.Created.Owner.AddressOwner)
		}
		if change.Data.Mutated != nil {
			fmt.Println("change.Data.Mutated.ObjectID: ", change.Data.Mutated.ObjectID)
			fmt.Println("change.Data.Mutated.ObjectType: ", change.Data.Mutated.ObjectType)
			fmt.Println("change.Data.Mutated.Owner.AddressOwner: ", change.Data.Mutated.Owner.AddressOwner)
		}
	}
}
