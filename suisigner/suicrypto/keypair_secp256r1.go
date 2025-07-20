package suicrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io"
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

func Secp256r1PubKeyFromBytes(b []byte) (*Secp256r1PubKey, error) {
	curve := elliptic.P256()

	byteLen := (curve.Params().BitSize + 7) >> 3

	if len(b) != 1+byteLen {
		return nil, fmt.Errorf("invalid compressed public key length")
	}

	prefix := b[0]
	if prefix != 0x02 && prefix != 0x03 {
		return nil, fmt.Errorf("invalid compressed public key prefix")
	}

	x := new(big.Int).SetBytes(b[1:])

	// y² = x³ - 3x + b mod p
	curveParams := curve.Params()
	x3 := new(big.Int).Exp(x, big.NewInt(3), curveParams.P) // x³
	threeX := new(big.Int).Mul(big.NewInt(3), x)            // 3x
	x3.Sub(x3, threeX)                                      // x³ - 3x
	x3.Add(x3, curveParams.B)                               // x³ - 3x + b
	x3.Mod(x3, curveParams.P)

	// Compute sqrt(y^2) mod p (Tonelli-Shanks)
	y := new(big.Int).ModSqrt(x3, curveParams.P)
	if y == nil {
		return nil, fmt.Errorf("invalid point: sqrt does not exist")
	}

	// If odd/even flag doesn't match, flip y
	if y.Bit(0) != uint(prefix&1) {
		y.Sub(curveParams.P, y)
		y.Mod(y, curveParams.P)
	}

	pubkey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	return (*Secp256r1PubKey)(pubkey), nil
}

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
func (p Secp256r1PubKey) String() string {
	return hex.EncodeToString(p.Bytes())
}
func (p Secp256r1PubKey) Ecdsa() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(&p)
}
func (p Secp256r1PubKey) MarshalBCS() ([]byte, error) {
	return p.Bytes(), nil
}
func (p *Secp256r1PubKey) UnmarshalBCS(r io.Reader) (int, error) {
	buf := make([]byte, KeypairSecp256r1PublicKeySize)
	n, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	_p, err := Secp256r1PubKeyFromBytes(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to convert bytes to Secp256r1PubKey: %w", err)
	}
	*p = *_p
	return n, nil
}
func (p Secp256r1PriKey) Bytes() []byte {
	panic("not implemented")
}
func (p Secp256r1PriKey) Ecdsa() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(&p)
}

func (p Secp256r1PubKey) Verify(data []byte, sig []byte) bool {
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	return VerifySecp256r1(p.Ecdsa(), data, r, s)
}
