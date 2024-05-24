package sui_types_test

import (
	"context"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/stretchr/testify/require"
)

func TestPTBMoveCall(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	// a random validator on devnet
	validatorAddress, err := sui_types.SuiAddressFromHex("0x1571d9d389181f509c3fd2bb0f0eb06dc8ba73152188567878a5dc85097a0eab")
	require.NoError(t, err)

	// case 1: split target amount
	ptb1 := sui_types.NewProgrammableTransactionBuilder()
	arg0, err := ptb1.Obj(sui_types.SuiSystemMutObj)
	require.NoError(t, err)
	// if sui is not enough, the transaction will fail
	splitAmountArg := ptb1.MustPure(uint64(1e9))
	arg1 := ptb1.Command(sui_types.Command{
		SplitCoins: &sui_types.ProgrammableSplitCoins{
			Coin:    sui_types.Argument{GasCoin: &serialization.EmptyEnum{}},
			Amounts: []sui_types.Argument{splitAmountArg},
		}},
	)
	arg2 := ptb1.MustPure(validatorAddress)
	ptb1.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:       sui_types.SuiPackageIdSuiSystem,
			Module:        sui_types.SuiSystemModuleName,
			Function:      sui_types.AddStakeFunName,
			TypeArguments: []sui_types.TypeTag{},
			Arguments:     []sui_types.Argument{arg0, arg1, arg2},
		}},
	)
	pt1 := ptb1.Finish()
	tx1 := sui_types.NewProgrammable(
		sender.Address,
		pt1,
		[]*sui_types.ObjectRef{coins[0].Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)
	txBytes1, err := bcs.Marshal(tx1)
	require.NoError(t, err)
	simulate1, err := client.DryRunTransaction(context.Background(), txBytes1)
	require.NoError(t, err)
	require.Empty(t, simulate1.Effects.Data.V1.Status.Error)
	require.True(t, simulate1.Effects.Data.IsSuccess())
	require.Equal(t, coins[0].CoinObjectID.String(), simulate1.Effects.Data.V1.GasObject.Reference.ObjectID)

	// case 2: direct stake the specified coin
	ptb2 := sui_types.NewProgrammableTransactionBuilder()
	coinArg := sui_types.CallArg{
		Object: &sui_types.ObjectArg{
			ImmOrOwnedObject: coins[1].Ref(),
		},
	}
	addrBytes := validatorAddress.Data()
	addrArg := sui_types.CallArg{
		Pure: &addrBytes,
	}
	err = ptb2.MoveCall(
		sui_types.SuiPackageIdSuiSystem,
		sui_types.SuiSystemModuleName,
		sui_types.AddStakeFunName,
		[]sui_types.TypeTag{},
		[]sui_types.CallArg{sui_types.SuiSystemMut, coinArg, addrArg},
	)
	require.NoError(t, err)
	pt2 := ptb2.Finish()
	tx2 := sui_types.NewProgrammable(
		sender.Address,
		pt2,
		[]*sui_types.ObjectRef{coins[0].Ref()},
		sui.DefaultGasBudget,
		sui.DefaultGasPrice,
	)

	txBytes2, err := bcs.Marshal(tx2)
	require.NoError(t, err)
	simulate2, err := client.DryRunTransaction(context.Background(), txBytes2)
	require.NoError(t, err)
	require.Empty(t, simulate2.Effects.Data.V1.Status.Error)
	require.True(t, simulate2.Effects.Data.IsSuccess())
	require.Equal(t, coins[0].CoinObjectID.String(), simulate2.Effects.Data.V1.GasObject.Reference.ObjectID)
}

func TestPTBTransferObject(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipient := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 1)[0]
	coinType := models.SuiCoinType
	limit := uint(2)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
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
		sender.Address,
		recipient.Address,
		transferCoin.CoinObjectID,
		gasCoin.CoinObjectID,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBTransferSui(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipient := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 1)[0]
	coinType := models.SuiCoinType
	limit := uint(1)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coin := models.Coins(coinPages.Data)[0]
	amount := uint64(123)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	ptb.TransferSui(recipient.Address, &amount)
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
		sender.Address,
		recipient.Address,
		coin.CoinObjectID,
		models.NewSafeSuiBigInt(amount),
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTBPayAllSui(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipient := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 1)[0]
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	ptb.PayAllSui(recipient.Address)
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
		sender.Address,
		recipient.Address,
		coins.ObjectIDs(),
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPaySui(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipients := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 2)
	coinType := models.SuiCoinType
	limit := uint(1)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coin := coinPages.Data[0]

	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PaySui(
		[]*sui_types.SuiAddress{recipients[0].Address, recipients[1].Address},
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
	require.Equal(t, coin.CoinObjectID.String(), simulate.Effects.Data.V1.GasObject.Reference.ObjectID)

	// 1 for Mutated, 2 created (the 2 transfer in pay_sui pt),
	require.Len(t, simulate.ObjectChanges, 3)
	for _, change := range simulate.ObjectChanges {
		if change.Data.Mutated != nil {
			require.Equal(t, coin.CoinObjectID, &change.Data.Mutated.ObjectID)
		} else if change.Data.Created != nil {
			require.Contains(t, []*sui_types.SuiAddress{recipients[0].Address, recipients[1].Address}, change.Data.Created.Owner.AddressOwner)
		}
	}

	// build with remote rpc
	txn, err := client.PaySui(
		context.Background(),
		sender.Address,
		[]*sui_types.ObjectID{coin.CoinObjectID},
		[]*sui_types.SuiAddress{recipients[0].Address, recipients[1].Address},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(uint64(123)),
			models.NewSafeSuiBigInt(uint64(456)),
		},
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}

func TestPTBPay(t *testing.T) {
	client, sender := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipients := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 2)
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), sender.Address, &coinType, nil, limit)
	require.NoError(t, err)
	coins := models.Coins(coinPages.Data)
	gasCoin := coins[0] // save the 1st element for gas fee
	transferCoins := coins[1:]
	amounts := []uint64{123, 567}
	totalBal := coins.TotalBalance().Uint64()

	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.Pay(
		transferCoins.CoinRefs(),
		[]*sui_types.SuiAddress{recipients[0].Address, recipients[1].Address},
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
	require.Equal(t, gasCoin.CoinObjectID.String(), simulate.Effects.Data.V1.GasObject.Reference.ObjectID)

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
		} else if balChange.Owner.AddressOwner == recipients[0].Address {
			require.Equal(t, amounts[0], balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipients[1].Address {
			require.Equal(t, amounts[1], balChange.Amount)
		}
	}

	// build with remote rpc
	txn, err := client.Pay(
		context.Background(),
		sender.Address,
		transferCoins.ObjectIDs(),
		[]*sui_types.SuiAddress{recipients[0].Address, recipients[1].Address},
		[]models.SafeSuiBigInt[uint64]{
			models.NewSafeSuiBigInt(amounts[0]),
			models.NewSafeSuiBigInt(amounts[1]),
		},
		gasCoin.CoinObjectID,
		models.NewSafeSuiBigInt(sui.DefaultGasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()
	require.Equal(t, txBytes, txBytesRemote)
}
