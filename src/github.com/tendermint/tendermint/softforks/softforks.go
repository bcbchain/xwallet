package softforks

import (
	"encoding/json"
	"fmt"
	"bcbchain.io/kms"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var AppForkInfo []ForkInfo

type ForkInfo struct {
	Tag			string	`json:"tag,omitempty"`
	EffectBlockHeight	int64	`json:"effectBlockHeight,omitempty"`
	Description		string	`json:"description,omitempty"`
}

func Init() {

	if len(AppForkInfo) == 0 {
		AppForkInfo = make([]ForkInfo, 1)
	} else {

		return
	}
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	dir := filepath.Dir(ex)
	if dir == "" {
		panic(errors.New("Failed to get path of forks file"))
	}

	forksFile := dir + "/tendermint-forks.json"
	if _, err = os.Stat(forksFile); err != nil {

		panic(err.Error())
	}
	sigFile := dir + "/tendermint-forks-signature.json"
	if _, err = os.Stat(forksFile); err != nil {

		panic(err.Error())
	}

	_, err = kms.VerifyFileSign(forksFile, sigFile)
	if err != nil {

		panic(err.Error())
	}

	data, err := ioutil.ReadFile(forksFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &AppForkInfo)
	if err != nil {
		panic(err.Error())
	}
}

func IsForkForV1023233(blockHeight int64) bool {
	for _, forks := range AppForkInfo {
		if forks.Tag == "fork-block#1.0.2.3233" && blockHeight >= forks.EffectBlockHeight {
			return true
		}
	}
	return false
}
