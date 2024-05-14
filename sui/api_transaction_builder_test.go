package sui_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"
	"github.com/stretchr/testify/require"
)

func TestBatchTransaction(t *testing.T) {
	t.Log("TestBatchTransaction TODO")
	// api := sui.NewSuiClient(conn.DevnetEndpointUrl)

	// txnBytes, err := api.BatchTransaction(context.Background(), signer, *coin1, *coin2, nil, 10000)
	// require.NoError(t, err)
	// dryRunTxn(t, api, txnBytes, M1Account(t))
}

func TestMergeCoins(t *testing.T) {
	t.Skip("FIXME create an account has at least two coin objects on chain")
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txn, err := api.MergeCoins(
		context.Background(), signer,
		coin1.CoinObjectID, coin2.CoinObjectID,
		coin3.CoinObjectID, coin3.Balance,
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestMoveCall(t *testing.T) {
	client, signer := sui.NewSuiClient(conn.LocalnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_MNEMONIC)
	// client := sui.NewSuiClient(conn.TestnetEndpointUrl)
	// signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	// require.NoError(t, err)

	// directly build (need sui toolchain)
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/sdk_tests/")
	require.NoError(t, err)
	// jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/sdk_tests/contract_base64.json")
	// require.NoError(t, err)

	// var modules utils.CompiledMoveModules
	// err = json.Unmarshal(jsonData, &modules)
	// require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		sui_signer.TEST_ADDRESS,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
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
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)

	packageID := txnResponse.GetPublishedPackageID()
	fmt.Println("packageID: ", packageID)
	type args struct {
		module    string
		function  string
		typeArgs  []string
		arguments []sui.SuiJsonArg
	}
	tests := []struct {
		name        string
		client      *sui.ImplSuiAPI
		args        args
		wantErr     error
		wantErrExec error
		afterFunc   func(t *testing.T, res *models.SuiTransactionBlockResponse)
	}{
		{
			name:   "byte_array_of_arrays",
			client: client,
			args: args{
				"sdk_tests",
				"input_byte_array_of_arrays",
				[]string{},
				// []sui.SuiJsonArg{sui.ToSuiJsonArg([]string{"haha", "abc"})}, // directly pass string array works
				// []sui.SuiJsonArg{sui.ToSuiJsonArg([]string{"68616861", "616263"})}, // encode each []byte to hex string works
				// []sui.SuiJsonArg{sui.ToSuiJsonArg([][]byte{[]byte{104, 97, 104, 97}, []byte{97, 98, 99}})},
				[]sui.SuiJsonArg{[][]byte{[]byte{104, 97, 104, 97}, []byte{97, 98, 99}}},
			},
			wantErr: nil,
			afterFunc: func(t *testing.T, res *models.SuiTransactionBlockResponse) {
				type InputByteArrayOfArrays struct {
					Data [][]byte
				}
				queryEventsRes, err := client.QueryEvents(
					context.Background(),
					&models.EventFilter{
						Transaction: &txnResponse.Digest,
					},
					nil,
					nil,
					false,
				)
				require.NoError(t, err)
				b, err := json.Marshal(queryEventsRes.Data[0].ParsedJson.(map[string]interface{}))
				require.NoError(t, err)
				var eventRes InputByteArrayOfArrays
				err = json.Unmarshal(b, &eventRes)
				require.NoError(t, err)
				fmt.Println("eventRes.Data: ", string(eventRes.Data[0]))
				fmt.Println("eventRes.Data: ", string(eventRes.Data[1]))
				require.Equal(t, []byte("haha"), eventRes.Data[0])
				require.Equal(t, []byte("abc"), eventRes.Data[1])
			},
		},
		// {
		// 	name:   "ints",
		// 	client: client,
		// 	args: args{
		// 		"sdk_tests",
		// 		"input_ints",
		// 		[]string{},
		// 		[]sui.SuiJsonArg{sui.ToSuiJsonArg(uint8(12)), sui.ToSuiJsonArg(uint16(511)), sui.ToSuiJsonArg(uint32(80000)), sui.ToSuiJsonArg(uint64(12))},
		// 	},
		// 	wantErr: nil,
		// 	afterFunc: func(t *testing.T, res *models.SuiTransactionBlockResponse) {
		// 		type InputByteArrayOfArrays struct {
		// 			Input0 uint8
		// 			Input1 uint16
		// 			Input2 uint32
		// 			Input3 uint64
		// 		}
		// 		queryEventsRes, err := client.QueryEvents(
		// 			context.Background(),
		// 			&models.EventFilter{
		// 				Transaction: &txnResponse.Digest,
		// 			},
		// 			nil,
		// 			nil,
		// 			false,
		// 		)
		// 		require.NoError(t, err)
		// 		b, err := json.Marshal(queryEventsRes.Data[0].ParsedJson.(map[string]interface{}))
		// 		require.NoError(t, err)
		// 		var eventRes InputByteArrayOfArrays
		// 		err = json.Unmarshal(b, &eventRes)
		// 		require.NoError(t, err)

		// 		require.Equal(t, uint8(12), eventRes.Input0)
		// 		require.Equal(t, uint16(511), eventRes.Input1)
		// 		require.Equal(t, uint32(80000), eventRes.Input2)
		// 		require.Equal(t, uint64(12), eventRes.Input3)
		// 	},
		// },
	}
	for _, tt := range tests {
		fmt.Println("string: ", []sui.SuiJsonArg{sui.ToSuiJsonArg([]string{"haha", "abc"})})
		fmt.Println("bytes: ", []sui.SuiJsonArg{[]string{"haha", "abc"}})
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.client.MoveCall(
					context.Background(),
					signer.Address,
					packageID,
					tt.args.module,
					tt.args.function,
					tt.args.typeArgs,
					tt.args.arguments,
					nil,
					models.NewSafeSuiBigInt(uint64(sui.DefaultGasBudget)),
				)
				require.ErrorIs(t, err, tt.wantErr)

				txnResponse, err = client.SignAndExecuteTransaction(
					context.Background(),
					signer,
					got.TxBytes,
					&models.SuiTransactionBlockResponseOptions{
						ShowEffects:       true,
						ShowObjectChanges: true,
					},
				)
				require.ErrorIs(t, err, tt.wantErrExec)
				// FIXME this should be false for the test cases which should be failed
				require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)

				tt.afterFunc(t, txnResponse)
			},
		)
	}
}

