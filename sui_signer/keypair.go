package sui_signer

import (
	"crypto/ed25519"
	"math"
)

type KeyPair byte

const (
	FlagEd25519 KeyPair = iota
	FlagSecp256k1

	FlagError = math.MaxUint8
)

func (k KeyPair) Byte() byte {
	return byte(k)
}

const (
	PublicKeyLengthEd25519   = 32
	PublicKeyLengthSecp256k1 = 33
)

const (
	DefaultAccountAddressLength = 16
	AccountAddress20Length      = 20
	AccountAddress32Length      = 32
)

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
