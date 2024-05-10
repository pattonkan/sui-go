package sui

import "github.com/howjmay/sui-go/sui/conn"

type ImplSuiAPI struct {
	http      *conn.HttpClient
	websocket *conn.WebsocketClient
}

func NewSuiClient(url string) *ImplSuiAPI {
	return &ImplSuiAPI{
		http: conn.NewHttpClient(url),
	}
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