func TestPay(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	limit := len(coins.Data) - 1 // need reserve a coin for gas

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, limit, 0)
	require.NoError(t, err)

	txn, err := api.Pay(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amount),
		},
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPayAllSui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := api.PayAllSui(
		context.Background(), signer, recipient,
		pickedCoins.CoinIds(),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPaySui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 0, 0)
	require.NoError(t, err)

	txn, err := api.PaySui(
		context.Background(), signer,
		pickedCoins.CoinIds(),
		[]*sui_types.SuiAddress{recipient},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amount),
		},
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestPublish(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	_, err = sui.RequestFundFromFaucet(signer.Address, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	// If local side has installed Sui-cli then the user can use the following func to build move contracts
	// modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin")
	// require.NoError(t, err)

	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/testcoin/contract_base64.json")
	require.NoError(t, err)

	var modules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &modules)
	require.NoError(t, err)

	coins, err := client.GetCoins(context.Background(), signer.Address, nil, nil, 10)
	require.NoError(t, err)
	gasBudget := uint64(1000000000)
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(1000000), gasBudget, 10, 10)
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		signer.Address,
		modules.Modules,
		modules.Dependencies,
		pickedCoins.CoinIds()[0],
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	txnResponse, err := client.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, &models.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	})
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)
}

func TestSplitCoin(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)
	splitCoins := []models.SafeSuiBigInt[uint64]{models.NewSafeSuiBigInt(amount / 2)}

	txn, err := api.SplitCoin(
		context.Background(), signer,
		pickedCoins.Coins[0].CoinObjectID,
		splitCoins,
		nil, models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, false)
}

func TestSplitCoinEqual(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.01).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), 0, 1, 0)
	require.NoError(t, err)

	txn, err := api.SplitCoinEqual(
		context.Background(), signer,
		pickedCoins.Coins[0].CoinObjectID,
		models.NewSafeSuiBigInt(uint64(2)),
		nil, models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestTransferObject(t *testing.T) {
	t.Skip("FIXME create an account has at least two coin objects on chain")
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 2)
	coin := coins.Data[0]

	txn, err := api.TransferObject(
		context.Background(),
		signer,
		recipient,
		coin.CoinObjectID,
		nil,
		models.NewSafeSuiBigInt(sui_types.SUI(0.01).Uint64()),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}

func TestTransferSui(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	signer := sui_signer.TEST_ADDRESS
	recipient := signer
	coins, err := api.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)

	amount := sui_types.SUI(0.0001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(0).SetUint64(amount), gasBudget, 1, 0)
	require.NoError(t, err)

	txn, err := api.TransferSui(
		context.Background(), signer, recipient,
		pickedCoins.Coins[0].CoinObjectID,
		models.NewSafeSuiBigInt(amount),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)

	dryRunTxn(t, api, txn.TxBytes, true)
}
