package common

import (
	"bcbXwallet/common/config"
	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/log"
	"bcbwallet/common"
	"os"
)

var (
	bcbXWalletConfig	config.Config
	logger			log.Logger
)

func InitAll() error {
	configFile := "./.config/bcbXwallet.yaml"
	moduleName := "bcbXwallet"

	err := bcbXWalletConfig.InitConfig(configFile)
	if err != nil {
		return errors.New("Init config fail err info : " + err.Error())
	}
	initLog(moduleName)

	if bcbXWalletConfig.ChainID == "" {
		return errors.New(" chainId cannot be empty")
	}
	crypto.SetChainId(bcbXWalletConfig.ChainID)

	return nil
}

func initLog(moduleName string) {
	l := log.NewTMLogger("./log", moduleName)
	l.SetOutputToFile(bcbXWalletConfig.LoggerFile)
	l.SetOutputToScreen(bcbXWalletConfig.LoggerScreen)
	l.AllowLevel(bcbXWalletConfig.LoggerLevel)
	logger = l
}

func GetConfig() config.Config {
	return bcbXWalletConfig
}

func GetLogger() log.Logger {
	return logger
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

func OutCertFileIsExist() (string, string) {
	crtPath := "./.config/server.crt"
	keyPath := "./.config/server.key"

	_, err := os.Stat(bcbXWalletConfig.OutCertPath + ".crt")
	if err != nil {
		return crtPath, keyPath
	}

	_, err = os.Stat(bcbXWalletConfig.OutCertPath + ".key")
	if err != nil {
		return crtPath, keyPath
	}

	return common.GetConfig().Out_Cert_Path + ".crt", bcbXWalletConfig.OutCertPath + ".key"
}
