package main

import (
	"context"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"
)

func main() {
	client, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)

	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: signer.Address,
		Limit: 4,
	})
	if err != nil {
		panic(err)
	}
	coins := suiclient.Coins(coinPages.Data)

	packageId := BuildDeployContract(client, signer, utils.GetGitRoot()+"/examples/recast/recast")
	fmt.Println("packageId: ", packageId.String())

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	argContainer := ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       packageId,
			Module:        "recast",
			Function:      "create_container",
			TypeArguments: []sui.TypeTag{},
			Arguments:     []suiptb.Argument{},
		}},
	)

	argU64 := ptb.MustPure(uint64(42949672970))
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       packageId,
			Module:        "recast",
			Function:      "try_recast",
			TypeArguments: []sui.TypeTag{},
			Arguments:     []suiptb.Argument{argU64, argContainer},
		}},
	)

	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		signer.Address,
		pt,
		[]*sui.ObjectRef{coins[2].Ref()},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx.V1.Kind)
	if err != nil {
		panic(err)
	}

	resp, err := client.DevInspectTransactionBlock(
		context.Background(),
		&suiclient.DevInspectTransactionBlockRequest{
			SenderAddress: signer.Address,
			TxKindBytes:   txBytes,
		},
	)
	// resp, err := client.DryRunTransaction(context.Background(), txBytes)
	// resp, err := client.SignAndExecuteTransaction(
	// 	context.Background(),
	// 	signer,
	// 	txBytes,
	// 	&SuiTransactionBlockResponseOptions{
	// 		ShowEffects: true,
	// 	},
	// )
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", resp)
}

func BuildDeployContract(
	client *suiclient.ClientImpl,
	signer *suisigner.Signer,
	contractPath string,
) *sui.PackageId {
	modules, err := utils.MoveBuild(contractPath)
	if err != nil {
		panic(err)
	}

	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	if err != nil {
		panic(err)
	}

	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&suiclient.SuiTransactionBlockResponseOptions{
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

	packageId, err := txnResponse.GetPublishedPackageId()
	if err != nil {
		panic(err)
	}

	return packageId
}
