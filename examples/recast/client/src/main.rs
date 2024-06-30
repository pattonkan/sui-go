mod faucet;

use shared_crypto::intent::{Intent, IntentMessage};
use sui_keys::keystore::{AccountKeystore, InMemKeystore, Keystore};
use sui_sdk::rpc_types::SuiTransactionBlockEffectsAPI;
use sui_sdk::{
    rpc_types::{
        SuiExecutionStatus, SuiObjectDataOptions, SuiObjectResponseQuery,
        SuiTransactionBlockResponseOptions,
    },
    types::Identifier,
    types::{
        programmable_transaction_builder::ProgrammableTransactionBuilder,
        quorum_driver_types::ExecuteTransactionRequestType,
        transaction::{Command, Transaction, TransactionData},
    },
    SuiClientBuilder,
};
use sui_types::{crypto::SignatureScheme, transaction::TransactionDataAPI};

const MNEMONIC_PHRASE: &str = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom";
const PACKAGE_ID_CONST: &str = "0x1234";

#[tokio::main]
async fn main() -> Result<(), anyhow::Error> {
    let mut keystore = Keystore::from(InMemKeystore::new_insecure_for_tests(0));
    let gen_addr = keystore
        .import_from_mnemonic(MNEMONIC_PHRASE, SignatureScheme::ED25519, None, None)
        .unwrap();

    // Sui local network
    let sui_client = SuiClientBuilder::default()
        .build("http://localhost:9000")
        .await?;

    faucet::request_fund_from_faucet(&gen_addr.to_string(), "http://localhost:9123/gas")
        .await
        .unwrap();

    let package_id = PACKAGE_ID_CONST.parse()?;

    // Get all my own objects
    let coins_response = &sui_client
        .read_api()
        .get_owned_objects(
            gen_addr,
            Some(SuiObjectResponseQuery::new_with_options(
                SuiObjectDataOptions::new().with_type(),
            )),
            None,
            None,
        )
        .await?;

    // Find a coin to use
    let coin = coins_response
        .data
        .iter()
        .find(|obj| obj.data.as_ref().unwrap().is_gas_coin())
        .unwrap();
    let coin = coin.data.as_ref().unwrap();

    let mut ptb = ProgrammableTransactionBuilder::new();

    let arg_u32 = ptb.command(Command::move_call(
        package_id,
        Identifier::new("recast")?,
        Identifier::new("create_container")?,
        vec![],
        vec![],
    ));

    let arg_u64 = ptb.pure(42949672970 as u64).expect("pure");
    ptb.command(Command::move_call(
        package_id,
        Identifier::new("recast")?,
        Identifier::new("try_recast")?,
        vec![],
        vec![arg_u64, arg_u32],
    ));

    let builder = ptb.finish();

    let gas_budget = 10_000_000;
    let gas_price = sui_client.read_api().get_reference_gas_price().await?;

    let tx_data = TransactionData::new_programmable(
        gen_addr,
        vec![coin.object_ref()],
        builder,
        gas_budget,
        gas_price,
    );

    let tx_data_kind = tx_data.kind();
    // Sign the transaction
    // let intent_msg = IntentMessage::new(Intent::sui_transaction(), &tx_data);
    // let raw_tx = bcs::to_bytes(&intent_msg).expect("bcs should not fail");
    // let mut hasher = sui_types::crypto::DefaultHash::default();
    // hasher.update(raw_tx.clone());
    // let digest = hasher.finalize().digest;

    // use SuiKeyPair to sign the digest.
    // let signature = keypair.sign(&digest);
    let dev_inspect_res = sui_client
        .read_api()
        .dev_inspect_transaction_block(gen_addr, tx_data_kind.clone(), None, None, None)
        .await
        .expect("");
    println!("dev_inspect_res.effects: {:#?}", dev_inspect_res.effects);
    // Submit the transaction
    // let transaction_response = sui_client
    //     .quorum_driver_api()
    //     .execute_transaction_block(
    //         Transaction::from_data(tx_data, vec![signature]),
    //         SuiTransactionBlockResponseOptions::full_content(),
    //         Some(ExecuteTransactionRequestType::WaitForLocalExecution),
    //     )
    //     .await?;
    // print!("done\nTransaction information: ");
    // println!("{:?}", transaction_response.digest);
    // if !transaction_response.status_ok().unwrap() {
    //     let a = transaction_response.effects.unwrap().status().clone();
    //     match a {
    //         SuiExecutionStatus::Success => {}
    //         SuiExecutionStatus::Failure { error } => {
    //             println!("error: {}", error)
    //         }
    //     }
    // }

    Ok(())
}
