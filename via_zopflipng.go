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
		iter = "333"
	}
	if srcFileLen > (128 * 1024) {
		iter = "222"
	}
	if srcFileLen > (640 * 1024) {
		iter = "111"
	}
	if srcFileLen > (8 * 1024 * 1024) {
		iter = "11"
	}
	if srcFileLen > (32 * 1024 * 1024) {
		iter = "4"
	}
	return viaCmd(printMsg, nil, "zopflipng", "-m", "--lossy_transparent", "--lossy_8bit", "--filters=01234mepb", "--iterations="+iter,
		srcFilePath,
		"$dstfilepath$")
}
