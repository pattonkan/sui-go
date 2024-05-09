package isc_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/isc"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
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

	_, err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/isc/contracts/isc/")
	require.NoError(t, err)

	txnBytes, err := client.Publish(context.Background(), sui_signer.TEST_ADDRESS, modules.Modules, modules.Dependencies, nil, models.NewSafeSuiBigInt(uint64(100000000)))
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, &models.SuiTransactionBlockResponseOptions{
		ShowEffects:       true,
		ShowObjectChanges: true,
	})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)

	packageID := txnResponse.GetPublishedPackageID()
	t.Log("packageID: ", packageID)

	startNewChainRes, err := client.StartNewChain(context.Background(), signer, packageID, 10000000, &models.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)
	t.Logf("StartNewChain response: %#v\n", startNewChainRes)
}
