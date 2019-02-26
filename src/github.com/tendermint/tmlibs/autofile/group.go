package autofile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

const (
	groupCheckDuration	= 5000 * time.Millisecond
	defaultHeadSizeLimit	= 10 * 1024 * 1024
	defaultTotalSizeLimit	= 1 * 1024 * 1024 * 1024
	maxFilesToRemove	= 4
)

type Group struct {
	cmn.BaseService

	ID		string
	Head		*AutoFile
	headBuf		*bufio.Writer
	Dir		string
	ticker		*time.Ticker
	mtx		sync.Mutex
	headSizeLimit	int64
	totalSizeLimit	int64
	minIndex	int
	maxIndex	int
}

func OpenGroup(headPath string) (g *Group, err error) {

	dir := path.Dir(headPath)
	head, err := OpenAutoFile(headPath)
	if err != nil {
		return nil, err
	}

	g = &Group{
		ID:		"group:" + head.ID,
		Head:		head,
		headBuf:	bufio.NewWriterSize(head, 4096*10),
		Dir:		dir,
		ticker:		time.NewTicker(groupCheckDuration),
		headSizeLimit:	defaultHeadSizeLimit,
		totalSizeLimit:	defaultTotalSizeLimit,
		minIndex:	0,
		maxIndex:	0,
	}
	g.BaseService = *cmn.NewBaseService(nil, "Group", g)

	gInfo := g.readGroupInfo()
	g.minIndex = gInfo.MinIndex
	g.maxIndex = gInfo.MaxIndex
	return
}

func (g *Group) OnStart() error {
	g.BaseService.OnStart()
	go g.processTicks()
	return nil
}

func (g *Group) OnStop() {
	g.BaseService.OnStop()
	g.ticker.Stop()
}

func (g *Group) SetHeadSizeLimit(limit int64) {
	g.mtx.Lock()
	g.headSizeLimit = limit
	g.mtx.Unlock()
}

func (g *Group) HeadSizeLimit() int64 {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.headSizeLimit
}

func (g *Group) SetTotalSizeLimit(limit int64) {
	g.mtx.Lock()
	g.totalSizeLimit = limit
	g.mtx.Unlock()
}

func (g *Group) TotalSizeLimit() int64 {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.totalSizeLimit
}

func (g *Group) MaxIndex() int {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.maxIndex
}

func (g *Group) MinIndex() int {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.minIndex
}

func (g *Group) Write(p []byte) (nn int, err error) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.headBuf.Write(p)
}

func (g *Group) WriteLine(line string) error {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	_, err := g.headBuf.Write([]byte(line + "\n"))
	return err
}

func (g *Group) Flush() error {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	err := g.headBuf.Flush()
	if err == nil {
		err = g.Head.Sync()
	}
	return err
}

func (g *Group) processTicks() {
	for {
		_, ok := <-g.ticker.C
		if !ok {
			return
		}
		g.checkHeadSizeLimit()
		g.checkTotalSizeLimit()
	}
}

func (g *Group) stopTicker() {
	g.ticker.Stop()
}

func (g *Group) checkHeadSizeLimit() {
	limit := g.HeadSizeLimit()
	if limit == 0 {
		return
	}
	size, err := g.Head.Size()
	if err != nil {
		panic(err)
	}
	if size >= limit {
		g.RotateFile()
	}
}

func (g *Group) checkTotalSizeLimit() {
	limit := g.TotalSizeLimit()
	if limit == 0 {
		return
	}

	gInfo := g.readGroupInfo()
	totalSize := gInfo.TotalSize
	for i := 0; i < maxFilesToRemove; i++ {
		index := gInfo.MinIndex + i
		if totalSize < limit {
			return
		}
		if index == gInfo.MaxIndex {

			log.Println("WARNING: Group's head " + g.Head.Path + "may grow without bound")
			return
		}
		pathToRemove := filePathForIndex(g.Head.Path, index, gInfo.MaxIndex)
		fileInfo, err := os.Stat(pathToRemove)
		if err != nil {
			log.Println("WARNING: Failed to fetch info for file @" + pathToRemove)
			continue
		}
		err = os.Remove(pathToRemove)
		if err != nil {
			log.Println(err)
			return
		}
		totalSize -= fileInfo.Size()
	}
}

func (g *Group) RotateFile() {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	headPath := g.Head.Path

	if err := g.Head.closeFile(); err != nil {
		panic(err)
	}

	indexPath := filePathForIndex(headPath, g.maxIndex, g.maxIndex+1)
	if err := os.Rename(headPath, indexPath); err != nil {
		panic(err)
	}

	g.maxIndex++
}

