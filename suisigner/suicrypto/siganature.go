package suicrypto

import (
	"errors"
)

func ExtractPubkeyAndSignature(rawSig []byte) (KeySchemeFlag, []byte, []byte, error) {
	schemeFlag := rawSig[0]
	var pubkey, signature []byte
	switch schemeFlag {
	case KeySchemeFlagEd25519.Byte():
		if len(rawSig) != SizeSuiSignatureEd25519 {
			return KeySchemeFlagError, nil, nil, errors.New("invalid ed25519 signature")
		}
		pubkey = make([]byte, KeypairEd25519PublicKeySize)
		signature = make([]byte, KeypairEd25519SignatureSize)
		copy(pubkey, rawSig[1:KeypairEd25519PublicKeySize+1])
		copy(signature, rawSig[1+KeypairEd25519PublicKeySize:])
		return KeySchemeFlagEd25519, pubkey, signature, nil

	case KeySchemeFlagSecp256k1.Byte():
		if len(rawSig) != SizeSuiSignatureSecp256k1 {
			return KeySchemeFlagError, nil, nil, errors.New("invalid secp256k1 signature")
		}
		pubkey = make([]byte, KeypairSecp256k1PublicKeySize)
		signature = make([]byte, KeypairSecp256k1SignatureSize)
		copy(pubkey, rawSig[1:KeypairSecp256k1PublicKeySize+1])
		copy(signature, rawSig[1+KeypairSecp256k1PublicKeySize:])
		return KeySchemeFlagSecp256k1, pubkey, signature, nil

	case KeySchemeFlagSecp256r1.Byte():
		if len(rawSig) != SizeSuiSignatureSecp256r1 {
			return KeySchemeFlagError, nil, nil, errors.New("invalid secp256r1 signature")
		}
		pubkey = make([]byte, KeypairSecp256r1PublicKeySize)
		signature = make([]byte, KeypairSecp256r1SignatureSize)
		copy(pubkey, rawSig[1:KeypairSecp256r1PublicKeySize+1])
		copy(signature, rawSig[1+KeypairSecp256r1PublicKeySize:])
		return KeySchemeFlagSecp256r1, pubkey, signature, nil

	default:
		return KeySchemeFlagError, nil, nil, errors.New("not supported signature")
	}

}
