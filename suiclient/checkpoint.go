package suiclient

import "github.com/pattonkan/sui-go/sui"

type (
	CheckpointSequenceNumber = uint64

	CheckpointCommitment    = ECMHLiveObjectSetDigest
	ECMHLiveObjectSetDigest = sui.Digest
)

type Checkpoint struct {
	Epoch                      *sui.BigInt            `json:"epoch"`
	SequenceNumber             *sui.BigInt            `json:"sequenceNumber"`
	Digest                     sui.Digest             `json:"digest"`
	NetworkTotalTransactions   *sui.BigInt            `json:"networkTotalTransactions"`
	PreviousDigest             *sui.Digest            `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary GasCostSummary         `json:"epochRollingGasCostSummary"`
	TimestampMs                *sui.BigInt            `json:"timestampMs"`
	Transactions               []*sui.Digest          `json:"transactions"`
	CheckpointCommitments      []CheckpointCommitment `json:"checkpointCommitments"`
	ValidatorSignature         sui.Base64Data         `json:"validatorSignature"`
}

type CheckpointPage = Page[*Checkpoint, sui.BigInt]
