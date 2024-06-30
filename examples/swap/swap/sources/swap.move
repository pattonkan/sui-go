/// Module: swap
module swap::swap {
    use sui::coin::{Self, Coin};
    use sui::balance::{Self, Supply, Balance};
    use sui::sui::SUI;
    use sui::math;

    /// For when supplied Coin is zero.
    const EZeroAmount: u64 = 0;

    /// For when pool fee is set incorrectly.
    /// Allowed values are: [0-10000).
    const EWrongFee: u64 = 1;

    /// For when someone tries to swap in an empty pool.
    const EReservesEmpty: u64 = 2;

    /// For when someone attempts to add more liquidity than u128 Math allows.
    const EPoolFull: u64 = 4;

    /// The integer scaling setting for fees calculation.
    const FEE_SCALING: u128 = 10000;

    /// The max value that can be held in one of the Balances of
    /// a Pool. U64 MAX / FEE_SCALING
    const MAX_POOL_VALUE: u64 = {
        18446744073709551615 / 10000
    };

    /// The Pool token that will be used to mark the pool share
    /// of a liquidity provider. The first type parameter stands
    /// for the witness type of a pool. The seconds is for the
    /// coin held in the pool.
    public struct LSP<phantom T> has drop {}

    /// The pool with exchange.
    ///
    /// - `fee_percent` should be in the range: [0-10000), meaning
    /// that 10000 is 100% and 1 is 0.1%
    public struct Pool<phantom T> has key {
        id: UID,
        sui: Balance<SUI>,
        token: Balance<T>,
        lsp_supply: Supply<LSP<T>>,
        /// Fee Percent is denominated in basis points.
        fee_percent: u64
    }

    #[allow(unused_function)]
    /// Module initializer is empty - to publish a new Pool one has
    /// to create a type which will mark LSPs.
    fun init(_: &mut TxContext) {}

    /// Create new `Pool` for token `T`. Each Pool holds a `Coin<T>`
    /// and a `Coin<SUI>`. Swaps are available in both directions.
    ///
    /// Share is calculated based on Uniswap's constant product formula:
    ///  liquidity = sqrt( X * Y )
    public fun create_pool<T>(
        token: Coin<T>,
        sui: Coin<SUI>,
        fee_percent: u64,
        ctx: &mut TxContext
    ): Coin<LSP<T>> {
        let sui_amt = sui.value();
        let tok_amt = token.value();

        assert!(sui_amt > 0 && tok_amt > 0, EZeroAmount);
        assert!(sui_amt < MAX_POOL_VALUE && tok_amt < MAX_POOL_VALUE, EPoolFull);
        assert!(fee_percent >= 0 && fee_percent < 10000, EWrongFee);

        // Initial share of LSP is the sqrt(a) * sqrt(b)
        let share = math::sqrt(sui_amt) * math::sqrt(tok_amt);
        let mut lsp_supply = balance::create_supply(LSP<T> {});
        let lsp = lsp_supply.increase_supply(share);

        transfer::share_object(Pool {
            id: object::new(ctx),
            token: token.into_balance(),
            sui: sui.into_balance(),
            lsp_supply,
            fee_percent
        });

        coin::from_balance(lsp, ctx)
    }


    /// Entrypoint for the `swap_sui` method. Sends swapped token
    /// to sender.
    entry fun swap_sui_<T>(
        pool: &mut Pool<T>, sui: Coin<SUI>, ctx: &mut TxContext
    ) {
        transfer::public_transfer(
            swap_sui(pool, sui, ctx),
            ctx.sender()
        )
    }

    /// Swap `Coin<SUI>` for the `Coin<T>`.
    /// Returns Coin<T>.
    public fun swap_sui<T>(
        pool: &mut Pool<T>, sui: Coin<SUI>, ctx: &mut TxContext
    ): Coin<T> {
        assert!(sui.value() > 0, EZeroAmount);

        let sui_balance = sui.into_balance();

        // Calculate the output amount - fee
        let (sui_reserve, token_reserve, _) = get_amounts(pool);

        assert!(sui_reserve > 0 && token_reserve > 0, EReservesEmpty);

        let output_amount = get_input_price(
            sui_balance.value(),
            sui_reserve,
            token_reserve,
            pool.fee_percent
        );

        pool.sui.join(sui_balance);
        coin::take(&mut pool.token, output_amount, ctx)
    }

    /// Entry point for the `swap_token` method. Sends swapped SUI
    /// to the sender.
    entry fun swap_token_<T>(
        pool: &mut Pool<T>, token: Coin<T>, ctx: &mut TxContext
    ) {
        transfer::public_transfer(
            swap_token(pool, token, ctx),
            ctx.sender()
        )
    }

    /// Swap `Coin<T>` for the `Coin<SUI>`.
    /// Returns the swapped `Coin<SUI>`.
    public fun swap_token<T>(
        pool: &mut Pool<T>, token: Coin<T>, ctx: &mut TxContext
    ): Coin<SUI> {
        assert!(token.value() > 0, EZeroAmount);

        let tok_balance = token.into_balance();
        let (sui_reserve, token_reserve, _) = get_amounts(pool);

        assert!(sui_reserve > 0 && token_reserve > 0, EReservesEmpty);

        let output_amount = get_input_price(
            tok_balance.value(),
            token_reserve,
            sui_reserve,
            pool.fee_percent
        );

        pool.token.join(tok_balance);
        coin::take(&mut pool.sui, output_amount, ctx)
    }

