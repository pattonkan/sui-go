#[test_only]
module contract_tests::contract_tests_tests {
    use contract_tests::contract_tests;
    use std::vector;
    // uncomment this line to import the module
    // use contract_tests::contract_tests;

    const ENotImplemented: u64 = 0;

    #[test]
    fun test_contract_tests() {
        let mut vec = vector::empty<vector<u8>>();
        vector::push_back(&mut vec, b"haha");
        vector::push_back(&mut vec, b"gogo");
        contract_tests::read_input_bytes_array(vec);
    }

    // #[test, expected_failure(abort_code = contract_tests::contract_tests_tests::ENotImplemented)]
    // fun test_contract_tests_fail() {
    //     abort ENotImplemented
    // }
}
