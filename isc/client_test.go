package isc_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/howjmay/sui-go/isc"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"

	"github.com/stretchr/testify/require"
)

type Client struct {
	API sui.ImplSuiAPI
}

func TestStartNewChain(t *testing.T) {
	t.Skip("only for localnet")
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	digest, err := sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)
	digest, err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)

	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/isc/contracts/isc/")
	require.NoError(t, err)

	coins, err := client.API.GetCoins(context.Background(), sui_signer.TEST_ADDRESS, nil, nil, 10)
	require.NoError(t, err)
	gasBudget := uint64(100000000)
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(100000), gasBudget, 10, 10)
	require.NoError(t, err)

	txnBytes, err := client.API.Publish(context.Background(), sui_signer.TEST_ADDRESS, modules.Modules, modules.Dependencies, pickedCoins.CoinIds()[0], models.NewSafeSuiBigInt(gasBudget))
	require.NoError(t, err)
	signature, err := signer.SignTransactionBlock(txnBytes.TxBytes, sui_signer.DefaultIntent())
	require.NoError(t, err)
	txnResponse, err := client.API.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes, []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	require.NoError(t, err)

	var packageID sui_types.PackageID
	fmt.Println("txnResponse.Effects.Data.V1.Status: ", txnResponse.Effects.Data.V1.Status)
	for _, change := range txnResponse.ObjectChanges {
		fmt.Println("change: ", change.Data)
		if change.Data.Published != nil {
			tmp := change.Data.Published.PackageId
			packageID = tmp
		}
		if change.Data.Mutated != nil {
			fmt.Println("change.Data.Mutated: ", change.Data.Mutated)
		}
	}
	fmt.Println("packageID: ", packageID.String())

	// packageID, anchorCap := isc.GetIscPackageIDAndAnchor(isc.GetGitRoot() + "/isc/contracts/isc/publish_receipt.json")

	res, err := client.StartNewChain(context.Background(), signer, &packageID, nil)
	require.NoError(t, err)
	require.Equal(t, res.Effects.Data.V1.Status.Status, "success")
	t.Logf("res.Effects.Data.V1.Status: %#v\n", res.Effects.Data.V1.Status)
	t.Logf("StartNewChain response: %#v\n", res)
	// for _, change := range res.ObjectChanges {
	// 	fmt.Println("change.Data.Created: ", change.Data.Created)
	// }
}
