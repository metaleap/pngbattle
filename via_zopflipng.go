package main

import (
	"image"
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
	iter := "222"
	if srcFileLen > (32 * 1024) {
		iter = "111"
	}
	if srcFileLen > (128 * 1024) {
		iter = "77"
	}
	if srcFileLen > (256 * 1024) {
		iter = "55"
	}
	if srcFileLen > (320 * 1024) {
		iter = "44"
	}
	if srcFileLen > (384 * 1024) {
		iter = "33"
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
	if srcFileLen > (1024 * 1024) {
		iter = "1"
	}
	return viaCmd(printMsg, nil, "zopflipng", "-m", "--lossy_transparent", "--lossy_8bit", "--filters=01234mepb", "--iterations="+iter,
		srcFilePath,
		"$dstfilepath$")
}
