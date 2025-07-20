package multisig

import (
	"bytes"
)

// pub struct MultisigCommittee {
//     /// A list of committee members and their corresponding weight.
//     #[cfg_attr(feature = "proptest", any(proptest::collection::size_range(0..=10).lift()))]
//     members: Vec<MultisigMember>,

//	    /// If the total weight of the public keys corresponding to verified signatures is larger than
//	    /// threshold, the Multisig is verified.
//	    threshold: ThresholdUnit,
//	}

type Committee struct {
	Members   []*Member
	Threshold ThresholdUnit
}

func NewCommittee(members []*Member, threshold ThresholdUnit) *Committee {
	return &Committee{
		Members:   members,
		Threshold: threshold,
	}
}

func (c *Committee) IsValid() bool {
	if c.Threshold == 0 {
		return false
	}
	if len(c.Members) == 0 || len(c.Members) > maxCommitteeSize {
		return false
	}

	sum := WeightUnit(0)
	for _, member := range c.Members {
		if member.Weight == 0 {
			return false
		}
		sum += member.Weight
	}
	if ThresholdUnit(sum) < c.Threshold {
		return false
	}

	m := make(map[string]bool)
	for _, member := range c.Members {
		m[member.PublicKey.String()] = true
	}
	return len(m) == len(c.Members)
}

func (c *Committee) ContainsMember(pubkey []byte) bool {
	for _, member := range c.Members {
		if bytes.Equal(member.PublicKey.Bytes(), pubkey) {
			return true
		}
	}
	return false
}

func (c *Committee) IndexOf(pubkey []byte) int8 {
	for i, member := range c.Members {
		if bytes.Equal(member.PublicKey.Bytes(), pubkey) {
			return int8(i)
		}
	}
	return -1
}
