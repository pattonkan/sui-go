package suicrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

type KeypairSecp256r1 struct {
	PriKey *Secp256r1PriKey
	PubKey *Secp256r1PubKey
}

func NewKeypairSecp256r1FromSeed(seed []byte) *KeypairSecp256r1 {
	curve := elliptic.P256()

	d := new(big.Int).SetBytes(seed[:])

	x, y := curve.ScalarBaseMult(seed)

	priv := ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}
	prikey := Secp256r1PriKey(priv)
	pubkey := Secp256r1PubKey(priv.PublicKey)
	return &KeypairSecp256r1{
		PriKey: &prikey,
		PubKey: &pubkey,
	}
}

func (k *KeypairSecp256r1) Sign(data []byte) ([]byte, error) {
	r, s, err := DeterministicSecp256r1Sign(k.PriKey.Ecdsa(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data by Secp256r1: %w", err)
	}
	sigRawRS := [64]byte{}
	copy(sigRawRS[:32], r.Bytes())
	copy(sigRawRS[32:], s.Bytes())
	return sigRawRS[:], nil
}

type Secp256r1PubKey ecdsa.PublicKey
type Secp256r1PriKey ecdsa.PrivateKey

// This function returns the compressed public key in bytes.
func (p Secp256r1PubKey) Bytes() []byte {
	byteLen := (p.Curve.Params().BitSize + 7) >> 3
	compressed := make([]byte, 1+byteLen)
	if p.Y.Bit(0) == 0 {
		compressed[0] = 0x02
	} else {
		compressed[0] = 0x03
	}
	xBytes := p.X.Bytes()
	copy(compressed[1+byteLen-len(xBytes):], xBytes)
	return compressed
}
func (p Secp256r1PubKey) Ecdsa() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(&p)
}
func (p Secp256r1PriKey) Bytes() []byte {
	panic("not implemented")
}
func (p Secp256r1PriKey) Ecdsa() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(&p)
}
