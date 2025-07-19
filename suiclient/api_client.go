package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/suisigner/suicrypto"
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
func (i *ClientImpl) WithSignerAndFund(seed []byte, flag suicrypto.KeySchemeFlag, index int) (*ClientImpl, *suisigner.Signer) {
	// there are only 256 different signers can be generated
	signer := suisigner.NewSignerByIndex(seed, flag, index)
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

func (i *ClientImpl) WithWebsocket(ctx context.Context, url string) {
	i.websocket = conn.NewWebsocketClient(ctx, url)
}

func NewSuiWebsocketClient(ctx context.Context, url string) *ClientImpl {
	return &ClientImpl{
		websocket: conn.NewWebsocketClient(ctx, url),
	}
}

func (i *ClientImpl) Close() error {
	if i.websocket != nil {
		return i.websocket.Close()
	}

	return nil
}

const (
	DefaultGasBudget uint64 = 10000000
	DefaultGasPrice  uint64 = 1000
)
