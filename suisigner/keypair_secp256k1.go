package suisigner

import (
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
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
	// sig := secp256k1.Sign(s.KeypairSecp256k1.PriKey, msg)

	// sigBuffer := bytes.NewBuffer([]byte{})
	// sigBuffer.WriteByte(byte(KeySchemeFlagEd25519))
	// sigBuffer.Write(sig[:])
	// sigBuffer.Write(s.KeypairEd25519.PubKey)

	// return &Ed25519SuiSignature{
	// 	Signature: [SizeEd25519SuiSignature]byte(sigBuffer.Bytes()),
	// }
	panic("TODO")
}
