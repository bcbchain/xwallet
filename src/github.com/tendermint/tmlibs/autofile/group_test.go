package autofile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cmn "github.com/tendermint/tmlibs/common"
)

func createTestGroup(t *testing.T, headSizeLimit int64) *Group {
	testID := cmn.RandStr(12)
	testDir := "_test_" + testID
	err := cmn.EnsureDir(testDir, 0700)
	require.NoError(t, err, "Error creating dir")
	headPath := testDir + "/myfile"
	g, err := OpenGroup(headPath)
	require.NoError(t, err, "Error opening Group")
	g.SetHeadSizeLimit(headSizeLimit)
	g.stopTicker()
	require.NotEqual(t, nil, g, "Failed to create Group")
	return g
}

func destroyTestGroup(t *testing.T, g *Group) {
	err := os.RemoveAll(g.Dir)
	require.NoError(t, err, "Error removing test Group directory")
}

func assertGroupInfo(t *testing.T, gInfo GroupInfo, minIndex, maxIndex int, totalSize, headSize int64) {
	assert.Equal(t, minIndex, gInfo.MinIndex)
	assert.Equal(t, maxIndex, gInfo.MaxIndex)
	assert.Equal(t, totalSize, gInfo.TotalSize)
	assert.Equal(t, headSize, gInfo.HeadSize)
}

func TestCheckHeadSizeLimit(t *testing.T) {
	g := createTestGroup(t, 1000*1000)

	assertGroupInfo(t, g.ReadGroupInfo(), 0, 0, 0, 0)

	for i := 0; i < 999; i++ {
		err := g.WriteLine(cmn.RandStr(999))
		require.NoError(t, err, "Error appending to head")
	}
	g.Flush()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 0, 999000, 999000)

	g.checkHeadSizeLimit()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 0, 999000, 999000)

	err := g.WriteLine(cmn.RandStr(999))
	require.NoError(t, err, "Error appending to head")
	g.Flush()

	g.checkHeadSizeLimit()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 1, 1000000, 0)

	err = g.WriteLine(cmn.RandStr(999))
	require.NoError(t, err, "Error appending to head")
	g.Flush()

	g.checkHeadSizeLimit()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 1, 1001000, 1000)

	for i := 0; i < 999; i++ {
		err = g.WriteLine(cmn.RandStr(999))
		require.NoError(t, err, "Error appending to head")
	}
	g.Flush()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 1, 2000000, 1000000)

	g.checkHeadSizeLimit()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 2, 2000000, 0)

	_, err = g.Head.Write([]byte(cmn.RandStr(999) + "\n"))
	require.NoError(t, err, "Error appending to head")
	g.Flush()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 2, 2001000, 1000)

	g.checkHeadSizeLimit()
	assertGroupInfo(t, g.ReadGroupInfo(), 0, 2, 2001000, 1000)

	destroyTestGroup(t, g)
}

func TestSearch(t *testing.T) {
	g := createTestGroup(t, 10*1000)

	for i := 0; i < 100; i++ {

		_, err := g.Head.Write([]byte(fmt.Sprintf("INFO %v %v\n", i, cmn.RandStr(123))))
		require.NoError(t, err, "Failed to write to head")
		g.checkHeadSizeLimit()
		for j := 0; j < 10; j++ {
			_, err1 := g.Head.Write([]byte(cmn.RandStr(123) + "\n"))
			require.NoError(t, err1, "Failed to write to head")
			g.checkHeadSizeLimit()
		}
	}

	makeSearchFunc := func(target int) SearchFunc {
		return func(line string) (int, error) {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return -1, errors.New("Line did not have 3 parts")
			}
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				return -1, errors.New("Failed to parse INFO: " + err.Error())
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

	for i := 0; i < 100; i++ {
		t.Log("Testing for i", i)
		gr, match, err := g.Search("INFO", makeSearchFunc(i))
		require.NoError(t, err, "Failed to search for line")
		assert.True(t, match, "Expected Search to return exact match")
		line, err := gr.ReadLine()
		require.NoError(t, err, "Failed to read line after search")
		if !strings.HasPrefix(line, fmt.Sprintf("INFO %v ", i)) {
			t.Fatal("Failed to get correct line")
		}

		cur := i + 1
		for {
			line, err := gr.ReadLine()
			if err == io.EOF {
				if cur == 99+1 {

					break
				} else {
					t.Fatal("Got EOF after the wrong INFO #")
				}
			} else if err != nil {
				t.Fatal("Error reading line", err)
			}
			if !strings.HasPrefix(line, "INFO ") {
				continue
			}
			if !strings.HasPrefix(line, fmt.Sprintf("INFO %v ", cur)) {
				t.Fatalf("Unexpected INFO #. Expected %v got:\n%v", cur, line)
			}
			cur++
		}
		gr.Close()
	}

	{
		gr, match, err := g.Search("INFO", makeSearchFunc(-999))
		require.NoError(t, err, "Failed to search for line")
		assert.False(t, match, "Expected Search to not return exact match")
		line, err := gr.ReadLine()
		require.NoError(t, err, "Failed to read line after search")
		if !strings.HasPrefix(line, "INFO 0 ") {
			t.Error("Failed to fetch correct line, which is the earliest INFO")
		}
		err = gr.Close()
		require.NoError(t, err, "Failed to close GroupReader")
	}

	{
		gr, _, err := g.Search("INFO", makeSearchFunc(999))
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, gr)
	}

	destroyTestGroup(t, g)
}

