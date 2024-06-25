package models

import (
	"github.com/howjmay/sui-go/sui_types"
)

type Checkpoint struct {
	Epoch                      *BigInt                          `json:"epoch"`
	SequenceNumber             *BigInt                          `json:"sequenceNumber"`
	Digest                     sui_types.Digest                 `json:"digest"`
	NetworkTotalTransactions   *BigInt                          `json:"networkTotalTransactions"`
	PreviousDigest             *sui_types.Digest                `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary GasCostSummary                   `json:"epochRollingGasCostSummary"`
	TimestampMs                *BigInt                          `json:"timestampMs"`
	Transactions               []*sui_types.Digest              `json:"transactions"`
	CheckpointCommitments      []sui_types.CheckpointCommitment `json:"checkpointCommitments"`
	ValidatorSignature         sui_types.Base64Data             `json:"validatorSignature"`
}

type CheckpointPage = Page[*Checkpoint, BigInt]
