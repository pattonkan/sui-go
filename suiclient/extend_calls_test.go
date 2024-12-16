package suiclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"

	"github.com/stretchr/testify/require"
)

func TestMintToken(t *testing.T) {
	client, signer := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)

	// module name is 'testcoin'
	tokenPackageId, treasuryCap := deployTestcoin(t, client, signer)
	mintAmount := uint64(1000000)
	txnRes, err := client.MintToken(
		context.Background(),
		signer,
		tokenPackageId,
		"testcoin",
		treasuryCap,
		mintAmount,
		&suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnRes.Effects.Data.IsSuccess())
	coinType := fmt.Sprintf("%s::testcoin::TESTCOIN", tokenPackageId.String())

	// all the minted tokens were sent to the signer, so we should find a single object contains all the minted token
	coins, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner:    signer.Address,
		CoinType: &coinType,
		Limit:    10,
	})
	require.NoError(t, err)
	require.Equal(t, mintAmount, coins.Data[0].Balance.Uint64())
}

func deployTestcoin(t *testing.T, client *suiclient.ClientImpl, signer *suisigner.Signer) (
	*sui.PackageId,
	*sui.ObjectId,
) {
	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/testcoin/contract_base64.json")
	require.NoError(t, err)

	var modules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &modules)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		&suiclient.PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(suiclient.DefaultGasBudget * 10),
		},
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), signer, txnBytes.TxBytes, &suiclient.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageId, err := txnResponse.GetPublishedPackageId()
	require.NoError(t, err)

	treasuryCap, _, err := txnResponse.GetCreatedObjectInfo("coin", "TreasuryCap")
	require.NoError(t, err)

	return packageId, treasuryCap
}

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	api := suiclient.NewClient(conn.DevnetEndpointUrl)

	options := suiclient.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
	}
	coinType := fmt.Sprintf("0x2::coin::Coin<%v>", sui.SuiCoinType)
	filterObject, err := api.BatchGetObjectsOwnedByAddress(context.TODO(), suisigner.TEST_ADDRESS, &options, coinType)
	require.NoError(t, err)
	t.Log(filterObject)
}
