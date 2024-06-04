package models

import (
	"github.com/howjmay/sui-go/sui_types"
)

type Checkpoint struct {
	Epoch                      SafeSuiBigInt[EpochId]                  `json:"epoch"`
	SequenceNumber             SafeSuiBigInt[sui_types.SequenceNumber] `json:"sequenceNumber"`
	Digest                     sui_types.Digest                        `json:"digest"`
	NetworkTotalTransactions   SafeSuiBigInt[uint64]                   `json:"networkTotalTransactions"`
	PreviousDigest             *sui_types.Digest                       `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary GasCostSummary                          `json:"epochRollingGasCostSummary"`
	TimestampMs                SafeSuiBigInt[uint64]                   `json:"timestampMs"`
	Transactions               []*sui_types.Digest                     `json:"transactions"`
	CheckpointCommitments      []sui_types.CheckpointCommitment        `json:"checkpointCommitments"`
	ValidatorSignature         sui_types.Base64Data                    `json:"validatorSignature"`
}

type CheckpointPage = Page[*Checkpoint, SafeSuiBigInt[uint64]]
