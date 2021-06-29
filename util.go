package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	itoa              = strconv.Itoa
	tmpFileNamePrefix = strconv.FormatInt(time.Now().UnixNano(), 36)
)

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func strSize(size int) string {
	return strSize64(int64(size))
}

func strSize64(size int64) string {
	if mb := int64(1024 * 1024); size >= mb {
		return ftoa(float64(size)*(1.0/float64(mb))) + "MB"
	} else if kb := int64(1024); size >= kb {
		return ftoa(float64(size)*(1.0/float64(kb))) + "KB"
	}
	return fmt.Sprintf("%vB", size)
}

func cmdOk(name string) string {
	s, err := exec.LookPath(name)
	if err != nil {
		println("Skipping '" + name + "': " + err.Error())
	}
	return s
}

func cmdRun(args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)
	if err = cmd.Start(); err == nil {
		err = cmd.Wait()
	}
	if err != nil {
		if errmsg := err.Error(); strings.HasPrefix(errmsg, "exit status ") {
			err = nil // ignore those because it's usually "tool X fails with header or chunk of file Y"..
		} else {
			err = fmt.Errorf("%v: %s", args, errmsg)
		}
	}
	return
}

func viaCmd(printMsg func(...interface{}), dstFileData []byte, cmdAndArgs ...string) (pngData []byte) {
	dstfilepath := tmpFilePath()
	if dstFileData != nil {
		if err := os.WriteFile(dstfilepath, dstFileData, os.ModePerm); err != nil {
			printMsg(err)
			return
		}
	}
	defer func() {
		_ = os.Remove(dstfilepath)
		delete(tmpFiles, dstfilepath)
	}()
	repl := strings.NewReplacer("$dstfilepath$", dstfilepath)
	for i, arg := range cmdAndArgs {
		cmdAndArgs[i] = repl.Replace(arg)
	}
	err := cmdRun(cmdAndArgs...)
	if err != nil {
		printMsg(err)
	} else if pngData, err = os.ReadFile(dstfilepath); err != nil {
		printMsg(err)
	}
	return
}

func tmpFilePath() (filePath string) {
	filePath = "/dev/shm/pb" + tmpFileNamePrefix + strconv.FormatInt(time.Now().UnixNano(), 36) + ".png"
	tmpFiles[filePath] = struct{}{}
	return
}
