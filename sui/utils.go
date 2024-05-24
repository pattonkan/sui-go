package sui

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"
	"github.com/stretchr/testify/require"
)

// test only requires `ShowObjectChanges: true`
func GetCreatedObjectIdAndType(
	txRes *models.SuiTransactionBlockResponse,
	moduleName string,
	objectName string,
) (*sui_types.ObjectID, string, error) {
	if txRes.ObjectChanges == nil {
		return nil, "", fmt.Errorf("no ObjectChanges")
	}
	for _, change := range txRes.ObjectChanges {
		if change.Data.Created != nil {
			// FIXME error-prone, we need to parse the object type to check
			// some possible examples
			// * 0x2::coin::TreasuryCap<0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::testcoin::TESTCOIN>
			// * 0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::anchor::Anchor
			if strings.Contains(change.Data.Created.ObjectType, fmt.Sprintf("%s::%s", moduleName, objectName)) {
				return &change.Data.Created.ObjectID, change.Data.Created.ObjectType, nil
			}
		}
	}
	return nil, "", fmt.Errorf("not found")
}

// test only
func BuildDeployContract(
	t *testing.T,
	client *ImplSuiAPI,
	signer *sui_signer.Signer,
	contractPath string,
) *sui_types.PackageID {
	modules, err := utils.MoveBuild(contractPath)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		signer.Address,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(DefaultGasBudget),
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		signer,
		txnBytes.TxBytes,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID, err := txnResponse.GetPublishedPackageID()
	require.NoError(t, err)

	return packageID
}

// test only
func BuildDeployMintCoin(
	t *testing.T,
	client *ImplSuiAPI,
	signer *sui_signer.Signer,
	contractPath string,
	mintAmount uint64,
	tokenName string,
) (*sui_types.PackageID, *sui_types.ObjectID) {
	modules, err := utils.MoveBuild(contractPath)
	require.NoError(t, err)

	txnBytes, err := client.Publish(context.Background(), signer.Address, modules.Modules, modules.Dependencies, nil, models.NewSafeSuiBigInt(uint64(100000000)))
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, &models.SuiTransactionBlockResponseOptions{
		ShowEffects:       true,
		ShowObjectChanges: true,
	})
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID, err := txnResponse.GetPublishedPackageID()
	require.NoError(t, err)

	treasuryCap, _, err := GetCreatedObjectIdAndType(txnResponse, "coin", "TreasuryCap")
	require.NoError(t, err)

	txnRes, err := client.MintToken(
		context.Background(),
		signer,
		packageID,
		tokenName,
		treasuryCap,
		mintAmount,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnRes.Effects.Data.IsSuccess())

	return packageID, treasuryCap
}
