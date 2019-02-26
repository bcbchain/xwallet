package config

import (
	"fmt"
	"bcbchain.io/client"
	"github.com/tendermint/tendermint/rpc/core/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Node_addrs	[]string	`yaml:"node_addrs"`

	Genesis_file	string	`yaml:"genesis_file"`
	Genesis		string	`yaml:"genesis"`
	Sign_Mode	string	`yaml:"sign_mode"`
	Sign_URL	string	`yaml:"sign_url"`
	Logger_screen	bool	`yaml:"logger_screen"`
	Logger_file	bool	`yaml:"logger_file"`
	Logger_level	string	`yaml:"logger_level"`
	KeyStoreDir	string	`yaml:"keyStoreDir"`
	Server_Addr	string	`yaml:"server_addr"`
	Out_Cert_Path	string	`yaml:"out_cer_path"`
	Ca_Path		string	`yaml:"ca_path"`
	Robot_Addr	string	`yaml:"robot_addr"`
	Coin_Type	string	`yaml:"coin_type"`
	Rpc_Mode	string	`yaml:"rpc_mode"`
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
	if len(c.Node_addrs) == 0 {
		c.Node_addrs = []string{"http://127.0.0.1:46657"}
	}

	if c.KeyStoreDir == "" {
		c.KeyStoreDir = "."
	}

	return c.initProtocal()
}

func (c *Config) initProtocal() error {
	result := new(core_types.ResultABCIInfo)
	for index, ip := range c.Node_addrs {
		if !strings.HasPrefix(ip, "http") {
			httpsip := "https://" + ip

			rpc := rpcclient.NewJSONRPCClientEx(httpsip, "", true)
			_, err := rpc.Call("abci_info", map[string]interface{}{}, result)
			if err == nil {
				c.Node_addrs[index] = httpsip
				return nil
			}

			httpip := "http://" + ip

			rpc = rpcclient.NewJSONRPCClientEx(httpip, "", true)
			_, err = rpc.Call("abci_info", map[string]interface{}{}, result)
			if err == nil {
				c.Node_addrs[index] = httpip
				return nil
			} else {
				c.Node_addrs = append(c.Node_addrs[:index], c.Node_addrs[index+1:]...)
			}
		}
	}

	return nil
}
