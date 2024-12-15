package movebcs

import "github.com/howjmay/sui-go/sui"

// / Sender is not @0x0 the system address.
const EClockNotSystemAddress = 0

type Clock struct {
	Id *sui.ObjectId
	// The clock's timestamp, which is set automatically by a
	// system transaction every time consensus commits a
	// schedule, or by `sui::clock::increment_for_testing` during
	// testing.
	TimestampMs uint64
}
