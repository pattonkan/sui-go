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
