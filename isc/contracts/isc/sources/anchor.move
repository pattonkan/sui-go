
/// Module: isc
module isc::anchor {
    use isc::{
        assets::{Self, Assets},
        request::{Request, RequestData},
    };
    use stardust::{
        nft::{Nft},
    };
    use std::type_name;
    use sui::{
        balance::{Balance},
        coin::{Self, Coin},
    };

    /// privileged function was called without authorization
    const EWrongCaller: u64 = 1;
    const EInvalidMutationID: u64 = 2;
    const EInvalidRequestID: u64 = 3;

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

    /// keeps track of an asset mutation
    public struct Mutation has copy, drop, store {
        /// sender of the asset
        sender: address,
        /// byte representation if the asset:
        /// - type name T for Coint<T>
        /// - NFT id for NFT
        asset_type: vector<u8>,
        /// amount of tokens, or zero when this is an NFT
        amount: u64,
    }

    public struct Anchor has key, store {
        id: UID,
        /// assets controlled by the anchor
        assets: Assets,
        /// pool of mutations to assets that have not been consumed by ISC chain yet
        mutations: vector<Mutation>,
        /// id of the next mutation, can be used to deduce the id of mutations[0]
        mutation_id: u64,
        /// pool of requests that have not been consumed by ISC chain yet
        requests: vector<RequestData>,
        /// id of the next request, can be used to deduce the id of requests[0]
        request_id: u64,
        // /// state index
        // state_receiveddex: u32,
        // // state metadata
        // state_metadata: vector<u8>,
        // // table that holds all the treasury caps of tokens minted by this chain
        // minted_token_treasuries: Bag,
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

    /// starts a new chain by creating a new `Anchor` for it
    public fun start_new_chain(_: &AnchorCap, ctx: &mut TxContext): GovernorCap {
        let anchor = Anchor{
            id: object::new(ctx),
            assets: assets::new(ctx),
            mutations: vector[],
            mutation_id: 0,
            requests: vector[],
            request_id: 0,
            // state_receiveddex: 0,
            // state_metadata: vector[],
            // minted_token_treasuries: bag::new(ctx),
         };

        let governor = GovernorCap{
            id: object::new(ctx),
            anchor_id: anchor.id.uid_to_inner(),
        };

        transfer::share_object(anchor);

        governor
    }
 
    /// client calls this to send tokens to the `Anchor`
    public fun send_tokens<T>(anchor: &mut Anchor, tokens: Coin<T>, ctx: &mut TxContext) {
        transfer::public_transfer(TokensReceipt<T> {
            id: object::new(ctx),
            sender: ctx.sender(),
            tokens: tokens.into_balance(),
        }, anchor.id.to_address());
    }
 
    /// client calls this to have the `Anchor` receive tokens and add them to its assets
    public fun receive_tokens<T>(anchor: &mut Anchor, tokens_receipt: transfer::Receiving<TokensReceipt<T>>, _ctx: &mut TxContext) {
        let TokensReceipt<T> { id, sender, tokens } = transfer::public_receive(&mut anchor.id, tokens_receipt);
        let token_type_name = type_name::get<T>().into_string();
        anchor.mutations.push_back(Mutation{
            sender: sender,
            asset_type: token_type_name.into_bytes(),
            amount: tokens.value(),
        });
        anchor.mutation_id = anchor.mutation_id + 1;
        anchor.assets.add_tokens(tokens);
        object::delete(id);
    }
 
    /// client calls this to send an NFT to the `Anchor`
    public fun send_nft(anchor: &mut Anchor, nft: Nft, ctx: &mut TxContext) {
       transfer::public_transfer(NftReceipt {
            id: object::new(ctx),
            sender: ctx.sender(),
            nft: nft,
        }, anchor.id.to_address());
     }
 
    /// client calls this to have the `Anchor` receive an NFT and add it to its assets
    public fun receive_nft(anchor: &mut Anchor, nft_received: transfer::Receiving<NftReceipt>, _ctx: &mut TxContext) {
        let NftReceipt { id, sender, nft } = transfer::public_receive(&mut anchor.id, nft_received);
        let nft_id = object::borrow_id<Nft>(&nft);
        anchor.mutations.push_back(Mutation{
            sender: sender,
            asset_type: nft_id.id_to_bytes(),
            amount: 0,
        });
        anchor.mutation_id = anchor.mutation_id + 1;
        anchor.assets.add_nft(nft);
        object::delete(id);
   }

