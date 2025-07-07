package suisigner

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
)

type Signature struct {
	*Ed25519SuiSignature
	*Secp256k1SuiSignature
	*Secp256r1SuiSignature
}

const (
	SizeEd25519SuiSignature   = ed25519.PublicKeySize + ed25519.SignatureSize + 1
	SizeSecp256k1SuiSignature = KeypairSecp256k1PublicKeySize + KeypairSecp256k1SignatureSize + 1
	SizeSecp256r1SuiSignature = KeypairSecp256r1PublicKeySize + KeypairSecp256r1SignatureSize + 1
)

type Ed25519SuiSignature struct {
	Signature [SizeEd25519SuiSignature]byte
}

type Secp256k1SuiSignature struct {
	Signature [SizeSecp256k1SuiSignature]byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Secp256r1SuiSignature struct {
	Signature [SizeSecp256r1SuiSignature]byte //secp256r1.pubKey + Secp256r1Signature + 1
}

func (s Signature) Bytes() []byte {
	switch {
	case s.Ed25519SuiSignature != nil:
		return s.Ed25519SuiSignature.Signature[:]
	case s.Secp256k1SuiSignature != nil:
		return s.Secp256k1SuiSignature.Signature[:]
	case s.Secp256r1SuiSignature != nil:
		return s.Secp256r1SuiSignature.Signature[:]
	default:
		return nil
	}
}

func (s Signature) MarshalJSON() ([]byte, error) {
	switch {
	case s.Ed25519SuiSignature != nil:
		return json.Marshal(s.Ed25519SuiSignature.Signature[:])
	case s.Secp256k1SuiSignature != nil:
		return json.Marshal(s.Secp256k1SuiSignature.Signature[:])
	case s.Secp256r1SuiSignature != nil:
		return json.Marshal(s.Secp256r1SuiSignature.Signature[:])
	default:
		return nil, errors.New("nil signature")
	}
}

func (s *Signature) UnmarshalJSON(data []byte) error {
	var signature []byte
	err := json.Unmarshal(data, &signature)
	if err != nil {
		return err
	}
	switch signature[0] {
	case KeySchemeFlagEd25519.Byte():
		if len(signature) != ed25519.PublicKeySize+ed25519.SignatureSize+1 {
			return errors.New("invalid ed25519 signature")
		}
		var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
		copy(signatureBytes[:], signature)
		s.Ed25519SuiSignature = &Ed25519SuiSignature{
			Signature: signatureBytes,
		}
	case KeySchemeFlagSecp256k1.Byte():
		if len(signature) != KeypairSecp256k1PublicKeySize+KeypairSecp256k1SignatureSize+1 {
			return errors.New("invalid secp256k1 signature")
		}
		var signatureBytes [KeypairSecp256k1PublicKeySize + KeypairSecp256k1SignatureSize + 1]byte
		copy(signatureBytes[:], signature)
		s.Secp256k1SuiSignature = &Secp256k1SuiSignature{
			Signature: signatureBytes,
		}
	default:
		return errors.New("not supported signature")
	}
	return nil
}
