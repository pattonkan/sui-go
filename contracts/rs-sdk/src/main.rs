use bip39::{Language, Mnemonic, Seed};
use serde_json;
use shared_crypto::intent::Intent;
use std::env;
use std::path::{Path, PathBuf};
use sui_keys::keystore::{AccountKeystore, FileBasedKeystore, InMemKeystore};
use sui_sdk::json::SuiJsonValue;
use sui_sdk::rpc_types::{SuiTransactionBlockEffects, SuiTransactionBlockEffectsAPI};
use sui_sdk::types::crypto;
use sui_sdk::{
    rpc_types::SuiTransactionBlockResponseOptions,
    types::base_types::ObjectID,
    types::{quorum_driver_types::ExecuteTransactionRequestType, transaction::Transaction},
    SuiClientBuilder,
};

const PHRASE:&str = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom";
const PACKAGE_ID: &str = "0x44cc91f9d6e807dcd4ea40731b1455146b56780ea77ce752bd226e756ff5ff51";

#[tokio::main]
async fn main() -> Result<(), anyhow::Error> {
    let mut keystore = sui_keys::keystore::Keystore::from(InMemKeystore::new_insecure_for_tests(0));
    let sui_address =
        keystore.import_from_mnemonic(&PHRASE, crypto::SignatureScheme::ED25519, None)?;
    // Sui local network
    let sui_client = SuiClientBuilder::default()
        .build("http://127.0.0.1:9000")
        .await?;
    println!("Sui local version: {}", sui_client.api_version());

    let gas_budget = 5_000_000;
    let gas_price = sui_client
        .read_api()
        .get_reference_gas_price()
        .await
        .unwrap();

    let tx_data = sui_client
        .transaction_builder()
        .move_call(
            sui_address,
            ObjectID::from_hex_literal(PACKAGE_ID).unwrap(),
            "sdk_tests",
            "input_byte_array_of_arrays",
            vec![],
            vec![SuiJsonValue::new(
                serde_json::to_value(vec![vec![104, 97, 104, 97], vec![97, 98, 99]]).unwrap(),
            )
            .unwrap()],
            None,
            gas_budget,
            Some(gas_price),
        )
        .await
        .unwrap();

    // Sign the transaction
    let signature = keystore
        .sign_secure(&sui_address, &tx_data, Intent::sui_transaction())
        .unwrap();

    // // Submit the transaction
    let transaction_response = sui_client
        .quorum_driver_api()
        .execute_transaction_block(
            Transaction::from_data(tx_data, vec![signature]),
            SuiTransactionBlockResponseOptions::full_content(),
            Some(ExecuteTransactionRequestType::WaitForLocalExecution),
        )
        .await
        .unwrap();
    print!("done\nTransaction information: ");
    println!("{:?}", transaction_response.digest);

    match transaction_response.effects.unwrap() {
        SuiTransactionBlockEffects::V1(v) => {
            println!("v.status: {:?}", v.status);
        }
    }

    Ok(())
}
