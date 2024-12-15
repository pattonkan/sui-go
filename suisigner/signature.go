package suisigner

import (
	"bytes"
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
	SizeEd25519SuiSignature = ed25519.PublicKeySize + ed25519.SignatureSize + 1
)

type Secp256k1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Secp256r1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Ed25519SuiSignature struct {
	Signature [SizeEd25519SuiSignature]byte
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
	case 0:
		if len(signature) != ed25519.PublicKeySize+ed25519.SignatureSize+1 {
			return errors.New("invalid ed25519 signature")
		}
		var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
		copy(signatureBytes[:], signature)
		s.Ed25519SuiSignature = &Ed25519SuiSignature{
			Signature: signatureBytes,
		}
	default:
		return errors.New("not supported signature")
	}
	return nil
}

func NewEd25519SuiSignature(s *Signer, msg []byte) *Ed25519SuiSignature {
	sig := ed25519.Sign(s.ed25519Keypair.PriKey, msg)

	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(KeySchemeFlagEd25519))
	sigBuffer.Write(sig[:])
	sigBuffer.Write(s.ed25519Keypair.PubKey)

	return &Ed25519SuiSignature{
		Signature: [SizeEd25519SuiSignature]byte(sigBuffer.Bytes()),
	}
}
