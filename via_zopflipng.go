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
	iter := "4"
	if srcFileLen < (5 * 1024 * 1024) {
		iter = "444"
	}
	return viaCmd(printMsg, nil, "zopflipng", "-m", "--lossy_transparent", "--lossy_8bit", "--filters=01234mepb", "--iterations="+iter,
		srcFilePath,
		"$dstfilepath$")
}
