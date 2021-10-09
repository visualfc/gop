package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	gobin string
)

func init() {
	var err error
	gobin, err = exec.LookPath("go")
	if err != nil {
		panic("not found go bin in PATH")
	}
}

func main() {
	home := os.Getenv("HOME")
	gopath, err := getGOPATH()
	if err != nil {
		panic(fmt.Sprintf("not found env GOPATH: %w", err))
	}
	goproot, err := getGOPROOT()
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(gopath, home) {
		gopath = strings.Replace(gopath, home, "$HOME", 1)
	}
	if strings.HasPrefix(goproot, home) {
		goproot = strings.Replace(goproot, home, "$HOME", 1)
	}
	fmt.Printf("GOPATH=%v\n", gopath)
	fmt.Printf("GOPROOT=%v\n", goproot)

	fpro := filepath.Join(home, ".profile")
	info, err := os.Stat(fpro)
	if err != nil {
		panic(err)
	}
	env := strings.NewReplacer("$GOPATH", gopath, "$GOPROOT", goproot).Replace(env_template)
	fenv := filepath.Join(home, ".gop/env")
	err = ioutil.WriteFile(fenv, []byte(env), info.Mode())
	if err != nil {
		panic(err)
	}
	fmt.Println("write", fenv)

	data, err := ioutil.ReadFile(fpro)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, `source "$HOME/.gop/env"`) {
			return
		}
	}
	lines = append(lines, `source "$HOME/.gop/env"`)
	err = ioutil.WriteFile(fpro, []byte(strings.Join(lines, "\n")+"\n"), info.Mode())
	if err != nil {
		panic(err)
	}
	fmt.Println("write", fenv)
}

var env_template = `
#!/bin/sh
case ":${PATH}:" in
    *:"$GOPATH/bin":*)
        ;;
    *)
        export PATH="$GOPATH/bin:$PATH"
        ;;
esac
export GOPROOT="$GOPROOT"
`

func getGOPATH() (string, error) {
	data, err := exec.Command(gobin, "env", "GOPATH").Output()
	if err != nil {
		return "", err
	}
	gopath := strings.Split(string(data), ":")[0]
	return gopath, nil
}

func command(gobin string, args ...string) (string, error) {
	data, err := exec.Command(gobin, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getGOPROOT() (root string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		modfile := filepath.Join(dir, "go.mod")
		if hasFile(modfile) {
			if isGoplus(modfile) {
				return dir, nil
			}
			return "", errors.New("current directory is not under goplus root")
		}
		next := filepath.Dir(dir)
		if dir == next {
			return "", errors.New("go.mod not found, please run under goplus root")
		}
		dir = next
	}
}

func isGoplus(modfile string) bool {
	b, err := ioutil.ReadFile(modfile)
	return err == nil && bytes.HasPrefix(b, goplusPrefix)
}

var (
	goplusPrefix = []byte("module github.com/goplus/gop")
)

func hasFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}
