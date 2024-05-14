// Module: sdk_tests
module sdk_tests::sdk_tests {
    use sui::{event, hex};

    public struct InputByteArrayOfArrays has drop, copy {
        data: vector<vector<u8>>,
    }
    public fun input_byte_array_of_arrays(vec: vector<vector<u8>>) {
        // assert!(b"haha" == hex::decode(*vector::borrow(&vec, 0)), 2);
        // assert!(b"abc" == hex::decode(*vector::borrow(&vec, 1)), 1);
        event::emit(InputByteArrayOfArrays { data: vec });
    }

    public struct InputInts has drop, copy {
        input0: u8,
        input1: u16,
        input2: u32,
        input3: u64,
        // FIXME add u128, u256
    }
    // signed ints are not supported
    public fun input_ints( input0: u8, input2: u16, input4: u32, input6: u64) {
        // assert!(b"haha" == vector::borrow(&vec, 0), 1);
        // assert!(b"abc" == vector::borrow(&vec, 1), 1);
        // event::emit(InputInts { data: vec });
    }
}
