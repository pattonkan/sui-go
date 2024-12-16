package suisigner_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/stretchr/testify/require"
)

func TestNewSigner(t *testing.T) {
	testMnemonic := "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
	testEd25519Address := sui.MustAddressFromHex("0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5")
	signer, err := suisigner.NewSignerWithMnemonic(testMnemonic, suisigner.KeySchemeFlagEd25519)
	require.NoError(t, err)
	require.Equal(t, testEd25519Address, signer.Address)
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	signer, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suisigner.KeySchemeFlagDefault)
	require.NoError(t, err)

	msg := "I want to have some bubble tea"
	msgBytes := []byte(msg)

	signature1, err := signer.SignTransactionBlock(msgBytes, suisigner.DefaultIntent())
	require.NoError(t, err)

	marshaledData, err := json.Marshal(signature1)
	require.NoError(t, err)

	var signature2 suisigner.Signature
	err = json.Unmarshal(marshaledData, &signature2)
	require.NoError(t, err)

	require.Equal(t, signature1, signature2)
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
