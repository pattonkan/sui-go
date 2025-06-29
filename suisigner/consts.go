package suisigner

import "crypto/ed25519"

const (
	SeedSizeEd25519   = ed25519.SeedSize
	SeedSizeSecp256k1 = 32 // secp256k1 private key is 32 bytes long
)
