package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteConfigVals(dir string, vals map[string]string) error {
	data := ""
	for k, v := range vals {
		data = data + fmt.Sprintf("%s = \"%s\"\n", k, v)
	}
	cfile := filepath.Join(dir, "config.toml")
	return ioutil.WriteFile(cfile, []byte(data), 0666)
}

func RunWithArgs(cmd Executable, args []string, env map[string]string) error {
	oargs := os.Args
	oenv := map[string]string{}

	defer func() {
		os.Args = oargs
		for k, v := range oenv {
			os.Setenv(k, v)
		}
	}()

	os.Args = args
	for k, v := range env {

		oenv[k] = os.Getenv(k)
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}

	return cmd.Execute()
}

func RunCaptureWithArgs(cmd Executable, args []string, env map[string]string) (stdout, stderr string, err error) {
	oldout, olderr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout, os.Stderr = wOut, wErr
	defer func() {
		os.Stdout, os.Stderr = oldout, olderr
	}()

	copyStd := func(reader *os.File) *(chan string) {
		stdC := make(chan string)
		go func() {
			var buf bytes.Buffer

			io.Copy(&buf, reader)
			stdC <- buf.String()
		}()
		return &stdC
	}
	outC := copyStd(rOut)
	errC := copyStd(rErr)

	err = RunWithArgs(cmd, args, env)

	wOut.Close()
	wErr.Close()
	stdout = <-*outC
	stderr = <-*errC
	return stdout, stderr, err
}
