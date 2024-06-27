package sui_types_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"

	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
)

func TestPTBMoveCall(t *testing.T) {
	t.Run(
		"access_multiple_return_values_from_move_func", func(t *testing.T) {
			client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
			var modules utils.CompiledMoveModules
			data, err := os.ReadFile(utils.GetGitRoot() + "/contracts/sdk_verify/contract_base64.json")
			require.NoError(t, err)
			err = json.Unmarshal(data, &modules)
			require.NoError(t, err)
			_, packageID, err := client.PublishContract(
				context.Background(),
				sender,
				modules.Modules,
				modules.Dependencies,
				sui.DefaultGasBudget,
				&models.SuiTransactionBlockResponseOptions{ShowObjectChanges: true, ShowEffects: true},
			)
			require.NoError(t, err)

			coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
				Owner: sender.Address,
				Limit: 3,
			})
			require.NoError(t, err)
			coins := models.Coins(coinPages.Data)

			ptb := sui_types.NewProgrammableTransactionBuilder()
			require.NoError(t, err)

			ptb.Command(
				sui_types.Command{
					MoveCall: &sui_types.ProgrammableMoveCall{
						Package:       packageID,
						Module:        "sdk_verify",
						Function:      "ret_two_1",
						TypeArguments: []sui_types.TypeTag{},
						Arguments:     []sui_types.Argument{},
					},
				},
			)
			ptb.Command(
				sui_types.Command{
					MoveCall: &sui_types.ProgrammableMoveCall{
						Package:       packageID,
						Module:        "sdk_verify",
						Function:      "ret_two_2",
						TypeArguments: []sui_types.TypeTag{},
						Arguments: []sui_types.Argument{
							{NestedResult: &sui_types.NestedResult{Cmd: 0, Result: 1}},
							{NestedResult: &sui_types.NestedResult{Cmd: 0, Result: 0}},
						},
					},
				},
			)
			pt := ptb.Finish()
			txData := sui_types.NewProgrammable(
				sender.Address,
				pt,
				[]*sui_types.ObjectRef{coins[0].Ref()},
				sui.DefaultGasBudget,
				sui.DefaultGasPrice,
			)
			txBytes, err := bcs.Marshal(txData)
			require.NoError(t, err)
			simulate, err := client.DryRunTransaction(context.Background(), txBytes)
			require.NoError(t, err)

			require.Empty(t, simulate.Effects.Data.V1.Status.Error)
			require.True(t, simulate.Effects.Data.IsSuccess())
			require.Equal(t, coins[0].CoinObjectID, simulate.Effects.Data.V1.GasObject.Reference.ObjectID)
		},
	)
}

func TestPTBTransferObject(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 2,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)
	gasCoin := coins[0]
	transferCoin := coins[1]

	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferObject(recipient.Address, transferCoin.Ref())
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		[]*sui_types.ObjectRef{gasCoin.Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.TransferObject(
		context.Background(),
		&models.TransferObjectRequest{
			Signer:    sender.Address,
			Recipient: recipient.Address,
			ObjectID:  transferCoin.CoinObjectID,
			Gas:       gasCoin.CoinObjectID,
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBTransferSui(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 1,
	})
	require.NoError(t, err)
	coin := models.Coins(coinPages.Data)[0]
	amount := uint64(123)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferSui(recipient.Address, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		[]*sui_types.ObjectRef{coin.Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.TransferSui(
		context.Background(),
		&models.TransferSuiRequest{
			Signer:    sender.Address,
			Recipient: recipient.Address,
			ObjectID:  coin.CoinObjectID,
			Amount:    models.NewBigInt(amount),
			GasBudget: models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTBPayAllSui(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 3,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PayAllSui(recipient.Address)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		coins.CoinRefs(),
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.PayAllSui(
		context.Background(),
		&models.PayAllSuiRequest{
			Signer:     sender.Address,
			Recipient:  recipient.Address,
			InputCoins: coins.ObjectIDs(),
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPaySui(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient1 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	_, recipient2 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 2)
	coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 1,
	})
	require.NoError(t, err)
	coin := coinPages.Data[0]

	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PaySui(
		[]*sui_types.SuiAddress{recipient1.Address, recipient2.Address},
		[]uint64{123, 456},
	)
	require.NoError(t, err)
	pt := ptb.Finish()

	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		[]*sui_types.ObjectRef{
			coin.Ref(),
		},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	require.Equal(t, coin.CoinObjectID, simulate.Effects.Data.V1.GasObject.Reference.ObjectID)

	// 1 for Mutated, 2 created (the 2 transfer in pay_sui pt),
	require.Len(t, simulate.ObjectChanges, 3)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, coin.CoinObjectID, &change.Data.Mutated.ObjectID)
		} else if change.Data.Created != nil {
			require.Contains(t, []*sui_types.SuiAddress{recipient1.Address, recipient2.Address}, change.Data.Created.Owner.AddressOwner)
		}
	}

	// build with remote rpc
	txn, err := client.PaySui(
		context.Background(),
		&models.PaySuiRequest{
			Signer:     sender.Address,
			InputCoins: []*sui_types.ObjectID{coin.CoinObjectID},
			Recipients: []*sui_types.SuiAddress{recipient1.Address, recipient2.Address},
			Amount:     []*models.BigInt{models.NewBigInt(123), models.NewBigInt(456)},
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPay(t *testing.T) {
	client, sender := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, recipient1 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	_, recipient2 := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 2)
	coinPages, err := client.GetCoins(context.Background(), &models.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 3,
	})
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)
	gasCoin := coins[0] // save the 1st element for gas fee
	transferCoins := coins[1:]
	amounts := []uint64{123, 567}
	totalBal := coins.TotalBalance().Uint64()

	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.Pay(
		transferCoins.CoinRefs(),
		[]*sui_types.SuiAddress{recipient1.Address, recipient2.Address},
		[]uint64{amounts[0], amounts[1]},
	)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		sender.Address,
		pt,
		[]*sui_types.ObjectRef{
			gasCoin.Ref(),
		},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	require.Equal(t, gasCoin.CoinObjectID, simulate.Effects.Data.V1.GasObject.Reference.ObjectID)

	// 2 for Mutated (1 gas coin and 1 merged coin in pay pt), 2 created (the 2 transfer in pay pt),
	require.Len(t, simulate.ObjectChanges, 5)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Contains(
				t,
				[]*sui_types.ObjectID{gasCoin.CoinObjectID, transferCoins[0].CoinObjectID},
				&change.Data.Mutated.ObjectID,
			)
		} else if change.Data.Deleted != nil {
			require.Equal(t, transferCoins[1].CoinObjectID, &change.Data.Deleted.ObjectID)
		}
	}
	require.Len(t, simulate.BalanceChanges, 3)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == sender.Address {
			require.Equal(t, totalBal-(amounts[0]+amounts[1]), balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipient1.Address {
			require.Equal(t, amounts[0], balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipient2.Address {
			require.Equal(t, amounts[1], balChange.Amount)
		}
	}

	// build with remote rpc
	txn, err := client.Pay(
		context.Background(),
		&models.PayRequest{
			Signer:     sender.Address,
			InputCoins: transferCoins.ObjectIDs(),
			Recipients: []*sui_types.SuiAddress{recipient1.Address, recipient2.Address},
			Amount:     []*models.BigInt{models.NewBigInt(amounts[0]), models.NewBigInt(amounts[1])},
			Gas:        gasCoin.CoinObjectID,
			GasBudget:  models.NewBigInt(sui.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}
