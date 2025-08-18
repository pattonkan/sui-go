package suicrypto

import (
	"crypto/ed25519"
)

// Sui signatures are composed of a scheme flag, a public key, and a signature.
// The byte layout is as follows:
// | scheme flag | public key in bytes.    | signature in bytes     |
// | 1 byte      | Keypair PublicKey Size  | Keypair Signature Size |
const (
	SizeSuiSignatureEd25519   = KeypairEd25519PublicKeySize + KeypairEd25519SignatureSize + 1
	SizeSuiSignatureSecp256k1 = KeypairSecp256k1PublicKeySize + KeypairSecp256k1SignatureSize + 1
	SizeSuiSignatureSecp256r1 = KeypairSecp256r1PublicKeySize + KeypairSecp256r1SignatureSize + 1
)

const (
	SeedSizeEd25519   = ed25519.SeedSize
	SeedSizeSecp256k1 = 32 // secp256k1 private key is 32 bytes long
)

const (
	KeypairEd25519PublicKeySize = ed25519.PublicKeySize
	KeypairEd25519SignatureSize = ed25519.SignatureSize
)

const (
	KeypairSecp256k1PublicKeySize = 33
	KeypairSecp256k1SignatureSize = 64
)

const (
	KeypairSecp256r1PublicKeySize = 33
	KeypairSecp256r1SignatureSize = 64
)
