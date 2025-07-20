package suicrypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

type KeypairEd25519 struct {
	PriKey *Ed25519PriKey
	PubKey *Ed25519PubKey
}

type Ed25519PriKey ed25519.PrivateKey
type Ed25519PubKey ed25519.PublicKey

func Ed25519PubKeyFromBytes(b []byte) (*Ed25519PubKey, error) {
	if len(b) != KeypairEd25519PublicKeySize {
		return nil, fmt.Errorf("invalid public key size")
	}
	pubkey := Ed25519PubKey(b)
	return &pubkey, nil
}
func (p Ed25519PubKey) Bytes() []byte {
	return p
}
func (p Ed25519PubKey) String() string {
	return hex.EncodeToString(p)
}
func (p Ed25519PriKey) Bytes() []byte {
	return p
}

func NewKeypairEd25519FromSeed(seed []byte) *KeypairEd25519 {
	if len(seed) < ed25519.SeedSize {
		return nil
	}
	primitivePrivKey := ed25519.NewKeyFromSeed(seed[:ed25519.SeedSize])
	prikey := Ed25519PriKey(primitivePrivKey)
	pubkey := Ed25519PubKey(primitivePrivKey.Public().(ed25519.PublicKey))
	return &KeypairEd25519{
		PriKey: &prikey,
		PubKey: &pubkey,
	}
}

func (k *KeypairEd25519) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(k.PriKey.Bytes(), data), nil
}

func (p Ed25519PubKey) Verify(data []byte, sig []byte) bool {
	return ed25519.Verify(p.Bytes(), data, sig)
}
