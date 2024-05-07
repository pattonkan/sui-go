/// Module: isc
module isc::request {
    use isc::{
        ledger::{Ledger},
    };
    use std::ascii::String;

    /// privileged function was called without authorization
    const EMultipleTimeLocks: u64 = 1;

    public struct Request has key, store {
        id: UID,
        contract: String,
        function: String,
        args: vector<vector<u8>>,
        sender: address,
        allowance: Option<Ledger>,
        time_lock: u64,
    }

    /// creates a request to call a specific SC function
    public fun create_request(contract: String, function: String, args: vector<vector<u8>>, ctx: &mut TxContext): Request {
        Request{
            id: object::new(ctx),
            allowance: option::none(),
            contract: contract,
            function: function,
            args: args,
            time_lock: 0,
            sender: ctx.sender(),
        }
    }

    public fun set_time_lock(req: &mut Request, time_lock: u64) {
        assert!(req.time_lock == 0, EMultipleTimeLocks);
        //TODO check time lock value to be in the future?
        req.time_lock = time_lock;
    }
}

