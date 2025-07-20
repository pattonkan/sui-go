package suicrypto_test

import (
	"testing"

	"github.com/pattonkan/sui-go/suisigner/suicrypto"

	"github.com/stretchr/testify/require"
)

func TestKeypairSecp256r1Bytes(t *testing.T) {
	seed3 := [32]byte{0x03}
	keypair := suicrypto.NewKeypairSecp256r1FromSeed(seed3[:])
	require.NotNil(t, keypair)

	pubkey, err := suicrypto.Secp256r1PubKeyFromBytes(keypair.PubKey.Bytes())
	require.NoError(t, err)
	require.Equal(t, keypair.PubKey, pubkey)
}