func (g *Group) NewReader(index int) (*GroupReader, error) {
	r := newGroupReader(g)
	err := r.SetIndex(index)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type SearchFunc func(line string) (int, error)

func (g *Group) Search(prefix string, cmp SearchFunc) (*GroupReader, bool, error) {
	g.mtx.Lock()
	minIndex, maxIndex := g.minIndex, g.maxIndex
	g.mtx.Unlock()

	for {
		curIndex := (minIndex + maxIndex + 1) / 2

		if minIndex == maxIndex {
			r, err := g.NewReader(maxIndex)
			if err != nil {
				return nil, false, err
			}
			match, err := scanUntil(r, prefix, cmp)
			if err != nil {
				r.Close()
				return nil, false, err
			}
			return r, match, err
		}

		r, err := g.NewReader(curIndex)
		if err != nil {
			return nil, false, err
		}
		foundIndex, line, err := scanNext(r, prefix)
		r.Close()
		if err != nil {
			return nil, false, err
		}

		val, err := cmp(line)
		if err != nil {
			return nil, false, err
		}
		if val < 0 {

			minIndex = foundIndex
		} else if val == 0 {

			r, err := g.NewReader(foundIndex)
			if err != nil {
				return nil, false, err
			}
			match, err := scanUntil(r, prefix, cmp)
			if !match {
				panic("Expected match to be true")
			}
			if err != nil {
				r.Close()
				return nil, false, err
			}
			return r, true, err
		} else {

			maxIndex = curIndex - 1
		}
	}

}

func scanNext(r *GroupReader, prefix string) (int, string, error) {
	for {
		line, err := r.ReadLine()
		if err != nil {
			return 0, "", err
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		index := r.CurIndex()
		return index, line, nil
	}
}

func scanUntil(r *GroupReader, prefix string, cmp SearchFunc) (bool, error) {
	for {
		line, err := r.ReadLine()
		if err != nil {
			return false, err
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		val, err := cmp(line)
		if err != nil {
			return false, err
		}
		if val < 0 {
			continue
		} else if val == 0 {
			r.PushLine(line)
			return true, nil
		} else {
			r.PushLine(line)
			return false, nil
		}
	}
}

func (g *Group) FindLast(prefix string) (match string, found bool, err error) {
	g.mtx.Lock()
	minIndex, maxIndex := g.minIndex, g.maxIndex
	g.mtx.Unlock()

	r, err := g.NewReader(maxIndex)
	if err != nil {
		return "", false, err
	}
	defer r.Close()

GROUP_LOOP:
	for i := maxIndex; i >= minIndex; i-- {
		err := r.SetIndex(i)
		if err != nil {
			return "", false, err
		}

		for {
			line, err := r.ReadLine()
			if err == io.EOF {
				if found {
					return match, found, nil
				}
				continue GROUP_LOOP
			} else if err != nil {
				return "", false, err
			}
			if strings.HasPrefix(line, prefix) {
				match = line
				found = true
			}
			if r.CurIndex() > i {
				if found {
					return match, found, nil
				}
				continue GROUP_LOOP
			}
		}
	}

	return
}

type GroupInfo struct {
	MinIndex	int
	MaxIndex	int
	TotalSize	int64
	HeadSize	int64
}

func (g *Group) ReadGroupInfo() GroupInfo {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.readGroupInfo()
}

func (g *Group) readGroupInfo() GroupInfo {
	groupDir := filepath.Dir(g.Head.Path)
	headBase := filepath.Base(g.Head.Path)
	var minIndex, maxIndex int = -1, -1
	var totalSize, headSize int64 = 0, 0

	dir, err := os.Open(groupDir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	fiz, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	for _, fileInfo := range fiz {
		if fileInfo.Name() == headBase {
			fileSize := fileInfo.Size()
			totalSize += fileSize
			headSize = fileSize
			continue
		} else if strings.HasPrefix(fileInfo.Name(), headBase) {
			fileSize := fileInfo.Size()
			totalSize += fileSize
			indexedFilePattern := regexp.MustCompile(`^.+\.([0-9]{3,})$`)
			submatch := indexedFilePattern.FindSubmatch([]byte(fileInfo.Name()))
			if len(submatch) != 0 {

				fileIndex, err := strconv.Atoi(string(submatch[1]))
				if err != nil {
					panic(err)
				}
				if maxIndex < fileIndex {
					maxIndex = fileIndex
				}
				if minIndex == -1 || fileIndex < minIndex {
					minIndex = fileIndex
				}
			}
		}
	}

	if minIndex == -1 {

		minIndex, maxIndex = 0, 0
	} else {

		maxIndex++
	}
	return GroupInfo{minIndex, maxIndex, totalSize, headSize}
}

func filePathForIndex(headPath string, index int, maxIndex int) string {
	if index == maxIndex {
		return headPath
	}
	return fmt.Sprintf("%v.%03d", headPath, index)
}

type GroupReader struct {
	*Group
	mtx		sync.Mutex
	curIndex	int
	curFile		*os.File
	curReader	*bufio.Reader
	curLine		[]byte
}

func newGroupReader(g *Group) *GroupReader {
	return &GroupReader{
		Group:		g,
		curIndex:	0,
		curFile:	nil,
		curReader:	nil,
		curLine:	nil,
	}
}

func (gr *GroupReader) Close() error {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()

	if gr.curReader != nil {
		err := gr.curFile.Close()
		gr.curIndex = 0
		gr.curReader = nil
		gr.curFile = nil
		gr.curLine = nil
		return err
	}
	return nil
}

func (gr *GroupReader) Read(p []byte) (n int, err error) {
	lenP := len(p)
	if lenP == 0 {
		return 0, errors.New("given empty slice")
	}

	gr.mtx.Lock()
	defer gr.mtx.Unlock()

	if gr.curReader == nil {
		if err = gr.openFile(gr.curIndex); err != nil {
			return 0, err
		}
	}

	var nn int
	for {
		nn, err = gr.curReader.Read(p[n:])
		n += nn
		if err == io.EOF {
			if n >= lenP {
				return n, nil
			}

			if err1 := gr.openFile(gr.curIndex + 1); err1 != nil {
				return n, err1
			}
		} else if err != nil {
			return n, err
		} else if nn == 0 {
			return n, err
		}
	}
}

func (gr *GroupReader) ReadLine() (string, error) {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()

	if gr.curLine != nil {
		line := string(gr.curLine)
		gr.curLine = nil
		return line, nil
	}

	if gr.curReader == nil {
		err := gr.openFile(gr.curIndex)
		if err != nil {
			return "", err
		}
	}

	var linePrefix string
	for {
		bytesRead, err := gr.curReader.ReadBytes('\n')
		if err == io.EOF {

			if err1 := gr.openFile(gr.curIndex + 1); err1 != nil {
				return "", err1
			}
			if len(bytesRead) > 0 && bytesRead[len(bytesRead)-1] == byte('\n') {
				return linePrefix + string(bytesRead[:len(bytesRead)-1]), nil
			}
			linePrefix += string(bytesRead)
			continue
		} else if err != nil {
			return "", err
		}
		return linePrefix + string(bytesRead[:len(bytesRead)-1]), nil
	}
}

func (gr *GroupReader) openFile(index int) error {

	gr.Group.mtx.Lock()
	defer gr.Group.mtx.Unlock()

	if index > gr.Group.maxIndex {
		return io.EOF
	}

	curFilePath := filePathForIndex(gr.Head.Path, index, gr.Group.maxIndex)
	curFile, err := os.Open(curFilePath)
	if err != nil {
		return err
	}
	curReader := bufio.NewReader(curFile)

	if gr.curFile != nil {
		gr.curFile.Close()
	}
	gr.curIndex = index
	gr.curFile = curFile
	gr.curReader = curReader
	gr.curLine = nil
	return nil
}

func (gr *GroupReader) PushLine(line string) {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()

	if gr.curLine == nil {
		gr.curLine = []byte(line)
	} else {
		panic("PushLine failed, already have line")
	}
}

func (gr *GroupReader) CurIndex() int {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	return gr.curIndex
}

func (gr *GroupReader) SetIndex(index int) error {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	return gr.openFile(index)
}

func MakeSimpleSearchFunc(prefix string, target int) SearchFunc {
	return func(line string) (int, error) {
		if !strings.HasPrefix(line, prefix) {
			return -1, errors.New(cmn.Fmt("Marker line did not have prefix: %v", prefix))
		}
		i, err := strconv.Atoi(line[len(prefix):])
		if err != nil {
			return -1, errors.New(cmn.Fmt("Failed to parse marker line: %v", err.Error()))
		}
		if target < i {
			return 1, nil
		} else if target == i {
			return 0, nil
		} else {
			return -1, nil
		}
	}
}
