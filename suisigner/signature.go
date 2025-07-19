package suisigner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pattonkan/sui-go/suisigner/suicrypto"
)

type Signature struct {
	*Ed25519SuiSignature
	*Secp256k1SuiSignature
	*Secp256r1SuiSignature
}

type Ed25519SuiSignature struct {
	Signature [suicrypto.SizeSuiSignatureEd25519]byte
}

type Secp256k1SuiSignature struct {
	Signature [suicrypto.SizeSuiSignatureSecp256k1]byte
}

type Secp256r1SuiSignature struct {
	Signature [suicrypto.SizeSuiSignatureSecp256r1]byte
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
	case suicrypto.KeySchemeFlagEd25519.Byte():
		if len(signature) != suicrypto.SizeSuiSignatureEd25519 {
			return errors.New("invalid ed25519 signature")
		}
		var signatureBytes [suicrypto.SizeSuiSignatureEd25519]byte
		copy(signatureBytes[:], signature)
		s.Ed25519SuiSignature = &Ed25519SuiSignature{
			Signature: signatureBytes,
		}
	case suicrypto.KeySchemeFlagSecp256k1.Byte():
		if len(signature) != suicrypto.SizeSuiSignatureSecp256k1 {
			return errors.New("invalid secp256k1 signature")
		}
		var signatureBytes [suicrypto.SizeSuiSignatureSecp256k1]byte
		copy(signatureBytes[:], signature)
		s.Secp256k1SuiSignature = &Secp256k1SuiSignature{
			Signature: signatureBytes,
		}
	case suicrypto.KeySchemeFlagSecp256r1.Byte():
		if len(signature) != suicrypto.SizeSuiSignatureSecp256r1 {
			return errors.New("invalid secp256r1 signature")
		}
		var signatureBytes [suicrypto.SizeSuiSignatureSecp256r1]byte
		copy(signatureBytes[:], signature)
		s.Secp256r1SuiSignature = &Secp256r1SuiSignature{
			Signature: signatureBytes,
		}
	default:
		return errors.New("not supported signature")
	}
	return nil
}

func NewEd25519SuiSignature(s *suicrypto.KeypairEd25519, data []byte) (*Ed25519SuiSignature, error) {
	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(suicrypto.KeySchemeFlagEd25519))

	sig, err := s.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data with ed25519: %w", err)
	}
	sigBuffer.Write(sig[:])
	sigBuffer.Write(s.PubKey.Bytes())

	return &Ed25519SuiSignature{
		Signature: [suicrypto.SizeSuiSignatureEd25519]byte(sigBuffer.Bytes()),
	}, nil
}
func NewSecp256k1SuiSignature(s *suicrypto.KeypairSecp256k1, data []byte) (*Secp256k1SuiSignature, error) {
	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(suicrypto.KeySchemeFlagSecp256k1))

	sig, err := s.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data with secp256k1: %w", err)
	}
	sigBuffer.Write(sig[:])
	sigBuffer.Write(s.PubKey.Bytes())

	return &Secp256k1SuiSignature{
		Signature: [suicrypto.SizeSuiSignatureSecp256k1]byte(sigBuffer.Bytes()),
	}, nil
}
func NewSecp256r1SuiSignature(s *suicrypto.KeypairSecp256r1, data []byte) (*Secp256r1SuiSignature, error) {
	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(suicrypto.KeySchemeFlagSecp256r1))

	sig, err := s.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data with secp256r1: %w", err)
	}
	sigBuffer.Write(sig[:])
	sigBuffer.Write(s.PubKey.Bytes())

	return &Secp256r1SuiSignature{
		Signature: [suicrypto.SizeSuiSignatureSecp256r1]byte(sigBuffer.Bytes()),
	}, nil
}
