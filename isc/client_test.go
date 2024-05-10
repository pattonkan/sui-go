package isc_test

import (
	"context"
	"fmt"
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

	startNewChainRes, err := client.StartNewChain(
		context.Background(),
		signer,
		packageID,
		sui.DefaultGasBudget,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, startNewChainRes.Effects.Data.V1.Status.Status)
	t.Logf("StartNewChain response: %#v\n", startNewChainRes)
}

func TestSendCoin(t *testing.T) {
	t.Skip("only for localnet")
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	_, err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := buildAndDeployIscContracts(t, client, signer)
	tokenPackageID, _ := buildDeployMintTestcoin(t, client, signer)

	// start a new chain
	startNewChainRes, err := client.StartNewChain(
		context.Background(),
		signer,
		iscPackageID,
		sui.DefaultGasBudget,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, startNewChainRes.Effects.Data.V1.Status.Status)

	anchorObjID, _, err := sui.GetCreatedObjectIdAndType(startNewChainRes, "anchor", "Anchor")
	coinType := fmt.Sprintf("%s::testcoin::TESTCOIN", tokenPackageID.String())
	require.NoError(t, err)
	// the signer should have only one coin object which belongs to testcoin type
	coins, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, 10)
	require.NoError(t, err)

	sendCoinRes, err := client.SendCoin(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		coinType,
		coins.Data[0].CoinObjectID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, sendCoinRes.Effects.Data.V1.Status.Status)

	getObjectRes, err := client.GetObject(context.Background(), coins.Data[0].CoinObjectID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, anchorObjID.String(), getObjectRes.Data.Owner.ObjectOwnerInternal.AddressOwner.String())
}

func TestReceiveCoin(t *testing.T) {
	t.Skip("only for localnet")
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	_, err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := buildAndDeployIscContracts(t, client, signer)
	tokenPackageID, _ := buildDeployMintTestcoin(t, client, signer)

	// start a new chain
	startNewChainRes, err := client.StartNewChain(
		context.Background(),
		signer,
		iscPackageID,
		sui.DefaultGasBudget,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, startNewChainRes.Effects.Data.V1.Status.Status)

	anchorObjID, _, err := sui.GetCreatedObjectIdAndType(startNewChainRes, "anchor", "Anchor")
	coinType := fmt.Sprintf("%s::testcoin::TESTCOIN", tokenPackageID.String())
	require.NoError(t, err)
	// the signer should have only one coin object which belongs to testcoin type
	coins, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, 10)
	require.NoError(t, err)

	sendCoinRes, err := client.SendCoin(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		coinType,
		coins.Data[0].CoinObjectID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, sendCoinRes.Effects.Data.V1.Status.Status)

	getObjectRes, err := client.GetObject(context.Background(), coins.Data[0].CoinObjectID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, anchorObjID.String(), getObjectRes.Data.Owner.ObjectOwnerInternal.AddressOwner.String())

	receiveCoinRes, err := client.ReceiveCoin(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		coinType,
		coins.Data[0].CoinObjectID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, receiveCoinRes.Effects.Data.V1.Status.Status)

	// TODO we should check the isc on-chain balance
}

func buildAndDeployIscContracts(t *testing.T, client *isc.Client, signer *sui_signer.Signer) *sui_types.PackageID {
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

	return packageID
}

func buildDeployMintTestcoin(t *testing.T, client *isc.Client, signer *sui_signer.Signer) (*sui_types.PackageID, *sui_types.ObjectID) {
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin/")
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

	treasuryCap, _, err := sui.GetCreatedObjectIdAndType(txnResponse, "coin", "TreasuryCap")
	require.NoError(t, err)

	mintAmount := uint64(1000000)
	txnRes, err := client.MintToken(
		context.Background(),
		signer,
		packageID,
		"testcoin",
		treasuryCap,
		mintAmount,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnRes.Effects.Data.V1.Status.Status)

	return packageID, treasuryCap
}
