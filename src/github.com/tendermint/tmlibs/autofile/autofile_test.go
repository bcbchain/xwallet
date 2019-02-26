package autofile

import (
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

func TestSIGHUP(t *testing.T) {

	file, name := cmn.Tempfile("sighup_test")
	if err := file.Close(); err != nil {
		t.Fatalf("Error creating tempfile: %v", err)
	}

	af, err := OpenAutoFile(name)
	if err != nil {
		t.Fatalf("Error creating autofile: %v", err)
	}

	_, err = af.Write([]byte("Line 1\n"))
	if err != nil {
		t.Fatalf("Error writing to autofile: %v", err)
	}
	_, err = af.Write([]byte("Line 2\n"))
	if err != nil {
		t.Fatalf("Error writing to autofile: %v", err)
	}

	err = os.Rename(name, name+"_old")
	if err != nil {
		t.Fatalf("Error moving autofile: %v", err)
	}

	oldSighupCounter := atomic.LoadInt32(&sighupCounter)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	for atomic.LoadInt32(&sighupCounter) == oldSighupCounter {
		time.Sleep(time.Millisecond * 10)
	}

	_, err = af.Write([]byte("Line 3\n"))
	if err != nil {
		t.Fatalf("Error writing to autofile: %v", err)
	}
	_, err = af.Write([]byte("Line 4\n"))
	if err != nil {
		t.Fatalf("Error writing to autofile: %v", err)
	}
	if err := af.Close(); err != nil {
		t.Fatalf("Error closing autofile")
	}

	if body := cmn.MustReadFile(name + "_old"); string(body) != "Line 1\nLine 2\n" {
		t.Errorf("Unexpected body %s", body)
	}
	if body := cmn.MustReadFile(name); string(body) != "Line 3\nLine 4\n" {
		t.Errorf("Unexpected body %s", body)
	}
}
