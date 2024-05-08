/// Module: isc
module isc::request {
    use isc::{
        ledger::{Ledger},
    };
    use std::ascii::String;

    /// privileged function was called without authorization
    const EMultipleTimeLocks: u64 = 1;

    public struct RequestData has copy, drop, store {
        id: ID,
        contract: String,
        function: String,
        args: vector<vector<u8>>,
        sender: address,
        allowance: Option<Ledger>,
        time_lock: u64,
    }

    public struct Request has key, store {
        id: UID,
        data: RequestData,
    }

    /// creates a request to call a specific SC function
    public fun create_request(contract: String, function: String, args: vector<vector<u8>>, ctx: &mut TxContext): Request {
        let id = object::new(ctx);
        let data = RequestData {
                id: id.uid_to_inner(),
                allowance: option::none(),
                contract: contract,
                function: function,
                args: args,
                time_lock: 0,
                sender: ctx.sender(),
            };
        Request{
            id: id,
            data: move data,
        }
    }

    public fun set_time_lock(req: &mut Request, time_lock: u64) {
        assert!(req.data.time_lock == 0, EMultipleTimeLocks);
        //TODO check time lock value to be in the future?
        req.data.time_lock = time_lock;
    }

    public fun as_data(req: Request): RequestData {
        let Request { id, data } = req;
        object::delete(id);
        move data
    }
}

