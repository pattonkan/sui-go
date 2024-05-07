package isc_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/isc"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

type Client struct {
	API sui.ImplSuiAPI
}

func TestStartNewChain(t *testing.T) {
	api := sui.NewSuiClient(conn.LocalnetEndpointUrl)
	client := isc.NewClient(api)

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Log("sui_signer: ", signer.Address)
	digest, err := sui.RequestFundFromFaucet(signer.Address.String(), conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)

	packageID, anchorCap := isc.GetIscPackageIDAndAnchor(isc.GetGitRoot() + "/isc/isc/publish_receipt.json")

	res, err := client.StartNewChain(context.Background(), signer, packageID, anchorCap)
	require.NoError(t, err)
	t.Logf("StartNewChain response: %#v\n", res)
	for _, change := range res.ObjectChanges {
		fmt.Println("change.Data.Created: ", change.Data.Created)
	}
}

func TestRegisterIscToken(t *testing.T) {
	api := sui.NewSuiClient(conn.LocalnetEndpointUrl)
	client := isc.NewClient(api)
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)
	t.Log("sui_signer: ", signer.Address)

	packageID, anchorCap := isc.GetIscPackageIDAndAnchor(isc.GetGitRoot() + "/isc/isc/publish_receipt.json")

	// sender := sui_signer.TEST_ADDRESS
	gasBudget := sui_types.SUI(0.1).Uint64()
	gasPrice := uint64(1000)

	_, err = sui.RequestFundFromFaucet(signer.Address.String(), conn.TestnetFaucetUrl)
	require.NoError(t, err)
	coins, err := client.API.GetCoins(context.Background(), signer.Address, nil, nil, 2)
	require.NoError(t, err)
	pcoin, err := models.PickupCoins(coins, *big.NewInt(1), 100, 10, 3)
	require.NoError(t, err)

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()

	argAnchor, err := ptb.Pure(anchorCap.String())
	require.NoError(t, err)
	arg1 := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       packageID,
				Module:        "isc",
				Function:      "start_new_chain",
				TypeArguments: []sui_types.TypeTag{},
				Arguments:     []sui_types.Argument{argAnchor},
			},
		},
	)

	arg2, err := ptb.Pure(anchorCap.String())
	if err != nil {
		panic(err)
	}
	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:  packageID,
				Module:   "iscanchor",
				Function: "start_new_chain",
				Arguments: []sui_types.Argument{
					arg1, arg2,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*signer.Address,
		[]*sui_types.ObjectRef{
			pcoin.CoinRefs()[0],
			pcoin.CoinRefs()[1],
		},
		pt, gasBudget, gasPrice,
	)

	txnBytes, err := bcs.Marshal(tx)
	require.NoError(t, err)

	signature, err := signer.SignTransactionBlock(txnBytes, sui_signer.DefaultIntent())
	require.NoError(t, err)
	txnResponse, err := api.ExecuteTransactionBlock(context.TODO(), txnBytes, []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	require.NoError(t, err)
	t.Log(txnResponse)
}

// func TestSeriesCall(t *testing.T) {
// 	sender := sui_signer.TEST_ADDRESS
// 	gasBudget := sui_types.SUI(0.1).Uint64()
// 	gasPrice := uint64(1000)

// 	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
// 	_, err := sui.RequestFundFromFaucet(sender.String(), conn.TestnetFaucetUrl)
// 	require.NoError(t, err)
// 	coins := getCoins(t, api, sender, 2)
// 	coin, coin2 := coins[0], coins[1]

// 	validatorAddress, err := sui_types.SuiAddressFromHex(ComingChatValidatorAddress)
// 	require.NoError(t, err)

// 	// build with BCS
// 	ptb := sui_types.NewProgrammableTransactionBuilder()

// 	// case 1: split target amount
// 	amtArg, err := ptb.Pure(sui_types.SUI(1).Uint64())
// 	require.NoError(t, err)
// 	arg1 := ptb.Command(
// 		sui_types.Command{
// 			SplitCoins: &struct {
// 				Argument  sui_types.Argument
// 				Arguments []sui_types.Argument
// 			}{
// 				Argument:  sui_types.Argument{GasCoin: &lib.EmptyEnum{}},
// 				Arguments: []sui_types.Argument{amtArg},
// 			},
// 		},
// 	) // the coin is split result argument
// 	arg2, err := ptb.Pure(validatorAddress)
// 	require.NoError(t, err)
// 	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
// 	require.NoError(t, err)
// 	ptb.Command(
// 		sui_types.Command{
// 			MoveCall: &sui_types.ProgrammableMoveCall{
// 				Package:  sui_types.SuiSystemAddress,
// 				Module:   sui_system_state.SuiSystemModuleName,
// 				Function: sui_types.AddStakeFunName,
// 				Arguments: []sui_types.Argument{
// 					arg0, arg1, arg2,
// 				},
// 			},
// 		},
// 	)
// 	pt := ptb.Finish()
// 	tx := sui_types.NewProgrammable(
// 		*sender, []*sui_types.ObjectRef{
// 			coin.Reference(),
// 			coin2.Reference(),
// 		},
// 		pt, gasBudget, gasPrice,
// 	)

// 	// build & simulate
// 	txBytesBCS, err := bcs.Marshal(tx)
// 	require.NoError(t, err)
// 	resp := dryRunTxn(t, api, txBytesBCS, true)
// 	t.Log(resp.Effects.Data.GasFee())
// }
