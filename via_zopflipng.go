package main

import (
	"image"
	"os"
)

func init() {
	if !withZopflipng {
		return
	}
	tools["zop"] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
		return viaZopfli(srcfilepath, len(srcfiledata), printmsg)
	}
}

func viaZopfli(srcFilePath string, srcFileLen int, printMsg func(...interface{})) []byte {
	iter, nofilt := "222", false
	if srcFileLen > (32 * 1024) {
		iter = "111"
	}
	if srcFileLen > (128 * 1024) {
		iter = "77"
	}
	if srcFileLen > (216 * 1024) {
		iter = "55"
	}
	if srcFileLen > (304 * 1024) {
		iter = "44"
	}
	if srcFileLen > (376 * 1024) {
		iter = "22"
	}
	if srcFileLen > (528 * 1024) {
		iter = "11"
	}
	if srcFileLen > (656 * 1024) {
		iter = "7"
	}
	if srcFileLen > (792 * 1024) {
		iter = "4"
	}
	if srcFileLen > (896 * 1024) {
		iter = "2"
	}
	if srcFileLen > (1016 * 1024) {
		iter, nofilt = "1", true
	}
	cmdandargs := []string{"zopflipng", "--filters=01234mepb", "--lossy_transparent", "--lossy_8bit", "--iterations=" + iter,
		srcFilePath,
		"$dstfilepath$"}
	if nofilt || os.Getenv("NOFILT") != "" {
		cmdandargs = append([]string{cmdandargs[0]}, cmdandargs[2:]...)
	}
	return viaCmd(printMsg, nil, cmdandargs...)
}
