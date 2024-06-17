// Module: sdk_verify
module sdk_verify::sdk_verify {
    use std::debug;
    use sui::event;

    public struct ReadInputBytesArrayEvent has drop, copy {
        data: vector<vector<u8>>,
    }
    public fun read_input_bytes_array(vec: vector<vector<u8>>) {
        debug::print(&vec);
        assert!(b"haha" == vector::borrow(&vec, 0), 1);
        assert!(b"gogo" == vector::borrow(&vec, 1), 1);
        event::emit(ReadInputBytesArrayEvent { data: vec });
    }

    public fun ret_two_1(): (u64, u32) {
        (1, 2)
    }
    public fun ret_two_2(arg0: u32, arg1: u64): (u32, u64) {
        (arg0, arg1)
    }
}
