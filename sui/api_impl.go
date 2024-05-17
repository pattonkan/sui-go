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

// test only func, which supports only testnet/devnet/localnet
func (i *ImplSuiAPI) WithSignerAndFund(mnemonic string) (*ImplSuiAPI, *sui_signer.Signer) {
	signer, err := sui_signer.NewSignerWithMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}
	var faucetUrl string
	switch i.http.Url() {
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
	return i, signer
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
