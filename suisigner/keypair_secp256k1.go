package suisigner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp256k1_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

// refer sui-rust-sdk/crates/sui-sdk-types/src/crypto/signature.rs
const (
	KeypairSecp256k1PublicKeySize = 33
	KeypairSecp256k1SignatureSize = 64
)

type KeypairSecp256k1 struct {
	PriKey *secp256k1.PrivateKey
	PubKey *secp256k1.PublicKey
}

func NewKeypairSecp256k1FromSeed(seed []byte) *KeypairSecp256k1 {
	if len(seed) < SeedSizeSecp256k1 {
		return nil
	}
	prikey := secp256k1.PrivKeyFromBytes(seed[:SeedSizeSecp256k1])
	pubkey := prikey.PubKey()
	return &KeypairSecp256k1{
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func NewKeypairSecp256k1(prikey *secp256k1.PrivateKey, pubkey *secp256k1.PublicKey) *KeypairSecp256k1 {
	return &KeypairSecp256k1{
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func NewSecp256k1SuiSignature(s *Signer, msg []byte) *Secp256k1SuiSignature {
	hash := sha256.Sum256(msg)
	sig := secp256k1_ecdsa.Sign(s.KeypairSecp256k1.PriKey, hash[:])

	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(KeySchemeFlagSecp256k1))
	// if we call func sig.Serialize() it will serialize the signature in DER format.
	// However, Sui requires the signature to be in raw R and S format.
	rawRS, err := dcrecconcatRS(sig)
	if err != nil {
		return nil
	}
	sigBuffer.Write(rawRS[:])
	sigBuffer.Write(s.KeypairSecp256k1.PubKey.SerializeCompressed())

	return &Secp256k1SuiSignature{
		Signature: [SizeSecp256k1SuiSignature]byte(sigBuffer.Bytes()),
	}
}

func dcrecconcatRS(sig *secp256k1_ecdsa.Signature) ([64]byte, error) {
	rawRS := [64]byte{}
	r, err := hex.DecodeString(sig.R().String())
	if err != nil {
		return [64]byte{}, fmt.Errorf("failed to decode R: %w", err)
	}
	s, err := hex.DecodeString(sig.S().String())
	if err != nil {
		return [64]byte{}, fmt.Errorf("failed to decode S: %w", err)
	}
	copy(rawRS[:32], r)
	copy(rawRS[32:], s)
	return rawRS, nil
}
