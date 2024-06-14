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
func (i *ImplSuiAPI) WithSignerAndFund(seed []byte, index int) (*ImplSuiAPI, *sui_signer.Signer) {
	keySchemeFlag := sui_signer.KeySchemeFlagEd25519
	// special case if localnet is used, then
	if i.http.Url() == conn.LocalnetEndpointUrl {
		keySchemeFlag = sui_signer.KeySchemeFlagIotaEd25519
	}
	// there are only 256 different signers can be generated
	seed[0] = seed[0] + byte(index)
	signer := sui_signer.NewSigner(seed, keySchemeFlag)
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
	DefaultGasPrice  uint64 = 1000
)