    /// client calls this to send a `Request` to the `Anchor`
    public fun send_request(anchor: &Anchor, req: Request, _ctx: &mut TxContext) {
        // we could implement specific checks here for the `Request`

        // send the request object to the `Anchor`'s id
        transfer::public_transfer(req, anchor.id.to_address());
    }

    /// client calls this to have the `Anchor` receive a request and add it to the `Request` pool
    public fun receive_request(anchor: &mut Anchor, request: transfer::Receiving<Request>, _ctx: &mut TxContext){
        let req = transfer::public_receive<Request>(&mut anchor.id, request);
        let data = req.as_data();
        anchor.requests.push_back(move data);
        anchor.request_id = anchor.request_id + 1
    }

    ////////////////////// Governor-specific functions ////////////////////////

    /// local function makes sure only the governer of the `Anchor` can call a function
    fun check_governor(anchor: &Anchor, governor: &GovernorCap) {
        assert!(governor.anchor_id == anchor.id.uid_to_inner(), EWrongCaller);
    }

    /// `Anchor` governor calls this get up to `limit` mutations from the asset mutation pool
    /// The function returns the id of the first mutation and an array slice of up to `limit` mutations
    public fun get_mutations(anchor: &mut Anchor, governor: &GovernorCap, limit: u64, _ctx: &mut TxContext): (u64, vector<Mutation>) {
        check_governor(anchor, governor);

        if (anchor.mutations.length() <= limit) {
            return (anchor.mutation_id - anchor.mutations.length(), anchor.mutations)
        };
        let mut slice: vector<Mutation> = vector[];
        let mut i = 0;
        while (i < limit) {
            slice.push_back(anchor.mutations[i]);
            i = i + 1
        };

        (anchor.mutation_id - anchor.mutations.length(), slice)
    }

    /// `Anchor` governor calls this to communicate it has processed all mutations in the asset mutation pool before id `next`
    /// The function will remove all mutations with a lower id from the mutation pool
    /// This handshake function is used to prevent mutations getting lost unless definitely processed by the governor
    public fun processed_mutations(anchor: &mut Anchor, governor: &GovernorCap, next: u64, _ctx: &mut TxContext) {
        check_governor(anchor, governor);
        assert!(next <= anchor.mutation_id, EInvalidMutationID);
        let mut base = anchor.mutation_id - anchor.mutations.length();
        while (base < next) {
            anchor.mutations.remove(0);
            base = base + 1
        }
    }

    /// `Anchor` governor calls this get up to `limit` requests from the request pool
    /// The function returns the id of the first request and an array slice of up to `limit` requests
    public fun get_requests(anchor: &mut Anchor, governor: &GovernorCap, limit: u64, _ctx: &mut TxContext): (u64, vector<RequestData>) {
        check_governor(anchor, governor);

        if (anchor.requests.length() <= limit) {
            return (anchor.request_id - anchor.requests.length(), anchor.requests)
        };
        let mut slice: vector<RequestData> = vector[];
        let mut i = 0;
        while (i < limit) {
            slice.push_back(anchor.requests[i]);
            i = i + 1
        };

        (anchor.request_id - anchor.requests.length(), slice)
    }

    /// `Anchor` governor calls this to communicate it has processed all requests in the request pool before id `next`
    /// The function will remove all requests with a lower id from the request pool
    /// This handshake function is used to prevent requests getting lost unless definitely processed by the governor
    public fun processed_requests(anchor: &mut Anchor, governor: &GovernorCap, next: u64, _ctx: &mut TxContext) {
        check_governor(anchor, governor);
        assert!(next <= anchor.request_id, EInvalidRequestID);
        let mut base = anchor.request_id - anchor.requests.length();
        while (base < next) {
            anchor.requests.remove(0);
            base = base + 1
        }
    }

    /// `Anchor` governor calls this to transfer tokens from its `Assets` to an address
    public fun transfer_tokens<T>(anchor: &mut Anchor, governor: &GovernorCap, to: address, amount: u64, ctx: &mut TxContext){
        check_governor(anchor, governor);
        let tokens = anchor.assets.take_tokens<T>(amount);
        transfer::public_transfer(coin::from_balance(tokens, ctx), to);
    }

    /// `Anchor` governor calls this to transfer an NFT from its `Assets` to an address
    public fun transfer_nft(anchor: &mut Anchor, governor: &GovernorCap, to: address, nft_id: ID, _ctx: &mut TxContext){
        check_governor(anchor, governor);
        let nft = anchor.assets.take_nft(nft_id);
        transfer::public_transfer(nft, to);
    }
}
