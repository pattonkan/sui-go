package schemes

import "math"

type KeySchemeFlag byte

var KeySchemeFlagDefault = KeySchemeFlagEd25519

const (
	KeySchemeFlagEd25519 KeySchemeFlag = iota
	KeySchemeFlagSecp256k1
	KeySchemeFlagSecp256r1
	KeySchemeFlagMultiSig
	KeySchemeFlagBLS12381
	KeySchemeFlagZkLogin
	KeySchemeFlagPasskey

	KeySchemeFlagError = math.MaxUint8
)

func (k KeySchemeFlag) Byte() byte {
	return byte(k)
}