    /// Entrypoint for the `add_liquidity` method. Sends `Coin<LSP>` to
    /// the transaction sender.
    entry fun add_liquidity_<T>(
        pool: &mut Pool<T>, sui: Coin<SUI>, token: Coin<T>, ctx: &mut TxContext
    ) {
        transfer::public_transfer(
            add_liquidity(pool, sui, token, ctx),
            ctx.sender()
        );
    }

    /// Add liquidity to the `Pool`. Sender needs to provide both
    /// `Coin<SUI>` and `Coin<T>`, and in exchange he gets `Coin<LSP>` -
    /// liquidity provider tokens.
    public fun add_liquidity<T>(
        pool: &mut Pool<T>, sui: Coin<SUI>, token: Coin<T>, ctx: &mut TxContext
    ): Coin<LSP<T>> {
        assert!(sui.value() > 0, EZeroAmount);
        assert!(sui.value() > 0, EZeroAmount);

        let sui_balance = sui.into_balance();
        let tok_balance = token.into_balance();

        let (sui_amount, tok_amount, lsp_supply) = get_amounts(pool);

        let sui_added = sui_balance.value();
        let tok_added = tok_balance.value();
        let share_minted = math::min(
            (sui_added * lsp_supply) / sui_amount,
            (tok_added * lsp_supply) / tok_amount
        );

        let sui_amt = pool.sui.join(sui_balance);
        let tok_amt = pool.token.join(tok_balance);

        assert!(sui_amt < MAX_POOL_VALUE, EPoolFull);
        assert!(tok_amt < MAX_POOL_VALUE, EPoolFull);

        let balance = pool.lsp_supply.increase_supply(share_minted);
        coin::from_balance(balance, ctx)
    }

    /// Entrypoint for the `remove_liquidity` method. Transfers
    /// withdrawn assets to the sender.
    entry fun remove_liquidity_<T>(
        pool: &mut Pool<T>,
        lsp: Coin<LSP<T>>,
        ctx: &mut TxContext
    ) {
        let (sui, token) = remove_liquidity(pool, lsp, ctx);
        let sender = ctx.sender();

        transfer::public_transfer(sui, sender);
        transfer::public_transfer(token, sender);
    }

    /// Remove liquidity from the `Pool` by burning `Coin<LSP>`.
    /// Returns `Coin<T>` and `Coin<SUI>`.
    public fun remove_liquidity<T>(
        pool: &mut Pool<T>,
        lsp: Coin<LSP<T>>,
        ctx: &mut TxContext
    ): (Coin<SUI>, Coin<T>) {
        let lsp_amount = lsp.value();

        // If there's a non-empty LSwe can
        assert!(lsp_amount > 0, EZeroAmount);

        let (sui_amt, tok_amt, lsp_supply) = get_amounts(pool);
        let sui_removed = (sui_amt * lsp_amount) / lsp_supply;
        let tok_removed = (tok_amt * lsp_amount) / lsp_supply;

        pool.lsp_supply.decrease_supply(lsp.into_balance());

        (
            coin::take(&mut pool.sui, sui_removed, ctx),
            coin::take(&mut pool.token, tok_removed, ctx)
        )
    }

    /// Public getter for the price of SUI in token T.
    /// - How much SUI one will get if they send `to_sell` amount of T;
    public fun sui_price<T>(pool: &Pool<T>, to_sell: u64): u64 {
        let (sui_amt, tok_amt, _) = get_amounts(pool);
        get_input_price(to_sell, tok_amt, sui_amt, pool.fee_percent)
    }

    /// Public getter for the price of token T in SUI.
    /// - How much T one will get if they send `to_sell` amount of SUI;
    public fun token_price<T>(pool: &Pool<T>, to_sell: u64): u64 {
        let (sui_amt, tok_amt, _) = get_amounts(pool);
        get_input_price(to_sell, sui_amt, tok_amt, pool.fee_percent)
    }


    /// Get most used values in a handy way:
    /// - amount of SUI
    /// - amount of token
    /// - total supply of LSP
    public fun get_amounts<T>(pool: &Pool<T>): (u64, u64, u64) {
        (
            pool.sui.value(),
            pool.token.value(),
            pool.lsp_supply.supply_value()
        )
    }

    /// Calculate the output amount minus the fee - 0.3%
    public fun get_input_price(
        input_amount: u64, input_reserve: u64, output_reserve: u64, fee_percent: u64
    ): u64 {
        // up casts
        let (
            input_amount,
            input_reserve,
            output_reserve,
            fee_percent
        ) = (
            (input_amount as u128),
            (input_reserve as u128),
            (output_reserve as u128),
            (fee_percent as u128)
        );

        let input_amount_with_fee = input_amount * (FEE_SCALING - fee_percent);
        let numerator = input_amount_with_fee * output_reserve;
        let denominator = (input_reserve * FEE_SCALING) + input_amount_with_fee;

        (numerator / denominator as u64)
    }

    #[test_only]
    public fun init_for_testing(ctx: &mut TxContext) {
        init(ctx)
    }
}