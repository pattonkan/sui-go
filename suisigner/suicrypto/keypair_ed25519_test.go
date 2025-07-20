package suicrypto_test

import (
	"testing"

	"github.com/pattonkan/sui-go/suisigner/suicrypto"

	"github.com/stretchr/testify/require"
)

func TestKeypairEd25519Bytes(t *testing.T) {
	seed3 := [32]byte{0x03}
	keypair := suicrypto.NewKeypairEd25519FromSeed(seed3[:])
	require.NotNil(t, keypair)

	pubkey, err := suicrypto.Ed25519PubKeyFromBytes(keypair.PubKey.Bytes())
	require.NoError(t, err)
	require.Equal(t, keypair.PubKey, pubkey)
}
