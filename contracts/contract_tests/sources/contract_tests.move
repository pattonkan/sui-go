// Module: contract_tests
module contract_tests::contract_tests {
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
}
