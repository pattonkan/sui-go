package suicrypto

import "math"

// KeySchemeFlag represents the type of key scheme used for signatures.
type KeySchemeFlag byte

var KeySchemeFlagDefault = KeySchemeFlagEd25519

const (
	KeySchemeFlagEd25519 KeySchemeFlag = iota
	KeySchemeFlagSecp256k1
	KeySchemeFlagSecp256r1
	KeySchemeFlagMultiSig
	KeySchemeFlagBLS12381
	KeySchemeFlagZkLoginAuthenticator
	KeySchemeFlagPasskeyAuthenticator

	KeySchemeFlagError = math.MaxUint8
)

func (k KeySchemeFlag) Byte() byte {
	return byte(k)
}

func (k KeySchemeFlag) String() string {
	switch k {
	case KeySchemeFlagEd25519:
		return "Ed25519"
	case KeySchemeFlagSecp256k1:
		return "Secp256k1"
	case KeySchemeFlagSecp256r1:
		return "Secp256r1"
	case KeySchemeFlagMultiSig:
		return "MultiSig"
	case KeySchemeFlagBLS12381:
		return "BLS12381"
	case KeySchemeFlagZkLoginAuthenticator:
		return "ZkLoginAuthenticator"
	case KeySchemeFlagPasskeyAuthenticator:
		return "PasskeyAuthenticator"
	default:
		return "Unknown"
	}
}
