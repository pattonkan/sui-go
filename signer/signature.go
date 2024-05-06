package signer

import "crypto/ed25519"

type Signature struct {
	*Ed25519SuiSignature
	*Secp256k1SuiSignature
	*Secp256r1SuiSignature
}

type Secp256k1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Secp256r1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Ed25519SuiSignature struct {
	Signature [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
}
