mod faucet;

use std::str::FromStr;
use sui_keys::keystore::{AccountKeystore, InMemKeystore, Keystore};
use sui_sdk::SuiClientBuilder;
use sui_types::{
    base_types::{ObjectDigest, ObjectID, ObjectRef, SequenceNumber},
    crypto::SignatureScheme,
    programmable_transaction_builder::ProgrammableTransactionBuilder,
    transaction::{ObjectArg, Transaction, TransactionData, TransactionDataAPI},
    Identifier, TypeTag,
};

const MNEMONIC_PHRASE: &str = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom";
const SWAP_PACKAGE_ID: &str = "0x123";
const TESTCOIN_OBJ_ID: &str = "0xtestcoin";

#[tokio::main]
async fn main() {
    let mut keystore = Keystore::from(InMemKeystore::new_insecure_for_tests(0));
    let gen_addr = keystore
        .import_from_mnemonic(MNEMONIC_PHRASE, SignatureScheme::ED25519, None, None)
        .unwrap();

    println!("generated address: {}", gen_addr);

    let sui_client = SuiClientBuilder::default().build_testnet().await.unwrap();

    let testcoin_obj_type = format!("{}::testcoin::TESTCOIN", TESTCOIN_OBJ_ID);

    faucet::request_fund_from_faucet(
        &gen_addr.to_string(),
        "https://faucet.testnet.sui.io/v1/gas",
    )
    .await
    .unwrap();

    // take the 1st object as gas object
    let get_coins_sui_res = sui_client
        .coin_read_api()
        .get_coins(gen_addr.clone(), None, None, None)
        .await
        .unwrap();

    let get_coins_testcoin_res = sui_client
        .coin_read_api()
        .get_coins(
            gen_addr.clone(),
            Some(testcoin_obj_type.clone()),
            None,
            None,
        )
        .await
        .unwrap();

    let mut builder = ProgrammableTransactionBuilder::new();

    let arg0 = builder
        .obj(ObjectArg::ImmOrOwnedObject(
            get_coins_testcoin_res.data[0].object_ref(),
        ))
        .unwrap();
    let arg1 = builder
        .obj(ObjectArg::ImmOrOwnedObject(
            get_coins_sui_res.data[1].object_ref(),
        ))
        .unwrap();
    let arg2 = builder.pure(3 as u64).unwrap();

    builder.programmable_move_call(
        ObjectID::from_hex_literal(SWAP_PACKAGE_ID).unwrap(),
        Identifier::from_str("swap").unwrap(),
        Identifier::from_str("create_pool").unwrap(),
        vec![TypeTag::from_str(&testcoin_obj_type).unwrap()],
        vec![arg0, arg1, arg2],
    );
    let pt = builder.finish();

    let tx_data = TransactionData::new_programmable(
        gen_addr,
        vec![get_coins_sui_res.data[0].object_ref()],
        pt,
        10000000,
        1000,
    );

    let transaction_response = sui_client
        .read_api()
        .dev_inspect_transaction_block(gen_addr, tx_data.into_kind(), None, None, None)
        .await
        .unwrap();

    println!("{:#?}", transaction_response);
}
