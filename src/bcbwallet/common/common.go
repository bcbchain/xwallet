package common

import (
	"encoding/json"
	"bcbchain.io/client"
	"bcbchain.io/kms"
	"bcbchain.io/smc"
	"bcbchain.io/tx"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tmlibs/log"
	"bcbwallet/common/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	giWalletConfig	config.Config
	logger		log.Logger
	basicToken	string
	address		smc.Address
)

const RESULTSUCCESS = 200

func InitAll(needGenesis bool) error {
	configFile := "./.config/wallet.yaml"
	moduleName := "bcbwallet"

	_, file := filepath.Split(os.Args[0])
	if strings.HasPrefix(file, "contract") {
		configFile = "./.config/contract.yaml"
		moduleName = "contract_rpc"
	}

	err := giWalletConfig.InitConfig(configFile)
	if err != nil {
		return errors.New("Init config fail err info : " + err.Error())
	}
	initLog(moduleName)
	kms.InitKMS(giWalletConfig.KeyStoreDir, giWalletConfig.Sign_Mode, giWalletConfig.Sign_URL, giWalletConfig.Ca_Path)

	if needGenesis {

		var genesisFilePath string
		if JudgeFile() != nil {
			genesisFilePath, err = GetGenesisFile()
			if err != nil {
				logger.Info("load genesis file err", "errInfo", err)
			}
		} else {
			genesisFilePath = "./.config/" + giWalletConfig.Genesis_file
		}

		err = tx.InitWrapper(genesisFilePath)
		if err != nil {
			logger.Info("load chainId from genesis file err", "errInfo", err)
		}

		err = parseGenesisFile(genesisFilePath)
		if err != nil {
			logger.Info("load basic token from genesis file err", "errInfo", err)
		}
	}

	return nil
}

func initLog(moduleName string) {
	l := log.NewTMLogger("./log", moduleName)
	l.SetOutputToFile(giWalletConfig.Logger_file)
	l.SetOutputToScreen(giWalletConfig.Logger_screen)
	l.AllowLevel(giWalletConfig.Logger_level)
	logger = l
}

func GetConfig() config.Config {
	return giWalletConfig
}

func GetLogger() log.Logger {
	return logger
}

func JudgeFile() error {
	if giWalletConfig.Genesis == "" {
		return errors.New("no cool wallet")
	}
	f, err := os.Open("./.keystore/" + giWalletConfig.Genesis)
	if err == nil {
		defer f.Close()
		return nil
	}
	return err
}

func GetGenesisFile() (string, error) {
	rpc := rpcclient.NewJSONRPCClientEx(giWalletConfig.Node_addrs[0], "", true)
	result := new(core_types.ResultGenesis)

	for {
		_, err := rpc.Call("genesis", map[string]interface{}{}, result)
		if err != nil {
			logger.Debug("Load genesis failed", "error", err)
			time.Sleep(time.Millisecond * 100)
			continue
		} else {
			break
		}
	}

	genesisByte, err := json.Marshal(result.Genesis)
	if err != nil {
		return "", err
	}

	filePath, err := WriteToFile("./.config/", "genesis-file.json", genesisByte)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func parseGenesisFile(genesisFile string) error {
	jsonBytes, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		return err
	}

	type GenesisFile struct {
		AppState struct {
			Token struct {
				Name	string	`json:"name"`
				Address	string	`json:"address"`
			} `json:"token"`
		} `json:"app_state"`
	}

	genesis := &GenesisFile{}
	err = json.Unmarshal(jsonBytes, genesis)
	if err != nil {
		return err
	}

	basicToken = genesis.AppState.Token.Name
	address = genesis.AppState.Token.Address
	return nil
}

func GetBasicTokenName() string {
	return basicToken
}

func GetBasicTokenAddress() smc.Address {
	return address
}

func FuncRecover(l log.Logger, errPtr *error) {
	if err := recover(); err != nil {
		msg := ""
		if errInfo, ok := err.(error); ok {
			msg = errInfo.Error()
		}

		if errInfo, ok := err.(string); ok {
			msg = errInfo
		}

		l.Error("FuncRecover", "error", msg)
		*errPtr = errors.New(msg)
	}
}