func TestRotateFile(t *testing.T) {
	g := createTestGroup(t, 0)
	g.WriteLine("Line 1")
	g.WriteLine("Line 2")
	g.WriteLine("Line 3")
	g.Flush()
	g.RotateFile()
	g.WriteLine("Line 4")
	g.WriteLine("Line 5")
	g.WriteLine("Line 6")
	g.Flush()

	body1, err := ioutil.ReadFile(g.Head.Path + ".000")
	assert.NoError(t, err, "Failed to read first rolled file")
	if string(body1) != "Line 1\nLine 2\nLine 3\n" {
		t.Errorf("Got unexpected contents: [%v]", string(body1))
	}

	body2, err := ioutil.ReadFile(g.Head.Path)
	assert.NoError(t, err, "Failed to read first rolled file")
	if string(body2) != "Line 4\nLine 5\nLine 6\n" {
		t.Errorf("Got unexpected contents: [%v]", string(body2))
	}

	destroyTestGroup(t, g)
}

func TestFindLast1(t *testing.T) {
	g := createTestGroup(t, 0)

	g.WriteLine("Line 1")
	g.WriteLine("Line 2")
	g.WriteLine("# a")
	g.WriteLine("Line 3")
	g.Flush()
	g.RotateFile()
	g.WriteLine("Line 4")
	g.WriteLine("Line 5")
	g.WriteLine("Line 6")
	g.WriteLine("# b")
	g.Flush()

	match, found, err := g.FindLast("#")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "# b", match)

	destroyTestGroup(t, g)
}

func TestFindLast2(t *testing.T) {
	g := createTestGroup(t, 0)

	g.WriteLine("Line 1")
	g.WriteLine("Line 2")
	g.WriteLine("Line 3")
	g.Flush()
	g.RotateFile()
	g.WriteLine("# a")
	g.WriteLine("Line 4")
	g.WriteLine("Line 5")
	g.WriteLine("# b")
	g.WriteLine("Line 6")
	g.Flush()

	match, found, err := g.FindLast("#")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "# b", match)

	destroyTestGroup(t, g)
}

func TestFindLast3(t *testing.T) {
	g := createTestGroup(t, 0)

	g.WriteLine("Line 1")
	g.WriteLine("# a")
	g.WriteLine("Line 2")
	g.WriteLine("# b")
	g.WriteLine("Line 3")
	g.Flush()
	g.RotateFile()
	g.WriteLine("Line 4")
	g.WriteLine("Line 5")
	g.WriteLine("Line 6")
	g.Flush()

	match, found, err := g.FindLast("#")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "# b", match)

	destroyTestGroup(t, g)
}

func TestFindLast4(t *testing.T) {
	g := createTestGroup(t, 0)

	g.WriteLine("Line 1")
	g.WriteLine("Line 2")
	g.WriteLine("Line 3")
	g.Flush()
	g.RotateFile()
	g.WriteLine("Line 4")
	g.WriteLine("Line 5")
	g.WriteLine("Line 6")
	g.Flush()

	match, found, err := g.FindLast("#")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Empty(t, match)

	destroyTestGroup(t, g)
}

func TestWrite(t *testing.T) {
	g := createTestGroup(t, 0)

	written := []byte("Medusa")
	g.Write(written)
	g.Flush()

	read := make([]byte, len(written))
	gr, err := g.NewReader(0)
	require.NoError(t, err, "failed to create reader")

	_, err = gr.Read(read)
	assert.NoError(t, err, "failed to read data")
	assert.Equal(t, written, read)

	destroyTestGroup(t, g)
}

func TestGroupReaderRead(t *testing.T) {
	g := createTestGroup(t, 0)

	professor := []byte("Professor Monster")
	g.Write(professor)
	g.Flush()
	g.RotateFile()
	frankenstein := []byte("Frankenstein's Monster")
	g.Write(frankenstein)
	g.Flush()

	totalWrittenLength := len(professor) + len(frankenstein)
	read := make([]byte, totalWrittenLength)
	gr, err := g.NewReader(0)
	require.NoError(t, err, "failed to create reader")

	n, err := gr.Read(read)
	assert.NoError(t, err, "failed to read data")
	assert.Equal(t, totalWrittenLength, n, "not enough bytes read")
	professorPlusFrankenstein := professor
	professorPlusFrankenstein = append(professorPlusFrankenstein, frankenstein...)
	assert.Equal(t, professorPlusFrankenstein, read)

	destroyTestGroup(t, g)
}

func TestGroupReaderRead2(t *testing.T) {
	g := createTestGroup(t, 0)

	professor := []byte("Professor Monster")
	g.Write(professor)
	g.Flush()
	g.RotateFile()
	frankenstein := []byte("Frankenstein's Monster")
	frankensteinPart := []byte("Frankenstein")
	g.Write(frankensteinPart)
	g.Flush()

	totalLength := len(professor) + len(frankenstein)
	read := make([]byte, totalLength)
	gr, err := g.NewReader(0)
	require.NoError(t, err, "failed to create reader")

	n, err := gr.Read(read)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, len(professor)+len(frankensteinPart), n, "Read more/less bytes than it is in the group")

	n, err = gr.Read([]byte("0"))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	destroyTestGroup(t, g)
}

func TestMinIndex(t *testing.T) {
	g := createTestGroup(t, 0)

	assert.Zero(t, g.MinIndex(), "MinIndex should be zero at the beginning")

	destroyTestGroup(t, g)
}

func TestMaxIndex(t *testing.T) {
	g := createTestGroup(t, 0)

	assert.Zero(t, g.MaxIndex(), "MaxIndex should be zero at the beginning")

	g.WriteLine("Line 1")
	g.Flush()
	g.RotateFile()

	assert.Equal(t, 1, g.MaxIndex(), "MaxIndex should point to the last file")

	destroyTestGroup(t, g)
}
