package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	idCounter uint32
	url       string
	conn      *websocket.Conn
}

type CallOp struct {
	Method string
	Params []interface{}
}

type SubscriptionResp struct {
	Error   string `json:"error,omitempty"`
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result"`
	Id      int64  `json:"id"`
}

var DefaultReceiveMsgChanSize = 10

func NewWebsocketClient(url string) *WebsocketClient {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to websocket server: %s, %s", err, url))
	}

	return &WebsocketClient{
		url:  url,
		conn: conn,
	}
}

func (c *WebsocketClient) Call(resultCh chan []byte, method JsonRpcMethod, args ...interface{}) error {
	ctx := context.Background()
	return c.CallContext(ctx, resultCh, method, args...)
}

func (c *WebsocketClient) CallContext(ctx context.Context, resultCh chan []byte, method JsonRpcMethod, args ...interface{}) error {
	msg, err := c.newMessage(method.String(), args...)
	if err != nil {
		return err
	}
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, reqBody)
	if nil != err {
		return err
	}

	_, msgData, err := c.conn.ReadMessage()
	if nil != err {
		return err
	}
	var resp SubscriptionResp
	if err = json.Unmarshal(msgData, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return fmt.Errorf("websocket CallContext error: %s", resp.Error)
	}

	go func(conn *websocket.Conn) {
		for {
			messageType, messageData, err := conn.ReadMessage()
			if nil != err {
				log.Fatal(err)
				break
			}
			switch messageType {
			case websocket.TextMessage:
				var respmsg jsonrpcMessage
				if err := json.Unmarshal(messageData, &respmsg); err != nil {
					log.Fatalf("could not unmarshal response body: %s", err)
				}
				if respmsg.Error != nil {
					log.Fatalf("sui returned error: %s", respmsg.Error)
				}
				if len(respmsg.Params) == 0 {
					log.Fatal(ErrNoResult)
				}
				var prams jsonrpcWebsocketParams
				if err := json.Unmarshal(respmsg.Params, &prams); err != nil {
					log.Fatalf("could not unmarshal respmsg.Params: %s", err)
				}
				resultCh <- prams.Result

			default:
				continue
			}
		}
	}(c.conn)

	return nil
}

func (c *WebsocketClient) newMessage(method string, paramsIn ...interface{}) (*jsonrpcMessage, error) {
	msg := &jsonrpcMessage{Version: version, Id: c.nextId(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (c *WebsocketClient) nextId() json.RawMessage {
	id := atomic.AddUint32(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}
