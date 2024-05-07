package sui_types

import "github.com/howjmay/sui-go/sui_types/serialization"

type Digest = serialization.Base58
type ObjectDigest = Digest
type TransactionDigest = Digest
type TransactionEffectsDigest = Digest
type TransactionEventsDigest = Digest
type CheckpointDigest = Digest
type CertificateDigest = Digest
type CheckpointContentsDigest = Digest

func NewDigest(str string) (*Digest, error) {
	return serialization.NewBase58(str)
}

func NewDigestMust(str string) *Digest {
	digest, err := serialization.NewBase58(str)
	if err != nil {
		panic(err)
	}
	return digest
}
