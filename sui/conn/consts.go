package conn

const (
	DevnetEndpointUrl   = "https://fullnode.devnet.sui.io"
	TestnetEndpointUrl  = "https://fullnode.testnet.sui.io"
	MainnetEndpointUrl  = "https://fullnode.mainnet.sui.io"
	LocalnetEndpointUrl = "http://localhost:9000"

	DevnetWebsocketEndpointUrl   = "wss://rpc.devnet.sui.io:443"
	TestnetWebsocketEndpointUrl  = "wss://rpc.testnet.sui.io:443"
	MainnetWebsocketEndpointUrl  = "wss://rpc.mainnet.sui.io:443"
	LocalnetWebsocketEndpointUrl = "ws://localhost:9000"

	DevnetFaucetUrl   = "https://faucet.devnet.sui.io/v1/gas"
	TestnetFaucetUrl  = "https://faucet.testnet.sui.io/v1/gas"
	LocalnetFaucetUrl = "http://localhost:9123/gas"
)

const (
	ChainIdentifierSuiDevnet  = "f4950f51"
	ChainIdentifierSuiTestnet = "4c78adac"
	ChainIdentifierSuiMainnet = "35834a8a"

	// localnet doesn't have a fixed ChainIdentifier
	// ChainIdentifierLocalnet
)
