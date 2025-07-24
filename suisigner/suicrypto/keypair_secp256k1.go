package suicrypto

import (
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

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

func Secp256k1PubKeyFromBytes(b []byte) (*Secp256k1PubKey, error) {
	if len(b) != KeypairSecp256k1PublicKeySize {
		return nil, fmt.Errorf("invalid public key size")
	}
	pubkey, err := secp256k1.ParsePubKey(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return (*Secp256k1PubKey)(pubkey), nil
}

// This function returns the compressed public key in bytes.
func (p Secp256k1PubKey) Bytes() []byte {
	return p.toDcrdPubKey().SerializeCompressed()
}
func (p Secp256k1PubKey) String() string {
	return hex.EncodeToString(p.Bytes())
}
func (p Secp256k1PubKey) MarshalBCS() ([]byte, error) {
	return p.Bytes(), nil
}
func (p *Secp256k1PubKey) UnmarshalBCS(r io.Reader) (int, error) {
	buf := make([]byte, KeypairSecp256k1PublicKeySize)
	n, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	_p, err := Secp256k1PubKeyFromBytes(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to convert bytes to Secp256k1PubKey: %w", err)
	}
	*p = *_p
	return n, nil
}

func (p Secp256k1PriKey) Bytes() []byte {
	return p.toDcrdPriKey().Serialize()
}

func (p Secp256k1PubKey) toDcrdPubKey() *secp256k1.PublicKey {
	return (*secp256k1.PublicKey)(&p)
}
func (p Secp256k1PriKey) toDcrdPriKey() *secp256k1.PrivateKey {
	return (*secp256k1.PrivateKey)(&p)
}

func (p Secp256k1PubKey) Verify(data []byte, sig []byte) bool {
	hash := sha256.Sum256(data)

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	der, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})
	if err != nil {
		return false
	}

	// Verify the signature using the public key and the hash of the data.
	derSig, err := secp256k1_ecdsa.ParseDERSignature(der)
	if err != nil {
		return false
	}
	return derSig.Verify(hash[:], p.toDcrdPubKey())
}
