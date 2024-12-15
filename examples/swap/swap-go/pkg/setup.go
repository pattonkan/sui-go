package pkg

import (
	"context"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suisigner"
	"github.com/howjmay/sui-go/utils"
)

func BuildAndPublish(client *suiclient.ClientImpl, signer *suisigner.Signer, path string) *sui.PackageId {
	modules, err := utils.MoveBuild(path)
	if err != nil {
		panic(err)
	}
	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(10 * suiclient.DefaultGasBudget),
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
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		panic(err)
	}
	packageId, err := txnResponse.GetPublishedPackageId()
	if err != nil {
		panic(err)
	}
	return packageId
}

func BuildDeployMintTestcoin(client *suiclient.ClientImpl, signer *suisigner.Signer) (
	*sui.PackageId,
	*sui.ObjectId,
) {
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin/")
	if err != nil {
		panic(err)
	}

	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(10 * suiclient.DefaultGasBudget),
		},
	)
	if err != nil {
		panic(err)
	}
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), signer, txnBytes.TxBytes, &suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		panic(err)
	}

	packageId, err := txnResponse.GetPublishedPackageId()
	if err != nil {
		panic(err)
	}

	treasuryCap, _, err := txnResponse.GetCreatedObjectInfo("coin", "TreasuryCap")
	if err != nil {
		panic(err)
	}

	mintAmount := uint64(1000000)
	txnResponse, err = client.MintToken(
		context.Background(),
		signer,
		packageId,
		"testcoin",
		treasuryCap,
		mintAmount,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		panic(err)
	}

	return packageId, treasuryCap
}
