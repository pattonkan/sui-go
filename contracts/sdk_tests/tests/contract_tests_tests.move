#[test_only]
module sdk_tests::sdk_tests_tests {
    use sdk_tests::sdk_tests;
    use std::vector;
    // uncomment this line to import the module
    // use sdk_tests::sdk_tests;

    const ENotImplemented: u64 = 0;

    #[test]
    fun test_sdk_tests() {
        let mut vec = vector::empty<vector<u8>>();
        vector::push_back(&mut vec, b"haha");
        vector::push_back(&mut vec, b"gogo");
        sdk_tests::read_input_bytes_array(vec);
    }

    // #[test, expected_failure(abort_code = sdk_tests::sdk_tests_tests::ENotImplemented)]
    // fun test_sdk_tests_fail() {
    //     abort ENotImplemented
    // }
}
