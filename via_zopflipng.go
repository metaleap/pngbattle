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
	iter := "444"
	if srcFileLen > (32 * 1024) {
		iter = "222"
	}
	if srcFileLen > (128 * 1024) {
		iter = "77"
	}
	if srcFileLen > (544 * 1024) {
		iter = "22"
	}
	if srcFileLen > (2 * 1024 * 1024) {
		iter = "5"
	}
	return viaCmd(printMsg, nil, "zopflipng", "-m", "--lossy_transparent", "--lossy_8bit", "--filters=01234mepb", "--iterations="+iter,
		srcFilePath,
		"$dstfilepath$")
}
