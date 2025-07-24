package multisig

type ThresholdUnit = uint16
type WeightUnit = uint8
type BitmapUnit = uint16

const bitmapSize = 16

// MAX_COMMITTEE_SIZE in 'crates/sui-sdk-types/src/crypto/multisig.rs'
const maxCommitteeSize = 10
