package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
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
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result"`
	ID      int64  `json:"id"`
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

	var rsp SubscriptionResp
	if gjson.ParseBytes(msgData).Get("error").Exists() {
		return fmt.Errorf(gjson.ParseBytes(msgData).Get("error").String())
	}

	err = json.Unmarshal([]byte(gjson.ParseBytes(msgData).String()), &rsp)
	if err != nil {
		return err
	}

	fmt.Printf("establish successfully, subscriptionID: %d, Waiting to accept data...\n", rsp.Result)

	go func(conn *websocket.Conn) {
		for {
			messageType, messageData, err := conn.ReadMessage()
			if nil != err {
				log.Println(err)
				break
			}
			switch messageType {
			case websocket.TextMessage:
				resultCh <- messageData

			default:
				continue
			}
		}
	}(c.conn)

	return nil
}

func (c *WebsocketClient) newMessage(method string, paramsIn ...interface{}) (*jsonrpcMessage, error) {
	msg := &jsonrpcMessage{Version: version, ID: c.nextID(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (c *WebsocketClient) nextID() json.RawMessage {
	id := atomic.AddUint32(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}
