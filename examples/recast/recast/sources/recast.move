/// Module: recast
module recast::recast {
    use sui::event;

    public struct U32_VALS has copy, drop {
        f0: u32,
        f1: u32,
    }

    public struct EventHere has copy, drop {
        u64_val: u64,
        u32_vals: U32_VALS,        
    }

    public fun create_container(): U32_VALS {
        U32_VALS {
            f0: 10,
            f1: 10,
        }
    }

    public fun try_recast(u64_val: &u64, u32_vals: &U32_VALS) {
        event::emit(EventHere {
            u64_val: *u64_val,
            u32_vals: *u32_vals,            
        })
    }
}
