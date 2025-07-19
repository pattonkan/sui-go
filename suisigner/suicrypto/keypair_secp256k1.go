package suicrypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp256k1_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

type KeypairSecp256k1 struct {
	PriKey *Secp256k1PriKey
	PubKey *Secp256k1PubKey
}

func NewKeypairSecp256k1FromSeed(seed []byte) *KeypairSecp256k1 {
	if len(seed) < SeedSizeSecp256k1 {
		return nil
	}

	dcrdPrivKey := secp256k1.PrivKeyFromBytes(seed[:SeedSizeSecp256k1])
	prikey := Secp256k1PriKey(*dcrdPrivKey)
	pubkey := Secp256k1PubKey(*dcrdPrivKey.PubKey())
	return &KeypairSecp256k1{
		PriKey: &prikey,
		PubKey: &pubkey,
	}
}

// The output is the signature in raw R and S format, not DER format.
func (k *KeypairSecp256k1) Sign(msg []byte) ([]byte, error) {
	hash := sha256.Sum256(msg)
	sig := secp256k1_ecdsa.Sign(k.PriKey.toDcrdPriKey(), hash[:])

	// if we call func sig.Serialize() it will serialize the signature in DER format.
	// However, Sui requires the signature to be in raw R and S format.
	// Therefore, we convert R and S into concatenated raw R and S format.
	rBytes, err := hex.DecodeString(sig.R().String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode R: %w", err)
	}
	sBytes, err := hex.DecodeString(sig.S().String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode S: %w", err)
	}
	sigRawRS := [64]byte{}
	copy(sigRawRS[:32], rBytes)
	copy(sigRawRS[32:], sBytes)
	return sigRawRS[:], nil
}

type Secp256k1PriKey secp256k1.PrivateKey
type Secp256k1PubKey secp256k1.PublicKey

// This function returns the compressed public key in bytes.
func (p Secp256k1PubKey) Bytes() []byte {
	return p.toDcrdPubKey().SerializeCompressed()
}
func (p Secp256k1PriKey) Bytes() []byte {
	return p.toDcrdPriKey().Serialize()
}

func (p Secp256k1PriKey) toDcrdPriKey() *secp256k1.PrivateKey {
	return (*secp256k1.PrivateKey)(&p)
}
func (p Secp256k1PubKey) toDcrdPubKey() *secp256k1.PublicKey {
	return (*secp256k1.PublicKey)(&p)
}

// func NewSecp256k1SuiSignature(s *Signer, msg []byte) *Secp256k1SuiSignature {
// 	hash := sha256.Sum256(msg)
// 	sig := secp256k1_ecdsa.Sign(s.KeypairSecp256k1.PriKey, hash[:])

// 	sigBuffer := bytes.NewBuffer([]byte{})
// 	sigBuffer.WriteByte(byte(KeySchemeFlagSecp256k1))
// 	// if we call func sig.Serialize() it will serialize the signature in DER format.
// 	// However, Sui requires the signature to be in raw R and S format.
// 	rawRS, err := dcrecconcatRS(sig)
// 	if err != nil {
// 		return nil
// 	}
// 	sigBuffer.Write(rawRS[:])
// 	sigBuffer.Write(s.KeypairSecp256k1.PubKey.SerializeCompressed())

// 	return &Secp256k1SuiSignature{
// 		Signature: [SizeSecp256k1SuiSignature]byte(sigBuffer.Bytes()),
// 	}
// }
