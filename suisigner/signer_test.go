package suisigner_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/suisigner/multisig"
	"github.com/pattonkan/sui-go/suisigner/suicrypto"
	"github.com/pattonkan/sui-go/utils/randtypes"

	"github.com/stretchr/testify/require"
)

func TestNewSigner(t *testing.T) {
	signer, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suicrypto.KeySchemeFlagEd25519)
	require.NoError(t, err)
	require.Equal(t, signer.Address, suisigner.TEST_ADDRESS)
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	signer, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suicrypto.KeySchemeFlagDefault)
	require.NoError(t, err)

	msg := "I want to have some bubble tea"
	msgBytes := []byte(msg)

	signature1, err := signer.SignDigest(msgBytes, suisigner.IntentTransaction())
	require.NoError(t, err)

	marshaledData, err := json.Marshal(signature1)
	require.NoError(t, err)

	var signature2 suisigner.Signature
	err = json.Unmarshal(marshaledData, &signature2)
	require.NoError(t, err)

	require.Equal(t, *signature1, signature2)
}

func TestSignSecp256k1Static(t *testing.T) {
	seed, err := hex.DecodeString("ac0c60b4cf8f6f975139f2594a059633386ff596b61e95b21dffbc1b30f197c1")
	require.NoError(t, err)
	targetSig, err := hex.DecodeString("0151405868ee0fcfb152873a8cc0b494aa8cb53fb8ba4967eefcfef0dbaff846d82a5e3b8f0cce5a7dfd99a0ecaaef3591eccb3e02f3cfaa3f01f484d0cf7e8d5c0295caed8aec9482012dffa084662f3e250f8d4bd4e63a13c1655d97d4a8e1182c")
	require.NoError(t, err)
	data := []byte("hello")

	keypair := suicrypto.NewKeypairSecp256k1FromSeed(seed)
	require.NotNil(t, keypair)
	signer := suisigner.NewSigner(seed, suicrypto.KeySchemeFlagSecp256k1)

	data, err = bcs.Marshal(data)
	require.NoError(t, err)

	sig, err := signer.SignDigest(data, suisigner.IntentPersonalMessage())
	require.NoError(t, err)
	require.Equal(t, targetSig, sig.Secp256k1SuiSignature.Signature[:])
}

func TestSignSecp256r1Static(t *testing.T) {
	seed, err := hex.DecodeString("a72c740033ddfeb3d22c5dc47d9b555e1d290d3c3b844554b0d94d311971767d")
	require.NoError(t, err)
	targetSig, err := hex.DecodeString("025a963dd9081366d9014aae2afd93761f62f8887ead54905c7fd8b20bf092678e251d182e32ff91a762be83ed89092f641435648f1e0d9c5ef1268b9e2c284fcf0387efef02e19fd54ae8b854dd79fa411169f76a9cbcc7711ec4de5ce444b31837")
	require.NoError(t, err)
	data := []byte("hello")

	keypair := suicrypto.NewKeypairSecp256r1FromSeed(seed)
	require.NotNil(t, keypair)
	signer := suisigner.NewSigner(seed, suicrypto.KeySchemeFlagSecp256r1)

	data, err = bcs.Marshal(data)
	require.NoError(t, err)

	sig, err := signer.SignDigest(data, suisigner.IntentPersonalMessage())
	require.NoError(t, err)
	require.Equal(t, targetSig, sig.Secp256r1SuiSignature.Signature[:])
}

