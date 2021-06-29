package main

import (
	"image"
)

func init() {
	if !withZopflipng {
		return
	}
	tools["zop"] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
		return viaZopfli(srcfilepath, printmsg)
	}
}

func viaZopfli(srcFilePath string, printMsg func(...interface{})) []byte {
	return viaCmd(printMsg, nil, "zopflipng", "-m", "--lossy_transparent", "--filters=01234mepb",
		srcFilePath,
		"$dstfilepath$")
}
