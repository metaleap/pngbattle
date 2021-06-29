package main

import (
	"fmt"
	"image"
)

func init() {
	if !withOxipng {
		return
	}
	for i := 0; i <= 1; i++ {
		for f := 0; f <= 5; f++ {
			for zc := 1; zc <= 9; zc++ {
				for zs := 0; zs <= 3; zs++ {
					toolname := fmt.Sprintf("oxi i%df%dzc%dzs%d", i, f, zc, zs)
					ti, tf, tzc, tzs := i, f, zc, zs
					tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
						return viaOxipng(srcfilepath, printmsg, ti, tf, tzc, tzs, false)
					}
				}
			}
			toolname := fmt.Sprintf("oxi i%df%d", i, f)
			ti, tf := i, f
			tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
				return viaOxipng(srcfilepath, printmsg, ti, tf, -1, -1, true)
			}
		}
	}
}

func viaOxipng(srcFilePath string, printMsg func(...interface{}), i int, f int, zc int, zs int, defl bool) []byte {
	cmdargs := []string{"oxipng", "--preserve", "--quiet", "--alpha",
		"--strip", "all",
		"--opt", "max",
		"--interlace", itoa(i), //0-1
		"--filters", itoa(f), //0-5
	}
	if defl {
		cmdargs = append(cmdargs, "--libdeflater")
	} else {
		cmdargs = append(cmdargs,
			"--zw", "32k",
			"--zc", itoa(zc), //1-9
			"--zs", itoa(zs), //0-3
		)
	}
	return viaCmd(printMsg, nil,
		append(cmdargs, "--out", "$dstfilepath$", srcFilePath)...)
}
