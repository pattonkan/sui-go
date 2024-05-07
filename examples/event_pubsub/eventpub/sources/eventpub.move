// Module: eventpub
module eventpub::eventpub {
    use sui::clock::{Self, Clock};
    use sui::event;

    public struct EventContainer has copy, drop {
        timestamp: u64,
    }

    public fun emit_clock(clock: &Clock) {
        event::emit(EventContainer{
            timestamp: clock::timestamp_ms(clock),
        })
    }
}
