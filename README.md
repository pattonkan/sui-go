# sui-go
Sui Golang SDK

[![Documentation (main)](https://img.shields.io/badge/docs-master-59f)](https://github.com/howjmay/sui-go)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://github.com/howjmay/sui-go/blob/main/LICENSE)

The Sui Golang SDK. We welcome other developers to participate in the development and testing of sui-sdk.

## Install

```sh
go get github.com/howjmay/sui-go
```

## Usage

### Signer

Singer is a struct which holds the keypair of a user and will be used to sign transactions.

```go
import "github.com/howjmay/sui-go/sui_signer"

// Create a sui_signer.Signer with mnemonic
mnemonic := "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
signer1, err := sui_signer.NewSignerWithMnemonic(mnemonic, sui_signer.KeySchemeFlagEd25519)
fmt.Printf("address   : %v\n", signer1.Address)

// create sui_signer.Signer with seed
seed, err := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
signer2 := sui_signer.NewSigner(seed, sui_signer.KeySchemeFlagDefault)

// Get private key, public key, address
fmt.Printf("privateKey: %x\n", signer2.PrivateKey()[:32])
fmt.Printf("publicKey : %x\n", signer2.PublicKey())
fmt.Printf("address   : %v\n", signer2.Address)
```

### JSON RPC Client

All data interactions on the Sui chain are implemented through the JSON RPC client.

```go
import "github.com/howjmay/sui-go/sui"
import "github.com/howjmay/sui-go/sui_types"

client := sui.NewSuiClient(rpcUrl) // some hardcoded endpoints are provided e.g. conn.TestnetEndpointUrl

// Call JSON RPC (e.g. call sui_getTransactionBlock)
digest, err := sui_types.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
require.NoError(t, err)
resp, err := client.GetTransactionBlock(
    context.Background(), *digest, models.SuiTransactionBlockResponseOptions{
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

## Reference

* [Programmable Transaction Blocks (official doc)](https://docs.sui.io/concepts/transactions/prog-txn-blocks)
* [Sui Programmable Transaction Blocks Basics (TypeScript SDK)](https://sdk.mystenlabs.com/typescript/transaction-building/basics)
* [Sui JSON RPC API Reference](https://docs.sui.io/sui-api-ref)
