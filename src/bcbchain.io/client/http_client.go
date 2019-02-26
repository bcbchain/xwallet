package rpcclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/pkg/errors"

	"bcbchain.io/types"
)

type HTTPClient interface {
	Call(method string, params map[string]interface{}, result interface{}) (interface{}, error)
}

func makeHTTPDialer(remoteAddr string) (string, func(string, string) (net.Conn, error)) {
	parts := strings.SplitN(remoteAddr, "://", 2)
	var protocol, address string
	if len(parts) == 1 {

		protocol, address = "tcp", remoteAddr
	} else if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	} else {

		msg := fmt.Sprintf("Invalid addr: %s", remoteAddr)
		return msg, func(_ string, _ string) (net.Conn, error) {
			return nil, errors.New(msg)
		}
	}

	if protocol == "http" {
		protocol = "tcp"
	}

	trimmedAddress := strings.Replace(address, "/", ".", -1)
	return trimmedAddress, func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, address)
	}
}

func makeHTTPClient(remoteAddr string) (string, *http.Client) {
	address, dialer := makeHTTPDialer(remoteAddr)
	return "http://" + address, &http.Client{
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
}

func makeHTTPSClient(remoteAddr string, pool *x509.CertPool, disableKeepAlive bool) (string, *http.Client) {

	tr := new(http.Transport)
	tr.DisableKeepAlives = disableKeepAlive
	tr.IdleConnTimeout = time.Second * 5
	if pool != nil {
		tr.TLSClientConfig = &tls.Config{RootCAs: pool}
	} else {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return remoteAddr, &http.Client{Transport: tr, Timeout: time.Duration(time.Second * 60)}
}

type JSONRPCClient struct {
	address	string
	client	*http.Client
}

func NewJSONRPCClient(remote string) *JSONRPCClient {
	address, client := makeHTTPClient(remote)
	return &JSONRPCClient{
		address:	address,
		client:		client,
	}
}

func NewJSONRPCClientEx(remote, certFile string, disableKeepAlive bool) *JSONRPCClient {
	var pool *x509.CertPool
	if certFile != "" {
		pool = x509.NewCertPool()
		caCert, err := ioutil.ReadFile(certFile)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		pool.AppendCertsFromPEM(caCert)
	}

	address, client := makeHTTPSClient(remote, pool, disableKeepAlive)

	return &JSONRPCClient{
		address:	address,
		client:		client,
	}
}

func (c *JSONRPCClient) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {

	request, err := types.MapToRequest("jsonrpc-client", method, params)
	if err != nil {

		return nil, err
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("lib client http_client error to json.Marshal(request)")

		return nil, err
	}

	requestBuf := bytes.NewBuffer(requestBytes)

	httpResponse, err := c.client.Post(c.address, "text/json", requestBuf)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	return unmarshalResponseBytes(responseBytes, result)
}

type URIClient struct {
	address	string
	client	*http.Client
}

func NewURIClient(remote string) *URIClient {
	address, client := makeHTTPClient(remote)
	return &URIClient{
		address:	address,
		client:		client,
	}
}

func (c *URIClient) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	values, err := argsToURLValues(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.PostForm(c.address+"/"+method, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalResponseBytes(responseBytes, result)
}

func unmarshalResponseBytes(responseBytes []byte, result interface{}) (interface{}, error) {

	var err error
	response := &types.RPCResponse{}
	err = CDC.UnmarshalJSON(responseBytes, response)
	if err != nil {
		return nil, errors.Errorf("Error unmarshalling rpc response: %v", err)
	}
	if response.Error != nil {
		return nil, errors.Errorf("Response error: %v", response.Error)
	}

	err = CDC.UnmarshalJSON(response.Result, result)
	if err != nil {
		return nil, errors.Errorf("Error unmarshalling rpc response result: %v", err)
	}
	return result, nil
}

func argsToURLValues(args map[string]interface{}) (url.Values, error) {
	values := make(url.Values)
	if len(args) == 0 {
		return values, nil
	}
	err := argsToJson(args)
	if err != nil {
		return nil, err
	}
	for key, val := range args {
		values.Set(key, val.(string))
	}
	return values, nil
}

func argsToJson(args map[string]interface{}) error {
	for k, v := range args {
		rt := reflect.TypeOf(v)
		isByteSlice := rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Uint8
		if isByteSlice {
			bytes := reflect.ValueOf(v).Bytes()
			args[k] = fmt.Sprintf("0x%X", bytes)
			continue
		}

		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		args[k] = string(data)
	}
	return nil
}
