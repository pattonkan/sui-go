/// Module: isc
module isc::assets {
    use stardust::{
        nft::{Nft},
    };
    use std::ascii::String;
    use std::type_name;
    use sui::{
        bag::{Self, Bag},
        balance::{Balance},
    };

    const ETokenNotFound: u64 = 1;
    const ENftNotFound: u64 = 2;
    const EInsufficientBalance: u64 = 3;

    /// keeps track of all assets in an `Anchor`
    public struct Assets has key, store {
        id: UID,
        tokens: Bag,
        nfts: vector<Nft>,
    }

    /// create empty container for assets
    public fun new(ctx: &mut TxContext): Assets {
        Assets {
            id: object::new(ctx),
            tokens: bag::new(ctx),
            nfts: vector[],
        }
    }

    /// add native tokens to the `Assets` 
    /// consolidates balances of the same T by `join`ing them together
    public fun add_tokens<T>(assets: &mut Assets, tokens: Balance<T>) {
        let token_type_name = type_name::get<T>().into_string();
        if (!assets.tokens.contains(token_type_name)) {
            assets.tokens.add(token_type_name, tokens);
        } else {
            assets.tokens.borrow_mut<String, Balance<T>>(token_type_name).join(tokens);
        }
     }

    /// add an `Nft` to the `Assets`
    public fun add_nft(assets: &mut Assets, nft: Nft) {
        assets.nfts.push_back(nft);
    }

    /// local helper, finds a specified nft in the `Assets`
    fun find_nft(assets: &Assets, nft_id: ID): (bool, u64) {
        let mut found = false;
        let mut i = assets.nfts.length();
        while (i > 0 && !found) {
            i = i - 1;
            let nft = assets.nfts.borrow(i);
            let id = object::borrow_id<Nft>(nft);
            found = (id == nft_id)
        };
        (found, i)
    }

    /// determines if a specified token T is present in the `Assets`
    public fun has_tokens<T>(assets: &mut Assets): bool {
        let token_type_name = type_name::get<T>().into_string();
        assets.tokens.contains(token_type_name)
    }

    /// determines if a specified nft is present in the `Assets`
    public fun has_nft(assets: &mut Assets, nft_id: ID): bool {
        let (found, _) = assets.find_nft(nft_id);
        found
    }

    /// takes a specified amount of native tokens from the `Assets`
    public fun take_tokens<T>(assets: &mut Assets, amount: u64): Balance<T> {
        let token_type_name = type_name::get<T>().into_string();
        assert!(assets.tokens.contains(token_type_name), ETokenNotFound);
        let tokens = assets.tokens.borrow_mut<String, Balance<T>>(token_type_name);
        assert!(tokens.value() >= amount, EInsufficientBalance);
        if (tokens.value() == amount) {
            return assets.tokens.remove(token_type_name)
        };
        tokens.split(amount)
    }

    /// takes a specified nft from the `Assets`
    public fun take_nft(assets: &mut Assets, nft_id: ID): Nft {
        let (found, i) = assets.find_nft(nft_id);
        assert!(found, ENftNotFound);
        assets.nfts.swap_remove(i)
    }

    public fun destroy_empty(assets: Assets) {
        let Assets {
            id,
            tokens,
            nfts,
        } = assets;
        tokens.destroy_empty();
        nfts.destroy_empty();
        object::delete(id);
    }
}
