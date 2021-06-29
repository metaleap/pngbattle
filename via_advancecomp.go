package main

import (
	"fmt"
	"image"
)

func init() {
	if !withAdvancecomp {
		return
	}
	for s := 0; s <= 4; s++ {
		toolname := fmt.Sprintf("adv %d", s)
		ts := s
		tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
			return viaAdvancecomp(srcfiledata, printmsg, ts)
		}
	}
}

func viaAdvancecomp(srcFileData []byte, printMsg func(...interface{}), s int) []byte {
	return viaCmd(printMsg, srcFileData, "advpng", "--recompress", "--quiet",
		fmt.Sprintf("-%d", s),
		"$dstfilepath$")
}
