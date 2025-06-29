package suisigner

import (
	"bytes"
	"crypto/ed25519"
)

type KeypairEd25519 struct {
	PriKey ed25519.PrivateKey
	PubKey ed25519.PublicKey
}

func NewKeypairEd25519FromSeed(seed []byte) *KeypairEd25519 {
	if len(seed) < ed25519.SeedSize {
		return nil
	}
	prikey := ed25519.NewKeyFromSeed(seed[:ed25519.SeedSize])
	pubkey := prikey.Public().(ed25519.PublicKey)
	return &KeypairEd25519{
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func NewKeypairEd25519(prikey ed25519.PrivateKey, pubkey ed25519.PublicKey) *KeypairEd25519 {
	return &KeypairEd25519{
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func NewEd25519SuiSignature(s *Signer, msg []byte) *Ed25519SuiSignature {
	sig := ed25519.Sign(s.KeypairEd25519.PriKey, msg)

	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(KeySchemeFlagEd25519))
	sigBuffer.Write(sig[:])
	sigBuffer.Write(s.KeypairEd25519.PubKey)

	return &Ed25519SuiSignature{
		Signature: [SizeEd25519SuiSignature]byte(sigBuffer.Bytes()),
	}
}
