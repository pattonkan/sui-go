package pkg

import (
	"context"

	"github.com/fardream/go-bcs/bcs"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/suiptb"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suisigner"
)

func CreatePool(
	suiClient *suiclient.ClientImpl,
	signer *suisigner.Signer,
	swapPackageId *sui.PackageId,
	testcoinId *sui.ObjectId,
	testCoin *suiclient.Coin,
	suiCoins []*suiclient.Coin,
) *sui.ObjectId {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	arg0 := ptb.MustObj(suiptb.ObjectArg{ImmOrOwnedObject: testCoin.Ref()})
	arg1 := ptb.MustObj(suiptb.ObjectArg{ImmOrOwnedObject: suiCoins[0].Ref()})
	arg2 := ptb.MustPure(uint64(3))

	lspArg := ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:  swapPackageId,
			Module:   "swap",
			Function: "create_pool",
			TypeArguments: []sui.TypeTag{{Struct: &sui.StructTag{
				Address: testcoinId,
				Module:  "testcoin",
				Name:    "TESTCOIN",
			}}},
			Arguments: []suiptb.Argument{arg0, arg1, arg2},
		}},
	)
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{lspArg},
			Address: ptb.MustPure(signer.Address),
		},
	})
	pt := ptb.Finish()
	txData := suiptb.NewTransactionData(
		signer.Address,
		pt,
		[]*sui.ObjectRef{suiCoins[1].Ref()},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		panic(err)
	}

	txnResponse, err := suiClient.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		panic(err)
	}
	for _, change := range txnResponse.ObjectChanges {
		if change.Data.Created != nil {
			resource, err := sui.NewResourceType(change.Data.Created.ObjectType)
			if err != nil {
				panic(err)
			}
			if resource.Contains(nil, "swap", "Pool") {
				return &change.Data.Created.ObjectId
			}
		}
	}

	return nil
}
