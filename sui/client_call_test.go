package sui_test

import (
	"encoding/json"
	"testing"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"

	"github.com/stretchr/testify/require"
)

func Test_TagJson_Owner(t *testing.T) {
	test := func(str string) serialization.TagJson[sui_types.Owner] {
		var s serialization.TagJson[sui_types.Owner]
		data := []byte(str)
		err := json.Unmarshal(data, &s)
		require.NoError(t, err)
		return s
	}
	{
		v := test(`"Immutable"`).Data
		require.Nil(t, v.AddressOwner)
		require.Nil(t, v.ObjectOwner)
		require.Nil(t, v.Shared)
		require.NotNil(t, v.Immutable)
	}
	{
		v := test(`{"AddressOwner": "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e"}`).Data
		require.NotNil(t, v.AddressOwner)
		require.Nil(t, v.ObjectOwner)
		require.Nil(t, v.Shared)
		require.Nil(t, v.Immutable)
	}
}

// TestExecuteTransactionSerializedSig
// This test case will affect the real coin in the test case of account
// temporary disabled
//func TestExecuteTransactionSerializedSig(t *testing.T) {
//	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
//	coins, err := api.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := api.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectID, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := api.ExecuteTransactionSerializedSig(context.TODO(), *signedTx, models.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}

//func TestExecuteTransaction(t *testing.T) {
//	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
//	coins, err := api.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := api.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectID, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := api.ExecuteTransaction(context.TODO(), *signedTx, models.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}
