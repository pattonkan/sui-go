package sui

import "github.com/howjmay/go-sui-sdk/sui/conn"

type ImplSuiAPI struct {
	http *conn.HttpClient
}

func NewSuiClient(http *conn.HttpClient) *ImplSuiAPI {
	return &ImplSuiAPI{
		http: http,
	}

}
