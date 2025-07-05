package suisigner_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/stretchr/testify/require"
)

func TestNewSigner(t *testing.T) {
	signer, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suisigner.KeySchemeFlagEd25519)
	require.NoError(t, err)
	require.Equal(t, signer.Address, suisigner.TEST_ADDRESS)
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	signer, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suisigner.KeySchemeFlagDefault)
	require.NoError(t, err)

	msg := "I want to have some bubble tea"
	msgBytes := []byte(msg)

	signature1, err := signer.SignDigest(msgBytes, suisigner.DefaultIntent())
	require.NoError(t, err)

	marshaledData, err := json.Marshal(signature1)
	require.NoError(t, err)

	var signature2 suisigner.Signature
	err = json.Unmarshal(marshaledData, &signature2)
	require.NoError(t, err)

	require.Equal(t, signature1, signature2)
}

func TestSignSecp256k1Static(t *testing.T) {
	seed, err := hex.DecodeString("ac0c60b4cf8f6f975139f2594a059633386ff596b61e95b21dffbc1b30f197c1")
	require.NoError(t, err)
	targetSig, err := hex.DecodeString("0151405868ee0fcfb152873a8cc0b494aa8cb53fb8ba4967eefcfef0dbaff846d82a5e3b8f0cce5a7dfd99a0ecaaef3591eccb3e02f3cfaa3f01f484d0cf7e8d5c0295caed8aec9482012dffa084662f3e250f8d4bd4e63a13c1655d97d4a8e1182c")
	require.NoError(t, err)
	data := []byte("hello")

	keypair := suisigner.NewKeypairSecp256k1FromSeed(seed)
	require.NotNil(t, keypair)
	signer := suisigner.NewSigner(seed, suisigner.KeySchemeFlagSecp256k1)

	intent := suisigner.Intent{
		Scope: suisigner.IntentScope{
			PersonalMessage: &sui.EmptyEnum{},
		},
		Version: suisigner.IntentVersion{
			V0: &sui.EmptyEnum{},
		},
		AppId: suisigner.AppId{
			Sui: &sui.EmptyEnum{},
		},
	}
	data, err = bcs.Marshal(data)
	require.NoError(t, err)

	sig, err := signer.SignDigest(data, intent)
	require.NoError(t, err)
	require.Equal(t, targetSig, sig.Secp256k1SuiSignature.Signature[:])
}

func TestSign(t *testing.T) {
	tests := []struct {
		name string
		flag suisigner.KeySchemeFlag
	}{
		{
			name: "successful, ed25519",
			flag: suisigner.KeySchemeFlagEd25519,
		},
		{
			name: "successful, secp256k1",
			flag: suisigner.KeySchemeFlagSecp256k1,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c, signer := suiclient.NewClient(conn.LocalnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, tt.flag, 0)

				coinPages, err := c.GetCoins(context.Background(), &suiclient.GetCoinsRequest{Owner: signer.Address})
				require.NoError(t, err)
				coins := suiclient.Coins(coinPages.Data)

				ptb := suiptb.NewTransactionDataTransactionBuilder()
				err = ptb.PayAllSui(signer.Address)
				require.NoError(t, err)
				pt := ptb.Finish()
				tx := suiptb.NewTransactionData(
					signer.Address,
					pt,
					coins.CoinRefs(),
					suiclient.DefaultGasBudget,
					suiclient.DefaultGasPrice,
				)
				txBytes, err := bcs.Marshal(tx)
				require.NoError(t, err)
				options := &suiclient.SuiTransactionBlockResponseOptions{ShowEffects: true}

				signature, err := signer.SignDigest(txBytes, suisigner.DefaultIntent())
				require.NoError(t, err)
				resp, err := c.ExecuteTransactionBlock(
					context.TODO(),
					&suiclient.ExecuteTransactionBlockRequest{
						TxDataBytes: txBytes,
						Signatures:  []*suisigner.Signature{&signature},
						Options:     options,
						RequestType: suiclient.TxnRequestTypeWaitForLocalExecution,
					},
				)
				require.NoError(t, err)
				require.True(t, resp.Effects.Data.IsSuccess())
			},
		)
	}
}

func ExampleSigner() {
	// Create a suisigner.Signer with mnemonic
	mnemonic := "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
	signer1, _ := suisigner.NewSignerWithMnemonic(mnemonic, suisigner.KeySchemeFlagDefault)
	fmt.Printf("address   : %v\n", signer1.Address)

	// Create suisigner.Signer with private key
	privKey, _ := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
	signer2 := suisigner.NewSigner(privKey, suisigner.KeySchemeFlagDefault)

	// Get private key, public key, address
	fmt.Printf("privateKey: %x\n", signer2.PrivateKey()[:32])
	fmt.Printf("publicKey : %x\n", signer2.PublicKey())
	fmt.Printf("address   : %v\n", signer2.Address)

	// Output:
	// address   : 0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5
	// privateKey: 4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b
	// publicKey : 9342fa65507f5cf61f1b8fb3b94a5aa80fa9b2e2c68963e30d68a2660a50c57e
	// address   : 0x579a9ef1ca86431df106abb86f1f129806cd336b28f5bc17d16ce247aa3a0623
}
