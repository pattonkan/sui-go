package multisig

import (
	"fmt"

	"github.com/pattonkan/sui-go/sui/suiptb"
)

type Aggregator struct {
	Committee    *Committee
	Signatures   map[uint16]*MemberSignature
	SignedWeight uint16
	Message      []byte
	Verifier     *Verifier
}

func NewAggregatorWithTransaction(
	committee *Committee,
	tx *suiptb.TransactionData,
) (*Aggregator, error) {
	digest, err := tx.SigningDigest()
	if err != nil {
		return nil, err
	}
	return &Aggregator{
		Committee:    committee,
		Signatures:   make(map[uint16]*MemberSignature),
		SignedWeight: 0,
		Message:      digest,
		Verifier:     &Verifier{},
	}, nil
}

// func NewAggregatorWithMessage(message []byte) *Aggregator {
// 	return &Aggregator{}
// }

// pub fn add_signature(&mut self, signature: UserSignature) -> Result<(), SignatureError> {
//     use std::collections::btree_map::Entry;

//     let (public_key, signature) = multisig_pubkey_and_signature_from_user_signature(signature)?;
//     let member_idx = self
//         .committee
//         .members()
//         .iter()
//         .position(|member| member.public_key() == &public_key)
//         .ok_or_else(|| {
//             SignatureError::from_source(
//                 "provided signature does not belong to committee member",
//             )
//         })?;

//     self.verifier()
//         .verify_member_signature(&self.message, &public_key, &signature)?;

//     match self.signatures.entry(member_idx) {
//         Entry::Vacant(v) => {
//             v.insert(signature);
//         }
//         Entry::Occupied(_) => {
//             return Err(SignatureError::from_source(
//                 "duplicate signature from same committee member",
//             ))
//         }
//     }

//     self.signed_weight += self.committee.members()[member_idx].weight() as u16;

//     Ok(())
// }

func (a *Aggregator) AddSignature(signature *MemberSignature) error {
	pubkey, _, err := signature.ExtractPubkeyAndSignature()
	if err != nil {
		return fmt.Errorf("failed to extract pubkey and signature: %w", err)
	}

	memberIdx := a.Committee.IndexOf(pubkey.Bytes())
	if memberIdx == -1 {
		return fmt.Errorf("signature does not belong to committee member")
	}

	err = a.Verifier.VerifyMemberSignature(a.Message, pubkey, signature)
	if err != nil {
		return fmt.Errorf("invalid member signature: %w", err)
	}
	_, ok := a.Signatures[uint16(memberIdx)]
	if ok {
		return fmt.Errorf("duplicate signature from same committee member")
	}
	a.Signatures[uint16(memberIdx)] = signature

	a.SignedWeight += uint16(a.Committee.Members[memberIdx].Weight)
	return nil
}

func (a *Aggregator) Finish() (*AggregatedSignature, error) {
	if a.SignedWeight < a.Committee.Threshold {
		return nil, fmt.Errorf("insufficient signature weight to reach threshold")
	}

	signatures := make([]*MemberSignature, 0, len(a.Signatures))
	bitmap := BitmapUnit(0)
	for memberIdx, signature := range a.Signatures {
		signatures = append(signatures, signature)
		bitmap |= 1 << memberIdx
	}

	return &AggregatedSignature{
		Signatures: signatures,
		Bitmap:     bitmap,
		Committee:  a.Committee,
	}, nil
}
