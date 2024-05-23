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
	"github.com/stretchr/testify/require"
)

func TestPTB_Pay(t *testing.T) {
	client, signer := sui.NewTestSuiClientWithSignerAndFund(conn.DevnetEndpointUrl, sui_signer.TEST_MNEMONIC)
	recipients := sui_signer.NewRandomSigners(sui_signer.KeySchemeFlagDefault, 2)
	coinType := models.SuiCoinType
	limit := uint(3)
	coinPages, err := client.GetCoins(context.Background(), signer.Address, &coinType, nil, limit)
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
		signer.Address,
		pt,
		[]*sui_types.ObjectRef{
			gasCoin.Reference(),
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
			require.Contains(t, []*sui_types.ObjectID{gasCoin.CoinObjectID, transferCoins[0].CoinObjectID}, &change.Data.Mutated.ObjectID)
		} else if change.Data.Deleted != nil {
			require.Equal(t, transferCoins[1].CoinObjectID, &change.Data.Deleted.ObjectID)
		}
	}
	require.Len(t, simulate.BalanceChanges, 3)
	for _, balChange := range simulate.BalanceChanges {
		if balChange.Owner.AddressOwner == signer.Address {
			require.Equal(t, totalBal-(amounts[0]+amounts[1]), balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipients[0].Address {
			require.Equal(t, amounts[0], balChange.Amount)
		} else if balChange.Owner.AddressOwner == recipients[1].Address {
			require.Equal(t, amounts[1], balChange.Amount)
		}
	}
}

func TestTransferSui(t *testing.T) {
	recipient, err := sui_types.SuiAddressFromHex("0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e")
	require.NoError(t, err)

	ptb := sui_types.NewProgrammableTransactionBuilder()
	amount := uint64(100000)
	err = ptb.TransferSui(recipient, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	digest := sui_types.MustNewDigest("HvbE2UZny6cP4KukaXetmj4jjpKTDTjVo23XEcu7VgSn")
	objectId, err := sui_types.SuiAddressFromHex("0x13c1c3d0e15b4039cec4291c75b77c972c10c8e8e70ab4ca174cf336917cb4db")
	require.NoError(t, err)
	tx := sui_types.NewProgrammable(
		recipient,
		pt,
		[]*sui_types.ObjectRef{
			{
				ObjectID: objectId,
				Version:  14924029,
				Digest:   digest,
			},
		},
		10000000,
		1000,
	)
	txByte, err := bcs.Marshal(tx)
	require.NoError(t, err)
	t.Logf("%x", txByte)
}
