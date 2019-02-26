package autofile

import (
	"os"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

const autoFileOpenDuration = 1000 * time.Millisecond

type AutoFile struct {
	ID	string
	Path	string
	ticker	*time.Ticker
	mtx	sync.Mutex
	file	*os.File
}

func OpenAutoFile(path string) (af *AutoFile, err error) {
	af = &AutoFile{
		ID:	cmn.RandStr(12) + ":" + path,
		Path:	path,
		ticker:	time.NewTicker(autoFileOpenDuration),
	}
	if err = af.openFile(); err != nil {
		return
	}
	go af.processTicks()
	sighupWatchers.addAutoFile(af)
	return
}

func (af *AutoFile) Close() error {
	af.ticker.Stop()
	err := af.closeFile()
	sighupWatchers.removeAutoFile(af)
	return err
}

func (af *AutoFile) processTicks() {
	for {
		_, ok := <-af.ticker.C
		if !ok {
			return
		}
		af.closeFile()
	}
}

func (af *AutoFile) closeFile() (err error) {
	af.mtx.Lock()
	defer af.mtx.Unlock()

	file := af.file
	if file == nil {
		return nil
	}
	af.file = nil
	return file.Close()
}

func (af *AutoFile) Write(b []byte) (n int, err error) {
	af.mtx.Lock()
	defer af.mtx.Unlock()

	if af.file == nil {
		if err = af.openFile(); err != nil {
			return
		}
	}

	n, err = af.file.Write(b)
	return
}

func (af *AutoFile) Sync() error {
	af.mtx.Lock()
	defer af.mtx.Unlock()

	if af.file == nil {
		if err := af.openFile(); err != nil {
			return err
		}
	}
	return af.file.Sync()
}

func (af *AutoFile) openFile() error {
	file, err := os.OpenFile(af.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	af.file = file
	return nil
}

func (af *AutoFile) Size() (int64, error) {
	af.mtx.Lock()
	defer af.mtx.Unlock()

	if af.file == nil {
		err := af.openFile()
		if err != nil {
			if err == os.ErrNotExist {
				return 0, nil
			}
			return -1, err
		}
	}
	stat, err := af.file.Stat()
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil

}
