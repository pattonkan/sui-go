
/// Module: isc
module isc::anchor {
    use isc::{
        assets::{Self, Assets},
        request::{Self, Request},
    };
    use stardust::{
        nft::{Nft},
    };
    use sui::{
        coin::{Self, Coin},
        event,
    };

    public struct Anchor has key, store {
        id: UID,
        /// assets controlled by the anchor
        assets: Assets,
    }

    /// `Anchor` owner uses the ID in this event to determine the type of the corresponding object:
    /// - type `Coint<T>`: update account of sender with the type/amount of tokens this coin
    ///   represents and call `receive_coin<T>()` to finalize adding the tokens to the Anchor assets.
    /// - type `Nft`: update account of sender with this NFT and call `receive_nft()` to finalize
    ///   adding the NFT to the Anchor assets.
    /// - type `Request`: extract the `RequestData` from the request for further processing and
    ///   call `receive_request()` to finalize receiving the request and subsequently destroy it.
    public struct AnchorEvent has copy, drop {
        id: ID,
        sender: address,
    }

    /// starts a new chain by creating a new `Anchor` for it
    public fun start_new_chain(ctx: &mut TxContext): Anchor {
        Anchor{
            id: object::new(ctx),
            assets: assets::new(ctx),
         }
    }
 
    /// Client calls this to send a `Coint<T>` to the `Anchor`'s address
    public fun send_coin<T>(anchor_address: address, coin: Coin<T>, ctx: &mut TxContext) {
        event::emit(AnchorEvent { id: object::id(&coin), sender: ctx.sender() });
        transfer::public_transfer(coin, anchor_address)
    }
 
    /// `Anchor` owner calls this to receive a `Coint<T>` when triggered by a `send_coin` event
    public fun receive_coin<T>(anchor: &mut Anchor, received_coin: transfer::Receiving<Coin<T>>, _ctx: &mut TxContext) {
        let coin = transfer::public_receive(&mut anchor.id, received_coin);
        anchor.assets.add_balance(coin.into_balance())
    }
 
    /// Client calls this to send an `Nft` to the `Anchor`'s address
    public fun send_nft(anchor_address: address, nft: Nft, ctx: &mut TxContext) {
        event::emit(AnchorEvent { id: object::id(&nft), sender: ctx.sender() });
        transfer::public_transfer(nft, anchor_address)
    }
 
    /// `Anchor` owner this to receive an `Nft` when triggered by a `send_nft` event
    public fun receive_nft(anchor: &mut Anchor, received_nft: transfer::Receiving<Nft>, _ctx: &mut TxContext) {
        let nft = transfer::public_receive(&mut anchor.id, received_nft);
        anchor.assets.add_nft(nft)
   }

    /// Client calls this to send a `Request` to the `Anchor`'s address
    public fun send_request(anchor_address: address, request: Request, ctx: &mut TxContext) {
        event::emit(AnchorEvent { id: object::id(&request), sender: ctx.sender() });
        transfer::public_transfer(request, anchor_address)
    }

    /// `Anchor` owner calls this to receive a `Request` when triggered by a `send_request` event
    public fun receive_request(anchor: &mut Anchor, received_request: transfer::Receiving<Request>, _ctx: &mut TxContext) {
        let req = transfer::public_receive(&mut anchor.id, received_request);
        request::destroy_request(req)
    }

    /// `Anchor` owner calls this to transfer a `Coin<T>` from its `Assets` to an L1 address
    public fun transfer_coin<T>(anchor: &mut Anchor, to: address, amount: u64, ctx: &mut TxContext){
        let coin = anchor.assets.take_coin<T>(amount);
        transfer::public_transfer(coin::from_balance(coin, ctx), to);
    }

    /// `Anchor` owner calls this to transfer an `Nft` from its `Assets` to an L1 address
    public fun transfer_nft(anchor: &mut Anchor, to: address, nft_id: ID, _ctx: &mut TxContext){
        let nft = anchor.assets.take_nft(nft_id);
        transfer::public_transfer(nft, to);
    }
}
