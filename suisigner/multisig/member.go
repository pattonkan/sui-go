package multisig

import "github.com/pattonkan/sui-go/suisigner/suicrypto"

type Member struct {
	PublicKey MemberPublicKey
	Weight    WeightUnit
}

func MemberMember(publicKey MemberPublicKey, weight WeightUnit) *Member {
	return &Member{
		PublicKey: publicKey,
		Weight:    weight,
	}
}

type MemberPublicKey struct {
	Ed25519PublicKey   *suicrypto.Ed25519PubKey   `bcs:"optional"`
	Secp256k1PublicKey *suicrypto.Secp256k1PubKey `bcs:"optional"`
	Secp256r1PublicKey *suicrypto.Secp256r1PubKey `bcs:"optional"`
	// ZkLoginPublicKey
	// PasskeyPublicKey
}

func (m MemberPublicKey) IsBcsEnum() {}

func (m MemberPublicKey) Bytes() []byte {
	if m.Ed25519PublicKey != nil {
		return m.Ed25519PublicKey.Bytes()
	}

	if m.Secp256k1PublicKey != nil {
		return m.Secp256k1PublicKey.Bytes()
	}

	if m.Secp256r1PublicKey != nil {
		return m.Secp256r1PublicKey.Bytes()
	}

	// ZkLoginPublicKey
	// PasskeyPublicKey
	return nil
}

// return public key encoded in hex string
func (m MemberPublicKey) String() string {
	if m.Ed25519PublicKey != nil {
		return m.Ed25519PublicKey.String()
	}

	if m.Secp256k1PublicKey != nil {
		return m.Secp256k1PublicKey.String()
	}

	if m.Secp256r1PublicKey != nil {
		return m.Secp256r1PublicKey.String()
	}

	// ZkLoginPublicKey
	// PasskeyPublicKey
	return ""
}
