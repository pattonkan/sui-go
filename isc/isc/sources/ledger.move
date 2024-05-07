/// Module: isc
module isc::ledger {
    use std::ascii::String;

    const EDuplicateNft: u64 = 1;

    /// A `Ledger` keeps track of administration of the objects in an `Assets` object
    /// It can be used to track mutations on the `Assets` object and to describe the
    /// asset allowance for an ISC SC request.
    /// Note that a `Ledger` can be passed around and copied even after the assets it
    /// describes have been moved out of an the `Assets` object.
    public struct Ledger has copy, drop, store {
        token_amounts: vector<u64>,
        token_types: vector<String>,
        nfts: vector<ID>,
    }

    public fun new(): Ledger {
        Ledger {
            token_amounts: vector[],
            token_types: vector[],
            nfts: vector[],
        }
    }

    public fun add_tokens(ledger: &mut Ledger, token_type: &String, amount: u64) {
        let (found, i) = ledger.token_types.index_of(token_type);
        if (found) {
            let tokens = ledger.token_amounts.borrow_mut(i);
            *tokens = *tokens + amount
        } else {
            ledger.token_types.push_back(*token_type);
            ledger.token_amounts.push_back(amount)
        }
    }

    public fun add_nft(ledger: &mut Ledger, nft: ID) {
        assert!(!ledger.nfts.contains(&nft), EDuplicateNft);
        ledger.nfts.push_back(nft)
    }

    public fun get_token_amount(ledger: &Ledger, token_type: &String): u64 {
        let (found, i) = ledger.token_types.index_of(token_type);
        if (!found) {
            return 0
        };
        ledger.token_amounts[i]
    }

    public fun get_token_types(ledger: &Ledger): &vector<String> {
        &ledger.token_types
    }

    public fun get_nfts(ledger: &Ledger): &vector<ID> {
        &ledger.nfts
    }

    public fun has_nft(ledger: &Ledger, nft: ID): bool {
        ledger.nfts.contains(&nft)
    }
}
