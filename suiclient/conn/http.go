package conn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ErrNoResult = errors.New("no result in JSON-RPC response")
)

// BatchElem is an element in a batch request.
type BatchElem struct {
	Method string
	Args   []interface{}
	// The result is unmarshaled into this field. Result must be set to a
	// non-nil pointer value of the desired type, otherwise the response will be
	// discarded.
	Result interface{}
	// Error is set if the server returns an error for this request, or if
	// unmarshaling into Result fails. It is not set for I/O errors.
	Error error
}

type HttpClient struct {
	idCounter uint32

	url    string
	client *http.Client
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url: strings.TrimRight(url, "/"),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    3,
				IdleConnTimeout: 30 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
	}
}

// Call performs a JSON-RPC call with the given arguments and unmarshals into
// the result if no error occurred.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (c *HttpClient) Call(result interface{}, method JsonRpcMethod, args ...interface{}) error {
	ctx := context.Background()
	return c.CallContext(ctx, result, method, args...)
}

// CallContext performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (c *HttpClient) CallContext(ctx context.Context, result interface{}, method JsonRpcMethod, args ...interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}
	msg, err := c.newMessage(method.String(), args...)
	if err != nil {
		return fmt.Errorf("failed to call newMessage: %w", err)
	}
	resp, err := c.doRequest(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to call doRequest:%w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}
	var respmsg jsonrpcMessage
	err = json.Unmarshal(resBody, &respmsg)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}
	if respmsg.Error != nil {
		return fmt.Errorf("sui returned error: %s", respmsg.Error)
	}
	if len(respmsg.Result) == 0 {
		return ErrNoResult
	}
	return json.Unmarshal(respmsg.Result, result)
}

// BatchCall sends all given requests as a single batch and waits for the server
// to return a response for all of them.
func (c *HttpClient) BatchCall(b []BatchElem) error {
	return c.BatchCallContext(context.Background(), b)
}

// BatchCallContext sends all given requests as a single batch and waits for the server
// to return a response for all of them. The wait duration is bounded by the
// context's deadline.
func (c *HttpClient) BatchCallContext(ctx context.Context, b []BatchElem) error {
	var (
		msgs = make([]*jsonrpcMessage, len(b))
		byId = make(map[string]int, len(b))
	)
	for i, elem := range b {
		msg, err := c.newMessage(elem.Method, elem.Args...)
		if err != nil {
			return err
		}
		msgs[i] = msg
		byId[string(msg.Id)] = i
	}
	resp, err := c.doRequest(ctx, msgs)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respmsgs []jsonrpcMessage
	if err := json.NewDecoder(resp.Body).Decode(&respmsgs); err != nil {
		return err
	}
	for idx, resp := range respmsgs {
		elem := &b[idx]
		if resp.Error != nil {
			elem.Error = resp.Error
			continue
		}
		if len(resp.Result) == 0 {
			elem.Error = ErrNoResult
			continue
		}
		elem.Error = json.Unmarshal(resp.Result, elem.Result)
	}
	return nil
}

func (c *HttpClient) Url() string {
	return c.url
}

func (c *HttpClient) nextId() json.RawMessage {
	id := atomic.AddUint32(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (c *HttpClient) newMessage(method string, paramsIn ...interface{}) (*jsonrpcMessage, error) {
	msg := &jsonrpcMessage{Version: version, Id: c.nextId(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (c *HttpClient) doRequest(ctx context.Context, msg interface{}) (*http.Response, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	req.Header.Set("Content-Type", "application/json")

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var body []byte
		if _, err := buf.ReadFrom(resp.Body); err == nil {
			body = buf.Bytes()
		}

		return nil, HTTPError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}
	return resp, nil
}
