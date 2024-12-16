package suiptb_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"
)

func TestPTBMoveCall(t *testing.T) {
	t.Run(
		"access_multiple_return_values_from_move_func", func(t *testing.T) {
			client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
			var modules utils.CompiledMoveModules
			data, err := os.ReadFile(utils.GetGitRoot() + "/contracts/sdk_verify/contract_base64.json")
			require.NoError(t, err)
			err = json.Unmarshal(data, &modules)
			require.NoError(t, err)
			_, packageId, err := client.PublishContract(
				context.Background(),
				sender,
				modules.Modules,
				modules.Dependencies,
				suiclient.DefaultGasBudget,
				&suiclient.SuiTransactionBlockResponseOptions{ShowObjectChanges: true, ShowEffects: true},
			)
			require.NoError(t, err)

			coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
				Owner: sender.Address,
				Limit: 3,
			})
			require.NoError(t, err)
			coins := suiclient.Coins(coinPages.Data)

			ptb := suiptb.NewTransactionDataTransactionBuilder()
			require.NoError(t, err)

			ptb.Command(
				suiptb.Command{
					MoveCall: &suiptb.ProgrammableMoveCall{
						Package:       packageId,
						Module:        "sdk_verify",
						Function:      "ret_two_1",
						TypeArguments: []sui.TypeTag{},
						Arguments:     []suiptb.Argument{},
					},
				},
			)
			ptb.Command(
				suiptb.Command{
					MoveCall: &suiptb.ProgrammableMoveCall{
						Package:       packageId,
						Module:        "sdk_verify",
						Function:      "ret_two_2",
						TypeArguments: []sui.TypeTag{},
						Arguments: []suiptb.Argument{
							{NestedResult: &suiptb.NestedResult{Cmd: 0, Result: 1}},
							{NestedResult: &suiptb.NestedResult{Cmd: 0, Result: 0}},
						},
					},
				},
			)
			pt := ptb.Finish()
			txData := suiptb.NewTransactionData(
				sender.Address,
				pt,
				[]*sui.ObjectRef{coins[0].Ref()},
				suiclient.DefaultGasBudget,
				suiclient.DefaultGasPrice,
			)
			txBytes, err := bcs.Marshal(txData)
			require.NoError(t, err)
			simulate, err := client.DryRunTransaction(context.Background(), txBytes)
			require.NoError(t, err)

			require.Empty(t, simulate.Effects.Data.V1.Status.Error)
			require.True(t, simulate.Effects.Data.IsSuccess())
			require.Equal(t, coins[0].CoinObjectId, simulate.Effects.Data.V1.GasObject.Reference.ObjectId)
		},
	)

	t.Run(
		"option<T> arguments", func(t *testing.T) {
			client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
			var modules utils.CompiledMoveModules
			data, err := os.ReadFile(utils.GetGitRoot() + "/contracts/sdk_verify/contract_base64.json")
			require.NoError(t, err)
			err = json.Unmarshal(data, &modules)
			require.NoError(t, err)
			_, packageId, err := client.PublishContract(
				context.Background(),
				sender,
				modules.Modules,
				modules.Dependencies,
				suiclient.DefaultGasBudget,
				&suiclient.SuiTransactionBlockResponseOptions{ShowObjectChanges: true, ShowEffects: true},
			)
			require.NoError(t, err)

			coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
				Owner: sender.Address,
				Limit: 3,
			})
			require.NoError(t, err)
			coins := suiclient.Coins(coinPages.Data)

			ptb := suiptb.NewTransactionDataTransactionBuilder()
			require.NoError(t, err)

			ptb.Command(
				suiptb.Command{
					MoveCall: &suiptb.ProgrammableMoveCall{
						Package:       packageId,
						Module:        "sdk_verify",
						Function:      "option_args",
						TypeArguments: []sui.TypeTag{},
						Arguments: []suiptb.Argument{
							ptb.MustPure(&bcs.Option[[]byte]{Some: []byte{1, 2}}),
							ptb.MustPure(&bcs.Option[uint32]{None: true}),
						},
					},
				},
			)
			pt := ptb.Finish()
			txData := suiptb.NewTransactionData(
				sender.Address,
				pt,
				[]*sui.ObjectRef{coins[0].Ref()},
				suiclient.DefaultGasBudget,
				suiclient.DefaultGasPrice,
			)
			txBytes, err := bcs.Marshal(txData)
			require.NoError(t, err)
			simulate, err := client.DryRunTransaction(context.Background(), txBytes)
			require.NoError(t, err)

			require.Empty(t, simulate.Effects.Data.V1.Status.Error)
			require.True(t, simulate.Effects.Data.IsSuccess())
			require.Equal(t, coins[0].CoinObjectId, simulate.Effects.Data.V1.GasObject.Reference.ObjectId)
		},
	)
}

