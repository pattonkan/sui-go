package suicrypto

import (
	"crypto/ed25519"
)

type KeypairEd25519 struct {
	PriKey Ed25519PriKey
	PubKey Ed25519PubKey
}

type Ed25519PriKey ed25519.PrivateKey
type Ed25519PubKey ed25519.PublicKey

func (p Ed25519PubKey) Bytes() []byte {
	return p
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
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func (k *KeypairEd25519) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(k.PriKey.Bytes(), data), nil
}
