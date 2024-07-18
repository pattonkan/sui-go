package main

import (
	"context"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"
)

func main() {
	client, signer := sui.NewSuiClient(conn.LocalnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)

	coinPages, err := client.GetCoins(context.Background(), &sui.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 4,
	})
	if err != nil {
		panic(err)
	}
	coins := models.Coins(coinPages.Data)

	packageID := BuildDeployContract(client, signer, utils.GetGitRoot()+"/examples/recast/recast")
	fmt.Println("packageID: ", packageID.String())

	ptb := sui_types.NewProgrammableTransactionBuilder()
	argContainer := ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:       packageID,
			Module:        "recast",
			Function:      "create_container",
			TypeArguments: []sui_types.TypeTag{},
			Arguments:     []sui_types.Argument{},
		}},
	)

	argU64 := ptb.MustPure(uint64(42949672970))
	ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:       packageID,
			Module:        "recast",
			Function:      "try_recast",
			TypeArguments: []sui_types.TypeTag{},
			Arguments:     []sui_types.Argument{argU64, argContainer},
		}},
	)

	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		signer.Address,
		pt,
		[]*sui_types.ObjectRef{coins[2].Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx.V1.Kind)
	if err != nil {
		panic(err)
	}

	resp, err := client.DevInspectTransactionBlock(
		context.Background(),
		&sui.DevInspectTransactionBlockRequest{
			SenderAddress: signer.Address,
			TxKindBytes:   txBytes,
		},
	)
	// resp, err := client.DryRunTransaction(context.Background(), txBytes)
	// resp, err := client.SignAndExecuteTransaction(
	// 	context.Background(),
	// 	signer,
	// 	txBytes,
	// 	&models.SuiTransactionBlockResponseOptions{
	// 		ShowEffects: true,
	// 	},
	// )
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", resp)
}

func BuildDeployContract(
	client *sui.ImplSuiAPI,
	signer *sui_signer.Signer,
	contractPath string,
) *sui_types.PackageID {
	modules, err := utils.MoveBuild(contractPath)
	if err != nil {
		panic(err)
	}

	txnBytes, err := client.Publish(
		context.Background(),
		&sui.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	if err != nil {
		panic(err)
	}

	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if !txnResponse.Effects.Data.IsSuccess() {
		panic(txnResponse.Effects.Data.V1.Status.Error)
	}

	packageID, err := txnResponse.GetPublishedPackageID()
	if err != nil {
		panic(err)
	}

	return packageID
}
