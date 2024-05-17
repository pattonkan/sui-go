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

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
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
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID := txnResponse.GetPublishedPackageID()
	t.Log("packageID: ", packageID)

	startNewChainRes, err := client.StartNewChain(
		context.Background(),
		signer,
		packageID,
		sui.DefaultGasBudget,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, startNewChainRes.Effects.Data.IsSuccess())
	t.Logf("StartNewChain response: %#v\n", startNewChainRes)
}

func TestSendCoin(t *testing.T) {
	t.Skip("only for localnet")
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := isc.BuildAndDeployIscContracts(t, client, signer)
	tokenPackageID, _ := isc.BuildDeployMintTestcoin(t, client, signer)

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
	require.True(t, startNewChainRes.Effects.Data.IsSuccess())

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
	require.True(t, sendCoinRes.Effects.Data.IsSuccess())

	getObjectRes, err := client.GetObject(context.Background(), coins.Data[0].CoinObjectID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, anchorObjID.String(), getObjectRes.Data.Owner.ObjectOwnerInternal.AddressOwner.String())
}

func TestReceiveCoin(t *testing.T) {
	t.Skip("only for localnet")
	var err error
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := isc.BuildAndDeployIscContracts(t, client, signer)
	tokenPackageID, _ := isc.BuildDeployMintTestcoin(t, client, signer)

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
	require.True(t, startNewChainRes.Effects.Data.IsSuccess())

	var assets1, assets2 *isc.Assets
	for _, change := range startNewChainRes.ObjectChanges {
		if change.Data.Created != nil {
			assets1, err = client.GetAssets(context.Background(), iscPackageID, &change.Data.Created.ObjectID)
			require.NoError(t, err)
			require.Len(t, assets1.Coins, 0)
		}
	}

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
	require.True(t, sendCoinRes.Effects.Data.IsSuccess())

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
	require.True(t, receiveCoinRes.Effects.Data.IsSuccess())
	assets2, err = client.GetAssets(context.Background(), iscPackageID, anchorObjID)
	require.NoError(t, err)
	require.Len(t, assets2.Coins, 1)
}

func TestCreateRequest(t *testing.T) {
	t.Skip("only for localnet")
	var err error
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := isc.BuildAndDeployIscContracts(t, client, signer)

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

	anchorObjID, _, err := sui.GetCreatedObjectIdAndType(startNewChainRes, "anchor", "Anchor")
	require.NoError(t, err)

	createReqRes, err := client.CreateRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		"isc_test_contract_name",
		"isc_test_func_name",
		[][]byte{}, // func input
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		})
	require.NoError(t, err)
	require.True(t, createReqRes.Effects.Data.IsSuccess())

	_, _, err = sui.GetCreatedObjectIdAndType(createReqRes, "request", "Request")
	require.NoError(t, err)
}

func TestSendRequest(t *testing.T) {
	t.Skip("only for localnet")
	var err error
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := isc.BuildAndDeployIscContracts(t, client, signer)

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

	anchorObjID, _, err := sui.GetCreatedObjectIdAndType(startNewChainRes, "anchor", "Anchor")
	require.NoError(t, err)

	createReqRes, err := client.CreateRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		"isc_test_contract_name",
		"isc_test_func_name",
		[][]byte{}, // func input
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, createReqRes.Effects.Data.IsSuccess())

	reqObjID, _, err := sui.GetCreatedObjectIdAndType(createReqRes, "request", "Request")
	require.NoError(t, err)
	getObjectRes, err := client.GetObject(context.Background(), reqObjID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, signer.Address, getObjectRes.Data.Owner.AddressOwner)

	sendReqRes, err := client.SendRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		reqObjID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, sendReqRes.Effects.Data.IsSuccess())

	getObjectRes, err = client.GetObject(context.Background(), reqObjID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, anchorObjID, getObjectRes.Data.Owner.AddressOwner)
}

func TestReceiveRequest(t *testing.T) {
	t.Skip("only for localnet")
	var err error
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	err = sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)

	iscPackageID := isc.BuildAndDeployIscContracts(t, client, signer)

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

	anchorObjID, _, err := sui.GetCreatedObjectIdAndType(startNewChainRes, "anchor", "Anchor")
	require.NoError(t, err)

	createReqRes, err := client.CreateRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		"isc_test_contract_name", // FIXME set up the proper ISC target contract name
		"isc_test_func_name",     // FIXME set up the proper ISC target func name
		[][]byte{},               // func input
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, createReqRes.Effects.Data.IsSuccess())

	reqObjID, _, err := sui.GetCreatedObjectIdAndType(createReqRes, "request", "Request")
	require.NoError(t, err)
	getObjectRes, err := client.GetObject(context.Background(), reqObjID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, signer.Address, getObjectRes.Data.Owner.AddressOwner)

	sendReqRes, err := client.SendRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		reqObjID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, sendReqRes.Effects.Data.IsSuccess())

	getObjectRes, err = client.GetObject(context.Background(), reqObjID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.Equal(t, anchorObjID, getObjectRes.Data.Owner.AddressOwner)

	receiveReqRes, err := client.ReceiveRequest(
		context.Background(),
		signer,
		iscPackageID,
		anchorObjID,
		reqObjID,
		sui.DefaultGasBudget, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, receiveReqRes.Effects.Data.IsSuccess())

	getObjectRes, err = client.GetObject(context.Background(), reqObjID, &models.SuiObjectDataOptions{ShowOwner: true})
	require.NoError(t, err)
	require.NotNil(t, getObjectRes.Error.Data.Deleted)
}
