package main

import (
	"fmt"
	"image"
)

func init() {
	if !withImagemagick {
		return
	}
	for l := 0; l <= 9; l++ {
		for f := 0; f <= 5; f++ {
			toolname := fmt.Sprintf("imc %v%v", l, f)
			tl, tf := l, f
			tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
				return viaImagemagick(srcfilepath, printmsg, tl, tf)
			}
		}
	}
}

func viaImagemagick(srcFilePath string, printMsg func(...interface{}), level int, filter int) []byte {
	return viaCmd(printMsg, nil, "convert", "-quiet",
		"-quality", fmt.Sprintf("%d%d", level, filter),
		srcFilePath,
		"$dstfilepath$")
}
