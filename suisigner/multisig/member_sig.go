package multisig

import (
	"fmt"

	"github.com/pattonkan/sui-go/suisigner/suicrypto"
)

type MemberSignature struct {
	Ed25519SuiSignature   *Ed25519SuiSignature   `bcs:"optional"`
	Secp256k1SuiSignature *Secp256k1SuiSignature `bcs:"optional"`
	Secp256r1SuiSignature *Secp256r1SuiSignature `bcs:"optional"`
}

func (s MemberSignature) IsBcsEnum() {}

func MemberSignatureFromBytesEd25519(pubkey *suicrypto.Ed25519PubKey, sigRawBytes []byte) *MemberSignature {
	if len(sigRawBytes) != suicrypto.KeypairEd25519SignatureSize {
		return nil
	}
	var sig Ed25519SuiSignature
	sig[0] = suicrypto.KeySchemeFlagEd25519.Byte()
	copy(sig[1:suicrypto.KeypairEd25519PublicKeySize+1], pubkey.Bytes())
	copy(sig[suicrypto.KeypairEd25519PublicKeySize+1:], sigRawBytes)
	return &MemberSignature{Ed25519SuiSignature: &sig}
}
func MemberSignatureFromBytesSecp256k1(pubkey *suicrypto.Secp256k1PubKey, sigRawBytes []byte) *MemberSignature {
	if len(sigRawBytes) != suicrypto.KeypairSecp256k1SignatureSize {
		return nil
	}
	var sig Secp256k1SuiSignature
	sig[0] = suicrypto.KeySchemeFlagSecp256k1.Byte()
	copy(sig[1:suicrypto.KeypairSecp256k1PublicKeySize+1], pubkey.Bytes())
	copy(sig[suicrypto.KeypairSecp256k1PublicKeySize+1:], sigRawBytes)
	return &MemberSignature{Secp256k1SuiSignature: &sig}
}
func MemberSignatureFromBytesSecp256r1(pubkey *suicrypto.Secp256r1PubKey, sigRawBytes []byte) *MemberSignature {
	if len(sigRawBytes) != suicrypto.KeypairSecp256r1SignatureSize {
		return nil
	}
	var sig Secp256r1SuiSignature
	sig[0] = suicrypto.KeySchemeFlagSecp256r1.Byte()
	copy(sig[1:suicrypto.KeypairSecp256r1PublicKeySize+1], pubkey.Bytes())
	copy(sig[suicrypto.KeypairSecp256r1PublicKeySize+1:], sigRawBytes)
	return &MemberSignature{Secp256r1SuiSignature: &sig}
}

func (s MemberSignature) Bytes() []byte {
	if s.Ed25519SuiSignature != nil {
		return s.Ed25519SuiSignature[:]
	}
	if s.Secp256k1SuiSignature != nil {
		return s.Secp256k1SuiSignature[:]
	}
	if s.Secp256r1SuiSignature != nil {
		return s.Secp256r1SuiSignature[:]
	}
	return nil
}

func (s MemberSignature) SigOnlyBytes() []byte {
	if s.Ed25519SuiSignature != nil {
		return s.Ed25519SuiSignature[suicrypto.KeypairEd25519PublicKeySize+1:]
	}
	if s.Secp256k1SuiSignature != nil {
		return s.Secp256k1SuiSignature[suicrypto.KeypairSecp256k1PublicKeySize+1:]
	}
	if s.Secp256r1SuiSignature != nil {
		return s.Secp256r1SuiSignature[suicrypto.KeypairSecp256r1PublicKeySize+1:]
	}
	return nil
}

type Ed25519SuiSignature [suicrypto.SizeSuiSignatureEd25519]byte

type Secp256k1SuiSignature [suicrypto.SizeSuiSignatureSecp256k1]byte

type Secp256r1SuiSignature [suicrypto.SizeSuiSignatureSecp256r1]byte

func (s MemberSignature) ExtractPubkeyAndSignature() (*MemberPublicKey, []byte, error) {
	if s.Ed25519SuiSignature != nil {
		_, pubkeyBytes, sig, err := suicrypto.ExtractPubkeyAndSignature(s.Ed25519SuiSignature[:])
		if err != nil {
			return nil, nil, fmt.Errorf("can't extract pubkey and signature from bytes: %w", err)
		}

		pubkeyInternal, err := suicrypto.Ed25519PubKeyFromBytes(pubkeyBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert bytes to pubkey: %w", err)
		}
		return &MemberPublicKey{Ed25519PublicKey: pubkeyInternal}, sig, nil
	}

	if s.Secp256k1SuiSignature != nil {
		_, pubkeyBytes, sig, err := suicrypto.ExtractPubkeyAndSignature(s.Secp256k1SuiSignature[:])
		if err != nil {
			return nil, nil, fmt.Errorf("can't extract pubkey and signature from bytes: %w", err)
		}
		pubkeyInternal, err := suicrypto.Secp256k1PubKeyFromBytes(pubkeyBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert bytes to pubkey: %w", err)
		}
		return &MemberPublicKey{Secp256k1PublicKey: pubkeyInternal}, sig, nil
	}

	if s.Secp256r1SuiSignature != nil {
		_, pubkeyBytes, sig, err := suicrypto.ExtractPubkeyAndSignature(s.Secp256r1SuiSignature[:])
		if err != nil {
			return nil, nil, fmt.Errorf("can't extract pubkey and signature from bytes: %w", err)
		}
		pubkeyInternal, err := suicrypto.Secp256r1PubKeyFromBytes(pubkeyBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert bytes to pubkey: %w", err)
		}
		return &MemberPublicKey{Secp256r1PublicKey: pubkeyInternal}, sig, nil
	}
	return nil, nil, fmt.Errorf("no valid signature found")
}
