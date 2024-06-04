package sui_types

type (
	CheckpointSequenceNumber = uint64
	CheckpointTimestamp      = uint64

	CheckpointCommitment    = ECMHLiveObjectSetDigest
	ECMHLiveObjectSetDigest = Digest
)
