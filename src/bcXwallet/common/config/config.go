package config

import (
	rpcclient "common/rpc/lib/client"
	"fmt"
	"github.com/tendermint/tendermint/rpc/core/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	NodeAddrSlice []string `yaml:"nodeAddrSlice"`
	ChainID       string   `yaml:"chainID"`
	ServerAddr    string   `yaml:"serverAddr"`
	UseHttps      bool     `yaml:"useHttps"`
	OutCertPath   string   `yaml:"outCerPath"`
	KeyStorePath  string   `yaml:"keyStorePath"`
	ChainVersion  string   `yaml:"chainVersion"`

	LoggerScreen bool   `yaml:"loggerScreen"`
	LoggerFile   bool   `yaml:"loggerFile"`
	LoggerLevel  string `yaml:"loggerLevel"`
}

func (c *Config) InitConfig(configFile string) error {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("yamlFile.Get err #%v\n ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v\n", err)
		return err
	}
	if len(c.NodeAddrSlice) == 0 {
		c.NodeAddrSlice = []string{"http://127.0.0.1:37827"}
	}

	if len(c.KeyStorePath) == 0 {
		c.KeyStorePath = "./.keystore"
	}

	return c.initProtocol()
}

func (c *Config) initProtocol() error {
	result := new(core_types.ResultABCIInfo)
	for index, ip := range c.NodeAddrSlice {
		if !strings.HasPrefix(ip, "http") {
			httpsIp := "https://" + ip

			rpc := rpcclient.NewJSONRPCClientEx(httpsIp, "", true)
			_, err := rpc.Call("abci_info", map[string]interface{}{}, result)
			if err == nil {
				c.NodeAddrSlice[index] = httpsIp
				return nil
			}

			httpIp := "http://" + ip

			rpc = rpcclient.NewJSONRPCClientEx(httpIp, "", true)
			_, err = rpc.Call("abci_info", map[string]interface{}{}, result)
			if err == nil {
				c.NodeAddrSlice[index] = httpIp
				return nil
			} else {
				c.NodeAddrSlice = append(c.NodeAddrSlice[:index], c.NodeAddrSlice[index+1:]...)
			}
		}
	}

	return nil
}
