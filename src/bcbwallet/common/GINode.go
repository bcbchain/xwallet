package common

import (
	"bcbwallet/common/config"
	"bcbwallet/common/http"
)

func GetNodeHttp(config config.Config) *http.HttpClient {
	return http.NewHttpClient(config.Node_addrs[0])
}

func GetNodeHttpFromStr(ipStr string) *http.HttpClient {
	return http.NewHttpClient(ipStr)
}
