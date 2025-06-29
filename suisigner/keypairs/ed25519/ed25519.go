package ed25519

import "crypto/ed25519"

type KeypairEd25519 struct {
	PriKey ed25519.PrivateKey
	PubKey ed25519.PublicKey
}

func NewKeypairEd25519(prikey ed25519.PrivateKey, pubkey ed25519.PublicKey) *KeypairEd25519 {
	return &KeypairEd25519{
		PriKey: prikey,
		PubKey: pubkey,
	}
}
