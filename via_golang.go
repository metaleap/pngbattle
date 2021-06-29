package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

func init() {
	if !withGo {
		return
	}
	for level := png.CompressionLevel(-1); level <= 9; level++ {
		toolname := fmt.Sprintf("go %v", level)
		tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
			return viaGolang(srcfilepath, srcfiledata, srcfileimg, png.CompressionLevel(level), printmsg)
		}
	}
}

func viaGolang(srcFilePath string, srcFileData []byte, srcFileImg image.Image, compressionLevel png.CompressionLevel, printMsg func(...interface{})) []byte {
	pngenc := png.Encoder{CompressionLevel: compressionLevel}
	buf := bytes.NewBuffer(make([]byte, 0, len(srcFileData)))
	if err := pngenc.Encode(buf, srcFileImg); err != nil {
		printMsg(err)
	}
	return buf.Bytes()
}
