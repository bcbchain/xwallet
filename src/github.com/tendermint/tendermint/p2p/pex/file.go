package pex

import (
	"encoding/json"
	"os"

	cmn "github.com/tendermint/tmlibs/common"
)

type addrBookJSON struct {
	Key	string		`json:"key"`
	Addrs	[]*knownAddress	`json:"addrs"`
}

func (a *addrBook) saveToFile(filePath string) {
	a.Logger.Info("Saving AddrBook to file", "size", a.Size())

	a.mtx.Lock()
	defer a.mtx.Unlock()

	addrs := []*knownAddress{}
	for _, ka := range a.addrLookup {
		addrs = append(addrs, ka)
	}

	aJSON := &addrBookJSON{
		Key:	a.key,
		Addrs:	addrs,
	}

	jsonBytes, err := json.MarshalIndent(aJSON, "", "\t")
	if err != nil {
		a.Logger.Error("Failed to save AddrBook to file", "err", err)
		return
	}
	err = cmn.WriteFileAtomic(filePath, jsonBytes, 0644)
	if err != nil {
		a.Logger.Error("Failed to save AddrBook to file", "file", filePath, "err", err)
	}
}

func (a *addrBook) loadFromFile(filePath string) bool {

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	r, err := os.Open(filePath)
	if err != nil {
		cmn.PanicCrisis(cmn.Fmt("Error opening file %s: %v", filePath, err))
	}
	defer r.Close()
	aJSON := &addrBookJSON{}
	dec := json.NewDecoder(r)
	err = dec.Decode(aJSON)
	if err != nil {
		cmn.PanicCrisis(cmn.Fmt("Error reading file %s: %v", filePath, err))
	}

	a.key = aJSON.Key

	for _, ka := range aJSON.Addrs {
		for _, bucketIndex := range ka.Buckets {
			bucket := a.getBucket(ka.BucketType, bucketIndex)
			bucket[ka.Addr.String()] = ka
		}
		a.addrLookup[ka.ID()] = ka
		if ka.BucketType == bucketTypeNew {
			a.nNew++
		} else {
			a.nOld++
		}
	}
	return true
}
