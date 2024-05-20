package sui

import (
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
)

type ImplSuiAPI struct {
	http      *conn.HttpClient
	websocket *conn.WebsocketClient
}

func NewSuiClient(url string) *ImplSuiAPI {
	return &ImplSuiAPI{
		http: conn.NewHttpClient(url),
	}
}

// test only. If localnet is used then iota network will be connect
func NewTestSuiClientWithSignerAndFund(url string, mnemonic string) (*ImplSuiAPI, *sui_signer.Signer) {
	client := &ImplSuiAPI{
		http: conn.NewHttpClient(url),
	}

	keySchemeFlag := sui_signer.KeySchemeFlagEd25519
	// special case if localnet is used, then
	if client.http.Url() == conn.LocalnetEndpointUrl {
		keySchemeFlag = sui_signer.KeySchemeFlagIotaEd25519
	}
	signer, err := sui_signer.NewSignerWithMnemonic(mnemonic, keySchemeFlag)
	if err != nil {
		panic(err)
	}
	var faucetUrl string
	switch client.http.Url() {
	case conn.TestnetEndpointUrl:
		faucetUrl = conn.TestnetFaucetUrl
	case conn.DevnetEndpointUrl:
		faucetUrl = conn.DevnetFaucetUrl
	case conn.LocalnetEndpointUrl:
		faucetUrl = conn.LocalnetFaucetUrl
	default:
		panic("not supported network")
	}
	err = RequestFundFromFaucet(signer.Address, faucetUrl)
	if err != nil {
		panic(err)
	}
	return client, signer
}

func (i *ImplSuiAPI) WithWebsocket(url string) {
	i.websocket = conn.NewWebsocketClient(url)
}

func NewSuiWebsocketClient(url string) *ImplSuiAPI {
	return &ImplSuiAPI{
		websocket: conn.NewWebsocketClient(url),
	}
}

const (
	DefaultGasBudget uint64 = 10000000
)
