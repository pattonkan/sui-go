package multisig

import (
	"fmt"
	"math/bits"
)

// pub struct MultisigAggregatedSignature {
//     /// The plain signature encoded with signature scheme.
//     ///
//     /// The signatures must be in the same order as they are listed in the committee.
//     #[cfg_attr(feature = "proptest", any(proptest::collection::size_range(0..=10).lift()))]
//     signatures: Vec<MultisigMemberSignature>,
//     /// A bitmap that indicates the position of which public key the signature should be
//     /// authenticated with.
//     bitmap: BitmapUnit,
//     /// Legacy encoding for the bitmap.
//     //TODO implement a strategy for legacy bitmap
//     #[cfg_attr(feature = "proptest", strategy(proptest::strategy::Just(None)))]
//     legacy_bitmap: Option<roaring::RoaringBitmap>,
//     /// The public key encoded with each public key with its signature scheme used along with the
//     /// corresponding weight.
//     committee: MultisigCommittee,
// }

type AggregatedSignature struct {
	Signatures []*MemberSignature
	Bitmap     BitmapUnit
	Committee  *Committee
}

func NewAggregatedSignature() *AggregatedSignature {
	return &AggregatedSignature{
		Signatures: []*MemberSignature{},
		Bitmap:     0,
		Committee: &Committee{
			Members:   []*Member{},
			Threshold: 0,
		},
	}
}

func (a *AggregatedSignature) Verify(message []byte) error {
	if !a.Committee.IsValid() {
		return fmt.Errorf("invalid MultisigCommittee")
	}

	if len(a.Signatures) != bits.OnesCount16(a.Bitmap) {
		return fmt.Errorf("number of signatures does not match bitmap")
	}

	if len(a.Signatures) > len(a.Committee.Members) {
		return fmt.Errorf("number of signatures does not match bitmap")
	}

	// calculate the sum of weights of all the members in committee
	v := Verifier{}
	var sum uint16 = 0
	var sigCount = -1
	for i := 0; i < bitmapSize; i++ {
		if a.Bitmap&(1<<i) == 0 {
			continue
		}

		sigCount += 1
		err := v.VerifyMemberSignature(message, &a.Committee.Members[i].PublicKey, a.Signatures[sigCount])
		if err == ErrKeySchemeMismatch {
			continue
		}
		if err != nil {
			return fmt.Errorf("member's sig is invalid: %w", err)
		}

		sum += uint16(a.Committee.Members[i].Weight)
	}

	if sum < a.Committee.Threshold {
		return fmt.Errorf("signature weight does not exceed threshold")
	}

	return nil
}
