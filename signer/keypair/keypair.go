package keypair

import (
	"math"
)

type KeyPair byte

const (
	Ed25519Flag   KeyPair = 0
	Secp256k1Flag KeyPair = 1
	ErrorFlag     byte    = math.MaxUint8
)

const (
	ed25519PublicKeyLength   = 32
	secp256k1PublicKeyLength = 33
)

const (
	DefaultAccountAddressLength = 16
	AccountAddress20Length      = 20
	AccountAddress32Length      = 32
)
