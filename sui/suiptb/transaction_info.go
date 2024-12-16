package suiptb

import (
	"github.com/pattonkan/sui-go/sui"
)

type GasData struct {
	Payment []*sui.ObjectRef
	Owner   *sui.Address
	Price   uint64
	Budget  uint64
}

type TransactionExpiration struct {
	None  *sui.EmptyEnum
	Epoch *sui.EpochId
}

func (t TransactionExpiration) IsBcsEnum() {}

type CheckpointTimestamp uint64

type ConsensusCommitPrologue struct {
	Epoch             uint64
	Round             uint64
	CommitTimestampMs CheckpointTimestamp
}

type ChangeEpoch struct {
	Epoch                   sui.EpochId
	ProtocolVersion         sui.ProtocolVersion
	StorageCharge           uint64
	ComputationCharge       uint64
	StorageRebate           uint64
	NonRefundableStorageFee uint64
	EpochStartTimestampMs   uint64
	SystemPackages          []*struct {
		SequenceNumber sui.SequenceNumber
		Bytes          [][]uint8
		Objects        []*sui.ObjectId
	}
}

type GenesisTransaction struct {
	Objects []GenesisObject
}

type GenesisObject struct {
	RawObject *struct {
		Data  sui.Data
		Owner sui.Owner
	}
}
