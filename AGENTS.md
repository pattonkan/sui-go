# AGENTS

This file helps AI coding agents and developers work safely and quickly in this repository.

## Scope

- Language: Go (`go 1.23` in `go.mod`)
- Main packages:
- `suiclient`: Sui RPC/GraphQL/WebSocket client APIs
- `sui`: base Sui types + programmable transaction builder in `sui/suiptb`
- `suisigner`: keypairs/signing/multisig
- `examples`: runnable samples for swap/recast/event pubsub

## Fast Start

- Install deps: `go mod download`
- Quick confidence checks: `make check`
- Unit-only test pass: `make test-unit`

## Test Model

The test suite has two classes:

1. Unit/offline tests (safe default, deterministic)
2. Integration/network tests (need local Sui node and/or public network)

Use these commands:

- `make test-unit`: excludes known network/localnet integration tests
- `make test-all`: runs all tests
- `make test-integration`: alias for full integration run

If integration tests fail with `localhost:9123` or DNS/network errors, this is expected without local services.

## Local Integration Environment

For localnet-based tests, start Sui with faucet:

- `sui start --force-regenesis --with-faucet=9123`

Then run:

- `make test-integration`

## Editing Rules

- Prefer minimal, focused changes.
- Avoid API-breaking renames unless required.
- Keep existing package boundaries (`sui`, `suiclient`, `suisigner`) intact.
- Add/update tests with behavioral changes.

## Common AI Tasks

- Add a new RPC method:
1. Add request/response types in `suiclient`
2. Add method wrapper in matching API file (read/write/coin/etc.)
3. Add unit tests if possible; integration test only if chain state is required

- Extend PTB builder:
1. Add types/arguments in `sui/suiptb`
2. Add serialization tests in `sui/suiptb/*_test.go`

- Add signer scheme behavior:
1. Update `suisigner/suicrypto`
2. Add fixtures/tests for sign/verify/address derivation

## PR Quality Bar

Before opening PR:

- `make fmt`
- `make lint`
- `make test-unit`

Run `make test-integration` if your change touches chain/network behavior.
