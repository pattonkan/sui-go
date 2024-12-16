# sui-go
Sui Golang SDK

[![Documentation (main)](https://img.shields.io/badge/docs-master-59f)](https://github.com/pattonkan/sui-go)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://github.com/pattonkan/sui-go/blob/main/LICENSE)

## Install

```sh
go get github.com/pattonkan/sui-go
```

## Features

* Full coverage of Sui's JSON RPC APIs for both HTTP and Websocket
* Websocket implementation for chain event/transaction subscriber
* Native Support of Sui's Programmable Transaction Blocks (PTB) by Programmable Transaction Builder (see see package [`suiptb`](sui/suiptb))
* Decoder for easy decoding of returned Move objects in BCS format (see package [`movebcs`](sui/movebcs))

## Usage

### Signer

Singer is a struct which holds the keypair of a user and will be used to sign transactions.

```go
import "github.com/pattonkan/sui-go/suisigner"

// Create a suisigner.Signer with mnemonic
mnemonic := "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
signer1, err := suisigner.NewSignerWithMnemonic(mnemonic, suisigner.KeySchemeFlagEd25519)
fmt.Printf("address   : %v\n", signer1.Address)

// create suisigner.Signer with seed
seed, err := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
signer2 := suisigner.NewSigner(seed, suisigner.KeySchemeFlagDefault)

// Get private key, public key, address
fmt.Printf("privateKey: %x\n", signer2.PrivateKey()[:32])
fmt.Printf("publicKey : %x\n", signer2.PublicKey())
fmt.Printf("address   : %v\n", signer2.Address)
```

### JSON RPC Client

All data interactions on the Sui chain are implemented through the JSON RPC client.

```go
import "github.com/pattonkan/sui-go/sui"
import "github.com/pattonkan/sui-go/suiclient"

client := suiclient.NewClient(rpcUrl) // some hardcoded endpoints are provided e.g. conn.TestnetEndpointUrl

// Call JSON RPC (e.g. call sui_getTransactionBlock)
digest, err := sui.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
require.NoError(t, err)
resp, err := client.GetTransactionBlock(
    context.Background(), *digest, sui.SuiTransactionBlockResponseOptions{
        ShowInput:          true,
        ShowEffects:        true,
        ShowObjectChanges:  true,
        ShowBalanceChanges: true,
        ShowEvents:         true,
    },
)
fmt.Println("transaction status = ", resp.Effects.Status)
fmt.Println("transaction timestamp = ", resp.TimestampMs)
```

### Programmable Transaction Blocks (PTB)

See `TestPTBMoveCall()` in [`programmable_transaction_builder_test.go`](sui/suiptb/programmable_transaction_builder_test.go)

### Decode move object in BCS

```go
import "github.com/pattonkan/sui-go/sui"
import "github.com/pattonkan/sui-go/suiclient"
import "github.com/pattonkan/sui-go/suisigner"

// get Coin object in BCS format by its ObjectRef
resGetObject, err := client.GetObject(context.TODO(), &suiclient.GetObjectRequest{
	ObjectId: targetCoinRef.ObjectId,
	Options: &suiclient.SuiObjectDataOptions{
		ShowBcs: true,
	},
})
require.NoError(t, err)
var moveCoin movebcs.MoveCoin
// here we get the Coin object for following usage
_, err = bcs.Unmarshal(resGetObject.GetMoveObjectInBcs(), &moveCoin)
require.NoError(t, err)
```

## Reference

* [Programmable Transaction Blocks (official doc)](https://docs.sui.io/concepts/transactions/prog-txn-blocks)
* [Sui Programmable Transaction Blocks Basics (TypeScript SDK)](https://sdk.mystenlabs.com/typescript/transaction-building/basics)
* [Sui JSON RPC API Reference](https://docs.sui.io/sui-api-ref)

## Acknowledgments

* https://github.com/coming-chat/go-sui-sdk
* https://github.com/block-vision/sui-go-sdk