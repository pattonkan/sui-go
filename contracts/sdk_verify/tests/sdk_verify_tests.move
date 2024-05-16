#[test_only]
module sdk_verify::sdk_verify_tests {
    use sdk_verify::sdk_verify;
    use std::vector;
    // uncomment this line to import the module
    // use sdk_verify::sdk_verify;

    const ENotImplemented: u64 = 0;

    #[test]
    fun test_sdk_verify() {
        let mut vec = vector::empty<vector<u8>>();
        vector::push_back(&mut vec, b"haha");
        vector::push_back(&mut vec, b"gogo");
        sdk_verify::read_input_bytes_array(vec);
    }

    // #[test, expected_failure(abort_code = sdk_verify::sdk_verify_tests::ENotImplemented)]
    // fun test_sdk_verify_fail() {
    //     abort ENotImplemented
    // }
}
