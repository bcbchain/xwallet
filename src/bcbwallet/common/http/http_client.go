package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	addr	string
	client	*http.Client
}

func NewHttpClient(addrStr string) *HttpClient {
	clientGo := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	return &HttpClient{
		addr:	addrStr,
		client:	clientGo,
	}
}

func (c *HttpClient) Post(path string, tx []byte) (*http.Response, error) {
	buf := bytes.NewBuffer(tx)
	return c.client.Post(c.addr+"/"+path, "text/json", buf)
}

func (c *HttpClient) Get(path string, paramName []string, paramData []string) (*http.Response, error) {
	var url = "/" + path
	if paramName != nil && paramData != nil && len(paramData) == len(paramName) {
		url += "?"
		for i := 0; i < len(paramName); i++ {
			dataStr := fmt.Sprintf("\"%s\"", paramData[i])
			url += paramName[i] + "=" + dataStr
			if i < len(paramName)-1 {
				url += "&"
			}
		}
	}

	return c.client.Get(c.addr + url)
}

func (c *HttpClient) Parse(res *http.Response) (string, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *HttpClient) QueryInfo(url string) (*http.Response, error) {
	return c.client.Get(url)
}
