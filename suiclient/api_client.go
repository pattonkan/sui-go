package suiclient

import (
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
)

type ClientImpl struct {
	http      *conn.HttpClient
	websocket *conn.WebsocketClient
}

func NewClient(url string) *ClientImpl {
	return &ClientImpl{
		http: conn.NewHttpClient(url),
	}
}

// test only. If localnet is used then iota network will be connect
func (i *ClientImpl) WithSignerAndFund(seed []byte, index int) (*ClientImpl, *suisigner.Signer) {
	keySchemeFlag := suisigner.KeySchemeFlagEd25519
	// special case if localnet is used, then
	if i.http.Url() == conn.LocalnetEndpointUrl {
		keySchemeFlag = suisigner.KeySchemeFlagIotaEd25519
	}
	// there are only 256 different signers can be generated
	signer := suisigner.NewSignerByIndex(seed, keySchemeFlag, index)
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
	err := RequestFundFromFaucet(signer.Address, faucetUrl)
	if err != nil {
		panic(err)
	}
	return i, signer
}

func (i *ClientImpl) WithWebsocket(url string) {
	i.websocket = conn.NewWebsocketClient(url)
}

func NewSuiWebsocketClient(url string) *ClientImpl {
	return &ClientImpl{
		websocket: conn.NewWebsocketClient(url),
	}
}

const (
	DefaultGasBudget uint64 = 10000000
	DefaultGasPrice  uint64 = 1000
)
