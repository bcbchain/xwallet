package common

import (
	"bcXwallet/common/config"
	"blockchain/tx2"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tmlibs/log"
)

var (
	bcbXWalletConfig config.Config
	logger           log.Logger
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
	tx2.Init(bcbXWalletConfig.ChainID)

	CheckChainVersion()

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

	return bcbXWalletConfig.OutCertPath + ".crt", bcbXWalletConfig.OutCertPath + ".key"
}

func CheckChainVersion() {
	cfg := bcbXWalletConfig

	if cfg.ChainVersion != "1" && cfg.ChainVersion != "2" && cfg.ChainVersion != "" {
		fmt.Println("Config file error, please check chainVersion!")
		return
	}

	if cfg.ChainVersion == "2" {
		return
	}

	ChainVersion, err := queryChainVersion()
	if err != nil {
		fmt.Println("Query ChainVersion failed, please check!")
		return
	}

	if ChainVersion == "0" {
		ChainVersion = "1"
	}

	if cfg.ChainVersion != ChainVersion {
		changeChainVersion(ChainVersion)
		bcbXWalletConfig.ChainVersion = ChainVersion
	}
}

func queryChainVersion() (chainVersion string, err error) {
	result := new(core_types.ResultHealth)
	params := map[string]interface{}{}
	err = DoHttpRequestAndParseEx(GetConfig().NodeAddrSlice, "health", params, result)
	if err != nil {
		return "", err
	}

	chainVersion = strconv.FormatInt(result.ChainVersion, 10)
	return
}

func changeChainVersion(chainversion string) {
	configFile := "./.config/bcbXwallet.yaml"

	f, err := os.Open(configFile)
	if err != nil {
		fmt.Println("OpenFile failed, please check!")
		return
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	var Str string
	for {
		line, err := buf.ReadString('\n')
		if strings.HasPrefix(line, "chainVersion:") {
			line = "chainVersion: " + chainversion + "\n"
		}
		Str = Str + line

		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}

	file2, err := os.Create(configFile)
	if err != nil {
		fmt.Println("CreateFile failed, please check!")
		return
	}
	defer file2.Close()

	_, err = file2.WriteString(Str)
	if err != nil {
		fmt.Println("WriteFile failed, please check!")
		return
	}
}
