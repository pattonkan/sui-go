#[test_only]
module isc::ledger_test {
    use isc::ledger::{Self,Ledger};
    use sui::sui::SUI;
    use std::type_name;

    #[test]
    fun test_ledger() {
        let mut ctx = tx_context::dummy();

        let uid = object::new(&mut ctx);
        let nft_id = uid.to_inner();
        object::delete(uid);
        
        let uid2 = object::new(&mut ctx);
        let nft_id2 = uid2.to_inner();
        object::delete(uid2);
        
        let mut a = ledger::new();
    
        let name = type_name::get<SUI>().into_string();
        std::debug::print(&name);
        a.add_tokens(&name, 100);

        a.add_nft(nft_id);

        assert!(a.get_token_amount(&name) == 100, 2);
        let dummy = type_name::get<Ledger>().into_string();
        assert!(a.get_token_amount(&dummy) == 0, 3);
        assert!(a.has_nft(nft_id), 4);
        assert!(!a.has_nft(nft_id2), 5);
    }

    #[test, expected_failure(abort_code = isc::ledger::EDuplicateNft)]
    fun test_ledger_duplicate_nft() {
        let mut ctx = tx_context::dummy();

        let uid = object::new(&mut ctx);
        let nft_id = uid.to_inner();
        object::delete(uid);

        let mut a = ledger::new();
        a.add_nft(nft_id);
        a.add_nft(nft_id);
    }
}
