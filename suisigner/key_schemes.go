package suisigner

import (
	"crypto/ed25519"
	"math"
)

type KeySchemeFlag byte

var KeySchemeFlagDefault = KeySchemeFlagEd25519

const (
	KeySchemeFlagEd25519 KeySchemeFlag = iota
	KeySchemeFlagSecp256k1
	KeySchemeFlagSecp256r1
	KeySchemeFlagMultiSig
	KeySchemeFlagBLS12381
	KeySchemeFlagZkLoginAuthenticator

	KeySchemeFlagIotaEd25519 = math.MaxUint8 - 1 // special case for iota ed25519
	KeySchemeFlagError       = math.MaxUint8
)

func (k KeySchemeFlag) Byte() byte {
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
