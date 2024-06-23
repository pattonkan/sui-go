import { Ed25519Keypair } from "@mysten/sui/keypairs/ed25519";
import { getFaucetHost, requestSuiFromFaucetV1 } from "@mysten/sui/faucet";
import { Transaction } from "@mysten/sui/transactions";
import { getFullnodeUrl, SuiClient } from "@mysten/sui/client";

const SWAP_PACKAGE_ID = "0x123";
const TESTCOIN_OBJ_ID = "0xtestcoin";

async function main() {
    let secretKey = new Uint8Array([
        50, 230, 119, 9, 86, 155, 106, 30, 245, 81, 234, 122, 116, 90, 172, 148, 59, 33, 88, 252, 134, 42, 231, 198,
        208, 141, 209, 116, 78, 21, 216, 24,
    ]);

    const keypair = Ed25519Keypair.fromSecretKey(secretKey);

    console.log("keypair.getSecretKey(): ", keypair.getSecretKey());
    console.log("keypair.toSuiAddress: ", keypair.toSuiAddress());

    requestSuiFromFaucetV1({ recipient: keypair.toSuiAddress(), host: getFaucetHost('testnet') });

    const rpcUrl = getFullnodeUrl("testnet");

    const client = new SuiClient({ url: rpcUrl });

    const testcoin_obj_type = TESTCOIN_OBJ_ID+"::testcoin::TESTCOIN";

    let get_coin_sui_res = await client.getCoins({owner: keypair.toSuiAddress()});
    let get_coin_testcoin_res = await client.getCoins({owner: keypair.toSuiAddress(), coinType: testcoin_obj_type});

    const tx = new Transaction();
    let arg = tx.moveCall({
        target: SWAP_PACKAGE_ID+"::swap::swap_sui",
        typeArguments: [testcoin_obj_type],
        arguments: [
            tx.object(get_coin_sui_res.data[1].coinObjectId),
            tx.object(get_coin_testcoin_res.data[0].coinObjectId),
            tx.pure.u64(3),
        ],
    });

    tx.transferObjects([arg], keypair.toSuiAddress());
    tx.setGasBudget(1000000);
    let res = await client.signAndExecuteTransaction({ signer: keypair, transaction: tx, options: {showEffects: true} });
    if (res.effects?.status.error != undefined) {
        console.log("fail to execute tx, "+ res.effects?.status.error);
    } else {
        console.log("successfully execute tx");
    }
}

main().catch(console.error);
