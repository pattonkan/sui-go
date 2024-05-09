/// Module: isc
module isc::allowance {
    use std::ascii::String;

    const EDuplicateNft: u64 = 1;

    /// An `Allowance` can be used to limit the assets that can be accessed by a `Request`
    public struct Allowance has copy, drop, store {
        coin_amounts: vector<u64>,
        coin_types: vector<String>,
        nfts: vector<ID>,
    }

    public fun new(): Allowance {
        Allowance {
            coin_amounts: vector[],
            coin_types: vector[],
            nfts: vector[],
        }
    }

    public fun add_coin(allowance: &mut Allowance, coin_type: &String, amount: u64) {
        let (found, i) = allowance.coin_types.index_of(coin_type);
        if (found) {
            let coin_amount = allowance.coin_amounts.borrow_mut(i);
            *coin_amount = *coin_amount + amount
        } else {
            allowance.coin_types.push_back(*coin_type);
            allowance.coin_amounts.push_back(amount)
        }
    }

    public fun add_nft(allowance: &mut Allowance, nft: ID) {
        assert!(!allowance.nfts.contains(&nft), EDuplicateNft);
        allowance.nfts.push_back(nft)
    }

    public fun get_coin_amount(allowance: &Allowance, coin_type: &String): u64 {
        let (found, i) = allowance.coin_types.index_of(coin_type);
        if (!found) {
            return 0
        };
        allowance.coin_amounts[i]
    }

    public fun get_coin_types(allowance: &Allowance): &vector<String> {
        &allowance.coin_types
    }

    public fun get_nfts(allowance: &Allowance): &vector<ID> {
        &allowance.nfts
    }

    public fun has_nft(allowance: &Allowance, nft: ID): bool {
        allowance.nfts.contains(&nft)
    }
}
