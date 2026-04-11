# Contributing

## Development Setup

1. Install Go 1.23+
2. Clone repo and enter root
3. Run `make deps`
4. Run `make check`

Optional for integration tests:

- Install Sui CLI
- Start localnet faucet:
- `sui start --force-regenesis --with-faucet=9123`

## Repository Map

- `suiclient/`: RPC, GraphQL, websocket APIs and request/response types
- `sui/`: core Sui types, BCS-related helpers, PTB builder
- `suisigner/`: key derivation, signing, multisig support
- `contracts/`: Move example contracts and publish artifacts
- `examples/`: runnable integration examples in Go/Rust/TypeScript
- `utils/`: utility helpers used by tests and packages

## Test Strategy

- `make test-unit` for deterministic/offline test loop
- `make test-integration` when changing network/localnet behavior
- `make test-all` before major PRs

## Pull Requests

1. Keep PRs focused and small when possible
2. Include tests for behavior changes
3. Run `make fmt lint test-unit` before opening PR
4. If touching chain behavior, include integration evidence (`make test-integration` logs)

## AI-Assisted Development

- Read `AGENTS.md` first
- Prefer `make` targets over ad-hoc shell commands
- Do not assume localnet or public network availability unless explicitly set up
