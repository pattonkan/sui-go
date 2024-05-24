/// Module: isc
module isc::request {
    use isc::{
        allowance::{Allowance},
    };
    use std::ascii::String;

    public struct RequestData has copy, drop, store {
        id: ID,
        contract: String,
        function: String,
        args: vector<vector<u8>>,
        sender: address,
        allowance: Option<Allowance>,
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
                sender: ctx.sender(),
            };
        Request{
            id: id,
            data: move data,
        }
    }

    /// sets an allowance for the `Request`
    public fun set_allowance(req: &mut Request, allowance: Allowance) {
        req.data.allowance = option::some(allowance);
    }

    public fun destroy_request(req: Request) {
        let Request { id, data: _ } = req;
        object::delete(id)
    }
}

