package suisigner

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

// refer sui-rust-sdk/crates/sui-sdk-types/src/crypto/signature.rs
const (
	KeypairSecp256r1PublicKeySize = 33
	KeypairSecp256r1SignatureSize = 64
)

type KeypairSecp256r1 struct {
	PriKey *ecdsa.PrivateKey
	PubKey *ecdsa.PublicKey
}

func NewKeypairSecp256r1FromSeed(seed []byte) *KeypairSecp256r1 {
	curve := elliptic.P256()

	d := new(big.Int).SetBytes(seed[:])

	x, y := curve.ScalarBaseMult(seed)

	priv := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}
	return &KeypairSecp256r1{
		PriKey: priv,
		PubKey: &priv.PublicKey,
	}
}

func NewKeypairSecp256r1(prikey *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey) *KeypairSecp256r1 {
	return &KeypairSecp256r1{
		PriKey: prikey,
		PubKey: pubkey,
	}
}

func NewSecp256r1SuiSignature(signer *Signer, msg []byte) *Secp256r1SuiSignature {
	r, s, err := DeterministicSecp256r1Sign(signer.KeypairSecp256r1.PriKey, msg)
	if err != nil {
		return nil
	}

	sigBuffer := bytes.NewBuffer([]byte{})
	sigBuffer.WriteByte(byte(KeySchemeFlagSecp256r1))

	rawRS := concatRS(r, s)

	sigBuffer.Write(rawRS[:])
	sigBuffer.Write(compressSecp256r1PublicKey(signer.KeypairSecp256r1.PubKey))

	return &Secp256r1SuiSignature{
		Signature: [SizeSecp256r1SuiSignature]byte(sigBuffer.Bytes()),
	}
}

func concatRS(r, s *big.Int) [64]byte {
	rawRS := [64]byte{}
	copy(rawRS[:32], r.Bytes())
	copy(rawRS[32:], s.Bytes())
	return rawRS
}

func compressSecp256r1PublicKey(pub *ecdsa.PublicKey) []byte {
	byteLen := (pub.Curve.Params().BitSize + 7) >> 3
	compressed := make([]byte, 1+byteLen)
	if pub.Y.Bit(0) == 0 {
		compressed[0] = 0x02
	} else {
		compressed[0] = 0x03
	}
	xBytes := pub.X.Bytes()
	copy(compressed[1+byteLen-len(xBytes):], xBytes)
	return compressed
}