func TestPTBTransferObject(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 2,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)
	gasCoin := coins[0]
	transferCoin := coins[1]

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	err = ptb.TransferObject(recipient.Address, transferCoin.Ref())
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		[]*sui.ObjectRef{gasCoin.Ref()},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.TransferObject(
		context.Background(),
		&suiclient.TransferObjectRequest{
			Signer:    sender.Address,
			Recipient: recipient.Address,
			ObjectId:  transferCoin.CoinObjectId,
			Gas:       gasCoin.CoinObjectId,
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBTransferSui(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 1,
	})
	require.NoError(t, err)
	coin := suiclient.Coins(coinPages.Data)[0]
	amount := uint64(123)

	// build with BCS
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	err = ptb.TransferSui(recipient.Address, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		[]*sui.ObjectRef{coin.Ref()},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.TransferSui(
		context.Background(),
		&suiclient.TransferSuiRequest{
			Signer:    sender.Address,
			Recipient: recipient.Address,
			ObjectId:  coin.CoinObjectId,
			Amount:    sui.NewBigInt(amount),
			GasBudget: sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTBPayAllSui(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 3,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)

	// build with BCS
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	err = ptb.PayAllSui(recipient.Address)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		coins.CoinRefs(),
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := client.PayAllSui(
		context.Background(),
		&suiclient.PayAllSuiRequest{
			Signer:     sender.Address,
			Recipient:  recipient.Address,
			InputCoins: coins.ObjectIds(),
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPaySui(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient1 := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	_, recipient2 := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 2)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 1,
	})
	require.NoError(t, err)
	coin := coinPages.Data[0]

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	err = ptb.PaySui(
		[]*sui.Address{recipient1.Address, recipient2.Address},
		[]uint64{123, 456},
	)
	require.NoError(t, err)
	pt := ptb.Finish()

	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		[]*sui.ObjectRef{
			coin.Ref(),
		},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	require.Equal(t, coin.CoinObjectId, simulate.Effects.Data.V1.GasObject.Reference.ObjectId)

	// 1 for Mutated, 2 created (the 2 transfer in pay_sui pt),
	require.Len(t, simulate.ObjectChanges, 3)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, coin.CoinObjectId, &change.Data.Mutated.ObjectId)
		} else if change.Data.Created != nil {
			require.Contains(t, []*sui.Address{recipient1.Address, recipient2.Address}, change.Data.Created.Owner.AddressOwner)
		}
	}

	// build with remote rpc
	txn, err := client.PaySui(
		context.Background(),
		&suiclient.PaySuiRequest{
			Signer:     sender.Address,
			InputCoins: []*sui.ObjectId{coin.CoinObjectId},
			Recipients: []*sui.Address{recipient1.Address, recipient2.Address},
			Amount:     []*sui.BigInt{sui.NewBigInt(123), sui.NewBigInt(456)},
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPay(t *testing.T) {
	client, sender := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, recipient1 := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	_, recipient2 := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 2)
	coinPages, err := client.GetCoins(context.Background(), &suiclient.GetCoinsRequest{
		Owner: sender.Address,
		Limit: 3,
	})
	require.NoError(t, err)
	coins := suiclient.Coins(coinPages.Data)
	gasCoin := coins[0] // save the 1st element for gas fee
	transferCoins := coins[1:]
	amounts := []uint64{123, 567}
	totalBal := coins.TotalBalance().Uint64()

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	err = ptb.Pay(
		transferCoins.CoinRefs(),
		[]*sui.Address{recipient1.Address, recipient2.Address},
		[]uint64{amounts[0], amounts[1]},
	)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		sender.Address,
		pt,
		[]*sui.ObjectRef{
			gasCoin.Ref(),
		},
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Empty(t, simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	require.Equal(t, gasCoin.CoinObjectId, simulate.Effects.Data.V1.GasObject.Reference.ObjectId)

	// 2 for Mutated (1 gas coin and 1 merged coin in pay pt), 2 created (the 2 transfer in pay pt),
	require.Len(t, simulate.ObjectChanges, 5)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Contains(
				t,
				[]*sui.ObjectId{gasCoin.CoinObjectId, transferCoins[0].CoinObjectId},
				&change.Data.Mutated.ObjectId,
			)
		} else if change.Data.Deleted != nil {
			require.Equal(t, transferCoins[1].CoinObjectId, &change.Data.Deleted.ObjectId)
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
		&suiclient.PayRequest{
			Signer:     sender.Address,
			InputCoins: transferCoins.ObjectIds(),
			Recipients: []*sui.Address{recipient1.Address, recipient2.Address},
			Amount:     []*sui.BigInt{sui.NewBigInt(amounts[0]), sui.NewBigInt(amounts[1])},
			Gas:        gasCoin.CoinObjectId,
			GasBudget:  sui.NewBigInt(suiclient.DefaultGasBudget),
		},
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}
