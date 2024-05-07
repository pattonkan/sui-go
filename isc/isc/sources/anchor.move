
/// Module: isc
module isc::anchor {
    use isc::{
        assets::{Self, Assets},
        request::{Request},
    };
    use stardust::{
        nft::{Nft},
    };
    use std::type_name;
    use sui::{
        bag::{Self, Bag},
        balance::{Balance},
        coin::{Self, Coin},
    };

    /// privileged function was called without authorization
    const EWrongCaller: u64 = 1;

    /// Only one who is allowed to create new `Anchor`s (start new chains).
    public struct AnchorCap has key { id: UID }

    /// Send `AnchorCap` to the publisher
    fun init(ctx: &mut TxContext) {
        transfer::transfer(AnchorCap {
            id: object::new(ctx)
        }, ctx.sender())
    }

    /// Governor of a chain `Anchor`
    public struct GovernorCap has key, store{
        id: UID,
        anchor_id: ID
    }

    public struct Mutation has store {
        sender: address,
        asset_type: vector<u8>,
        amount: u64,
    }

    public struct Anchor has key, store {
        id: UID,
        /// assets controlled by the anchor
        assets: Assets,
        /// pool of mutations to assets that have not been consumed by ISC chain yet
        mutations: vector<Mutation>,
        /// pool of requests that have not been consumed by ISC chain yet
        requests: vector<Request>,
        /// state index
        state_receiveddex: u32,
        // state metadata
        state_metadata: vector<u8>,
        // table that holds all the treasury caps of tokens minted by this chain
        minted_token_treasuries: Bag,
    }

    public struct TokensReceipt<phantom T> has key, store {
        id: UID,
        sender: address,
        tokens: Balance<T>,
    }

    public struct NftReceipt has key, store {
        id: UID,
        sender: address,
        nft: Nft,
    }

    /// starts a new chain
    public fun start_new_chain(_: &AnchorCap, ctx: &mut TxContext): GovernorCap {
        let anchor = Anchor{
            id: object::new(ctx),
            assets: assets::new(ctx),
            mutations: vector[],
            requests: vector[],
            state_receiveddex: 0,
            state_metadata: vector[],
            minted_token_treasuries: bag::new(ctx),
         };

        let governor = GovernorCap{
            id: object::new(ctx),
            anchor_id: anchor.id.uid_to_inner(),
        };

        transfer::share_object(anchor);

        governor
    }
 
    /// client calls this to send tokens to the anchor
    public fun send_tokens<T>(anchor: &mut Anchor, tokens: Coin<T>, ctx: &mut TxContext) {
        transfer::public_transfer(TokensReceipt<T> {
            id: object::new(ctx),
            sender: ctx.sender(),
            tokens: tokens.into_balance(),
        }, anchor.id.to_address());
    }
 
    /// client calls this to have the anchor receive tokens and add them to its assets
    public fun receive_tokens<T>(anchor: &mut Anchor, tokens_receipt: transfer::Receiving<TokensReceipt<T>>, _ctx: &mut TxContext) {
        let TokensReceipt<T> { id, sender, tokens } = transfer::public_receive(&mut anchor.id, tokens_receipt);
        let token_type_name = type_name::get<T>().into_string();
        anchor.mutations.push_back(Mutation{
            sender: sender,
            asset_type: token_type_name.into_bytes(),
            amount: tokens.value(),
        });
        anchor.assets.add_tokens(tokens);
        object::delete(id);
    }
 
    public fun send_nft(anchor: &mut Anchor, nft: Nft, ctx: &mut TxContext) {
       transfer::public_transfer(NftReceipt {
            id: object::new(ctx),
            sender: ctx.sender(),
            nft: nft,
        }, anchor.id.to_address());
     }
 
    public fun receive_nft(anchor: &mut Anchor, nft_received: transfer::Receiving<NftReceipt>, _ctx: &mut TxContext) {
        let NftReceipt { id, sender, nft } = transfer::public_receive(&mut anchor.id, nft_received);
        let nft_id = object::borrow_id<Nft>(&nft);
        anchor.mutations.push_back(Mutation{
            sender: sender,
            asset_type: nft_id.id_to_bytes(),
            amount: 0,
        });
        anchor.assets.add_nft(nft);
        object::delete(id);
   }

    /// client calls this to send a request to the anchor
    public fun send_request(anchor: &Anchor, req: Request, _ctx: &mut TxContext) {
        // we could implement specific checks here for the request

        // send the request object to the chain's id
        transfer::public_transfer(req, anchor.id.to_address());
    }

    /// client calls this to have the anchor receive a request and add it to the request pool
    public fun receive_request(anchor: &mut Anchor, request: transfer::Receiving<Request>, _ctx: &mut TxContext){
        let req = transfer::public_receive(&mut anchor.id, request);
        anchor.requests.push_back<Request>(req);
    }


    fun check_governor(anchor: &Anchor, governor: &GovernorCap) {
        assert!(governor.anchor_id == anchor.id.uid_to_inner(), EWrongCaller);
     }

    /// chain governor calls this to transfer tokens from its `Assets` to an address
    public fun transfer_tokens<T>(anchor: &mut Anchor, governor: &GovernorCap, to: address, amount: u64, ctx: &mut TxContext){
        check_governor(anchor, governor);
        let tokens = anchor.assets.take_tokens<T>(amount);
        transfer::public_transfer(coin::from_balance(tokens, ctx), to);
    }

    /// chain governor calls this to transfer an nft from its `Assets` to an address
    public fun transfer_nft(anchor: &mut Anchor, governor: &GovernorCap, to: address, nft_id: ID, _ctx: &mut TxContext){
        check_governor(anchor, governor);
        let nft = anchor.assets.take_nft(nft_id);
        transfer::public_transfer(nft, to);
    }
}
