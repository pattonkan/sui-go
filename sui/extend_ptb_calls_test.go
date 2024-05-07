package sui_test

import (
	"context"
	"fmt"
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

func TestPTB_PaySui(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	recipient, _ := sui_types.SuiAddressFromHex("0x123456")
	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coin := getCoins(t, api, sender, 1)[0]

	gasPrice := uint64(1000)
	// gasPrice, err := api.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PaySui([]*sui_types.SuiAddress{recipient, recipient}, []uint64{amount, amount})
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	resp := dryRunTxn(t, api, txBytesBCS, true)
	gasFee := resp.Effects.Data.GasFee()
	t.Log(gasFee)

	// build with remote rpc
	// txn, err := api.PaySui(context.Background(), *sender, []sui_types.ObjectID{coin.CoinObjectID},
	// 	[]*sui_types.SuiAddress{recipient, recipient},
	// 	[]models.SafeSuiBigInt[uint64]{models.NewSafeSuiBigInt(amount), models.NewSafeSuiBigInt(amount)},
	// 	models.NewSafeSuiBigInt(gasBudget))
	// require.NoError(t, err)
	// txBytesRemote := txn.TxBytes.Data()

	// XXX: Fail when there are multiple recipients
	// require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTB_TransferObject(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	recipient := sui_signer.TEST_ADDRESS
	gasBudget := sui_types.SUI(0.1).Uint64()

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coins := getCoins(t, api, sender, 2)
	coin, gas := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := api.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferObject(recipient, []*sui_types.ObjectRef{coin.Reference()})
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			gas.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := api.TransferObject(
		context.Background(), sender, recipient,
		&coin.CoinObjectID,
		&gas.CoinObjectID,
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTB_TransferSui(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	recipient := sender
	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coin := getCoins(t, api, sender, 1)[0]

	gasPrice := uint64(1000)
	// gasPrice, err := api.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferSui(recipient, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := api.TransferSui(
		context.Background(), sender, recipient, &coin.CoinObjectID,
		models.NewSafeSuiBigInt(amount),
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTB_PayAllSui(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	recipient := sender
	gasBudget := sui_types.SUI(0.01).Uint64()

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coins := getCoins(t, api, sender, 2)
	coin, coin2 := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := api.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PayAllSui(recipient)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
			coin2.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := api.PayAllSui(
		context.Background(), sender, recipient,
		[]sui_types.ObjectID{
			coin.CoinObjectID, coin2.CoinObjectID,
		},
		models.NewSafeSuiBigInt(gasBudget),
	)
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTB_Pay(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	recipient, _ := sui_types.SuiAddressFromHex("0x123456")
	amount := sui_types.SUI(0.001).Uint64()
	gasBudget := sui_types.SUI(0.01).Uint64()

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coins := getCoins(t, api, sender, 2)
	coin, gas := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := api.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.Pay(
		[]*sui_types.ObjectRef{coin.Reference()},
		[]*sui_types.SuiAddress{recipient, recipient},
		[]uint64{amount, amount},
	)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			gas.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	resp := dryRunTxn(t, api, txBytesBCS, true)
	gasfee := resp.Effects.Data.GasFee()
	t.Log(gasfee)

	// build with remote rpc
	// txn, err := api.Pay(context.Background(), *sender,
	// 	[]sui_types.ObjectID{coin.CoinObjectID},
	// 	[]*sui_types.SuiAddress{recipient, recipient},
	// 	[]models.SafeSuiBigInt[uint64]{models.NewSafeSuiBigInt(amount), models.NewSafeSuiBigInt(amount)},
	// 	&gas.CoinObjectID,
	// 	models.NewSafeSuiBigInt(gasBudget))
	// require.NoError(t, err)
	// txBytesRemote := txn.TxBytes.Data()

	// XXX: Fail when there are multiple recipients
	// require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestPTB_MoveCall(t *testing.T) {
	sender := sui_signer.TEST_ADDRESS
	gasBudget := sui_types.SUI(0.1).Uint64()
	gasPrice := uint64(1000)

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	_, err := sui.RequestFundFromFaucet(sender, conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coins := getCoins(t, api, sender, 2)
	coin, coin2 := coins[0], coins[1]

	validatorAddress, err := sui_types.SuiAddressFromHex(ComingChatValidatorAddress)
	require.NoError(t, err)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()

	// case 1: split target amount
	amtArg, err := ptb.Pure(sui_types.SUI(1).Uint64())
	require.NoError(t, err)
	arg1 := ptb.Command(
		sui_types.Command{
			SplitCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  sui_types.Argument{GasCoin: &serialization.EmptyEnum{}},
				Arguments: []sui_types.Argument{amtArg},
			},
		},
	) // the coin is split result argument
	arg2, err := ptb.Pure(validatorAddress)
	require.NoError(t, err)
	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
	require.NoError(t, err)
	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:  sui_types.SuiSystemAddress,
				Module:   sui_types.SuiSystemModuleName,
				Function: sui_types.AddStakeFunName,
				Arguments: []sui_types.Argument{
					arg0, arg1, arg2,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
			coin2.Reference(),
		},
		pt, gasBudget, gasPrice,
	)

	// case 2: direct stake the specified coin
	// coinArg := sui_types.CallArg{
	// 	Object: &sui_types.ObjectArg{
	// 		ImmOrOwnedObject: coin.Reference(),
	// 	},
	// }
	// addrBytes := validatorAddress.Data()
	// addrArg := sui_types.CallArg{
	// 	Pure: &addrBytes,
	// }
	// err = ptb.MoveCall(
	// 	*sui_types.SuiSystemAddress,
	// 	sui_system_state.SuiSystemModuleName,
	// 	sui_types.AddStakeFunName,
	// 	[]sui_types.TypeTag{},
	// 	[]sui_types.CallArg{
	// 		sui_types.SuiSystemMut,
	// 		coinArg,
	// 		addrArg,
	// 	},
	// )
	// require.NoError(t, err)
	// pt := ptb.Finish()
	// tx := sui_types.NewProgrammable(
	// 	*sender, []*sui_types.ObjectRef{
	// 		coin2.Reference(),
	// 	},
	// 	pt, gasBudget, gasPrice,
	// )

	// build & simulate
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)
	resp := dryRunTxn(t, api, txBytesBCS, true)
	t.Log(resp.Effects.Data.GasFee())
}

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)

	options := models.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
	}
	coinType := fmt.Sprintf("0x2::coin::Coin<%v>", models.SuiCoinType)
	filterObject, err := api.BatchGetObjectsOwnedByAddress(context.TODO(), sui_signer.TEST_ADDRESS, &options, coinType)
	require.NoError(t, err)
	t.Log(filterObject)
}

func getCoins(t *testing.T, api *sui.ImplSuiAPI, sender *sui_types.SuiAddress, needCoinObjNum int) []models.Coin {
	coins, err := api.GetCoins(context.Background(), sender, nil, nil, uint(needCoinObjNum))
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= needCoinObjNum)
	return coins.Data
}
