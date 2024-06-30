package pkg

import (
	"context"

	"github.com/fardream/go-bcs/bcs"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

func CreatePool(
	suiClient *sui.ImplSuiAPI,
	signer *sui_signer.Signer,
	swapPackageID *sui_types.PackageID,
	testcoinID *sui_types.ObjectID,
	testCoin *models.Coin,
	suiCoins []*models.Coin,
) *sui_types.ObjectID {
	ptb := sui_types.NewProgrammableTransactionBuilder()

	arg0 := ptb.MustObj(sui_types.ObjectArg{ImmOrOwnedObject: testCoin.Ref()})
	arg1 := ptb.MustObj(sui_types.ObjectArg{ImmOrOwnedObject: suiCoins[0].Ref()})
	arg2 := ptb.MustPure(uint64(3))

	lspArg := ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:  swapPackageID,
			Module:   "swap",
			Function: "create_pool",
			TypeArguments: []sui_types.TypeTag{{Struct: &sui_types.StructTag{
				Address: testcoinID,
				Module:  "testcoin",
				Name:    "TESTCOIN",
			}}},
			Arguments: []sui_types.Argument{arg0, arg1, arg2},
		}},
	)
	ptb.Command(sui_types.Command{
		TransferObjects: &sui_types.ProgrammableTransferObjects{
			Objects: []sui_types.Argument{lspArg},
			Address: ptb.MustPure(signer.Address),
		},
	})
	pt := ptb.Finish()
	txData := sui_types.NewProgrammable(
		signer.Address,
		pt,
		[]*sui_types.ObjectRef{suiCoins[1].Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		panic(err)
	}

	txnResponse, err := suiClient.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		panic(err)
	}
	for _, change := range txnResponse.ObjectChanges {
		if change.Data.Created != nil {
			resource, err := models.NewResourceType(change.Data.Created.ObjectType)
			if err != nil {
				panic(err)
			}
			if resource.Contains(nil, "swap", "Pool") {
				return &change.Data.Created.ObjectID
			}
		}
	}

	return nil
}
