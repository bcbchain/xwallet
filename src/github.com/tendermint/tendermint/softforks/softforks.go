package softforks

// initial version copied from gichain

import (
	"common/sig"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var AppForkInfo []ForkInfo

//具体含义请参考 gichain.yaml
type ForkInfo struct {
	Tag               string `json:"tag,omitempty"`               //Tag, contains the former released version
	EffectBlockHeight int64  `json:"effectBlockHeight,omitempty"` // Effect Block Height
	Description       string `json:"description,omitempty"`       // Description for the fork
}

// 不自動init了，會給很多其他引用模塊帶來麻煩
func Init() {

	if len(AppForkInfo) == 0 {
		AppForkInfo = make([]ForkInfo, 1)
	} else {
		// already initialized, return
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
		//File doesn't exist, terminate the process
		panic(err.Error())
	}
	sigFile := dir + "/tendermint-forks.json.sig"
	if _, err = os.Stat(forksFile); err != nil {
		//File doesn't exist, terminate the process
		panic(err.Error())
	}
	// Verify Fork.json
	_, err = sig.VerifyTextFile(forksFile, sigFile)
	if err != nil {
		// Failed to verify, terminate the process
		panic(err.Error())
	}

	// Notes: be careful of permission of the file, should be 444 or 644
	data, err := ioutil.ReadFile(forksFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &AppForkInfo)
	if err != nil {
		panic(err.Error())
	}
}

// Adds the softfork to show all of rewards in block
func IsForkForV1023233(blockHeight int64) bool {
	for _, forks := range AppForkInfo {
		if forks.Tag == "fork-block#1.0.2.3233" && blockHeight >= forks.EffectBlockHeight {
			return true
		}
	}
	return false
}