func TestSign(t *testing.T) {
	tests := []struct {
		name string
		flag suicrypto.KeySchemeFlag
	}{
		{
			name: "successful, ed25519",
			flag: suicrypto.KeySchemeFlagEd25519,
		},
		{
			name: "successful, secp256k1",
			flag: suicrypto.KeySchemeFlagSecp256k1,
		},
		{
			name: "successful, secp256r1",
			flag: suicrypto.KeySchemeFlagSecp256r1,
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

				signature, err := signer.SignDigest(txBytes, suisigner.IntentTransaction())
				require.NoError(t, err)
				resp, err := c.ExecuteTransactionBlock(
					context.TODO(),
					&suiclient.ExecuteTransactionBlockRequest{
						TxDataBytes: txBytes,
						Signatures:  []*suisigner.Signature{signature},
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
	signer1, _ := suisigner.NewSignerWithMnemonic(mnemonic, suicrypto.KeySchemeFlagDefault)
	fmt.Printf("address   : %v\n", signer1.Address)

	// Create suisigner.Signer with private key
	privKey, _ := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
	signer2 := suisigner.NewSigner(privKey, suicrypto.KeySchemeFlagDefault)

	// Get private key, public key, address
	fmt.Printf("privateKey: %x\n", signer2.PrivateKeyBytes())
	fmt.Printf("publicKey : %x\n", signer2.PublicKeyBytes())
	fmt.Printf("address   : %v\n", signer2.Address)

	// Output:
	// address   : 0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5
	// privateKey: 4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b9342fa65507f5cf61f1b8fb3b94a5aa80fa9b2e2c68963e30d68a2660a50c57e
	// publicKey : 9342fa65507f5cf61f1b8fb3b94a5aa80fa9b2e2c68963e30d68a2660a50c57e
	// address   : 0x579a9ef1ca86431df106abb86f1f129806cd336b28f5bc17d16ce247aa3a0623
}

func TestMultisig(t *testing.T) {
	seed1 := [32]byte{0x01}
	keypair1 := suicrypto.NewKeypairSecp256k1FromSeed(seed1[:])
	require.NotNil(t, keypair1)
	seed2 := [32]byte{0x02}
	keypair2 := suicrypto.NewKeypairSecp256r1FromSeed(seed2[:])
	require.NotNil(t, keypair2)
	seed3 := [32]byte{0x03}
	keypair3 := suicrypto.NewKeypairEd25519FromSeed(seed3[:])
	require.NotNil(t, keypair3)

	committee := multisig.NewCommittee(
		[]*multisig.Member{
			{PublicKey: multisig.MemberPublicKey{Secp256k1PublicKey: keypair1.PubKey}, Weight: 2},
			{PublicKey: multisig.MemberPublicKey{Secp256r1PublicKey: keypair2.PubKey}, Weight: 2},
			{PublicKey: multisig.MemberPublicKey{Ed25519PublicKey: keypair3.PubKey}, Weight: 3},
		},
		3,
	)

	t.Run("pass threshold", func(t *testing.T) {
		tx := randtypes.RandomTransactionData()
		aggregator, err := multisig.NewAggregatorWithTransaction(committee, tx)
		require.NoError(t, err)
		digest, err := tx.SigningDigest()
		require.NoError(t, err)
		sig1, err := keypair1.Sign(digest)
		require.NoError(t, err)
		err = aggregator.AddSignature(multisig.MemberSignatureFromBytesSecp256k1(keypair1.PubKey, sig1))
		require.NoError(t, err)
		sig2, err := keypair2.Sign(digest)
		require.NoError(t, err)
		err = aggregator.AddSignature(multisig.MemberSignatureFromBytesSecp256r1(keypair2.PubKey, sig2))
		require.NoError(t, err)
		sig3, err := keypair3.Sign(digest)
		require.NoError(t, err)
		err = aggregator.AddSignature(multisig.MemberSignatureFromBytesEd25519(keypair3.PubKey, sig3))
		require.NoError(t, err)

		aggregatedSig1, err := aggregator.Finish()
		require.NoError(t, err)
		aggregatedSigBytes, err := bcs.Marshal(aggregatedSig1)
		require.NoError(t, err)

		verifier := multisig.Verifier{}
		aggregatedSig2 := multisig.NewAggregatedSignature()
		_, err = bcs.Unmarshal(aggregatedSigBytes, aggregatedSig2)
		require.NoError(t, err)
		err = verifier.VerifyAggregatedSignature(digest, aggregatedSig2)
		require.NoError(t, err)
	})

	t.Run("insufficient weight", func(t *testing.T) {
		tx := randtypes.RandomTransactionData()
		aggregator, err := multisig.NewAggregatorWithTransaction(committee, tx)
		require.NoError(t, err)
		digest, err := tx.SigningDigest()
		require.NoError(t, err)
		sig1, err := keypair1.Sign(digest)
		require.NoError(t, err)
		err = aggregator.AddSignature(multisig.MemberSignatureFromBytesSecp256k1(keypair1.PubKey, sig1))
		require.NoError(t, err)

		_, err = aggregator.Finish()
		require.Equal(t, fmt.Errorf("insufficient signature weight to reach threshold"), err)
	})
}
