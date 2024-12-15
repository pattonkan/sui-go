package movebcs

const (
	// For when trying to destroy a non-zero balance.
	EBalanceNonZero = 0
	// For when an overflow is happening on Supply operations.
	EBalanceOverflow = 1
	// For when trying to withdraw more than there is.
	EBalanceNotEnough = 2
	// Sender is not @0x0 the system address.
	EBalanceNotSystemAddress = 3
	// System operation performed for a coin other than SUI
	EBalanceNotSUI = 4
)

type MoveBalance struct {
	Value uint64
}

type MoveSupply struct {
	Value uint64
}
