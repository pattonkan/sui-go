package multisig

import (
	"fmt"
)

type Verifier struct{}

var ErrKeySchemeMismatch = fmt.Errorf("key scheme mismatch")

func (v *Verifier) VerifyMemberSignature(
	message []byte,
	memberPublicKey *MemberPublicKey,
	memberSignature *MemberSignature,
) error {
	if memberSignature.Ed25519SuiSignature != nil {
		if memberPublicKey.Ed25519PublicKey == nil {
			return ErrKeySchemeMismatch
		}
		valid := memberPublicKey.Ed25519PublicKey.Verify(message, memberSignature.SigOnlyBytes())
		if valid {
			return nil
		}
		return fmt.Errorf("ed25519 signature verification failed")
	}
	if memberSignature.Secp256r1SuiSignature != nil {
		if memberPublicKey.Secp256r1PublicKey == nil {
			return ErrKeySchemeMismatch
		}
		valid := memberPublicKey.Secp256r1PublicKey.Verify(message, memberSignature.SigOnlyBytes())
		if valid {
			return nil
		}
		return fmt.Errorf("secp256r1 signature verification failed")
	}
	if memberSignature.Secp256k1SuiSignature != nil {
		if memberPublicKey.Secp256k1PublicKey == nil {
			return ErrKeySchemeMismatch
		}
		valid := memberPublicKey.Secp256k1PublicKey.Verify(message, memberSignature.SigOnlyBytes())
		if valid {
			return nil
		}
		return fmt.Errorf("secp256k1 signature verification failed")
	}
	return fmt.Errorf("no valid signature")
}

func (v *Verifier) VerifyAggregatedSignature(message []byte, a *AggregatedSignature) error {
	return a.Verify(message)
}
