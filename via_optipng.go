package main

import (
	"fmt"
	"image"
)

func init() {
	if !withOptipng {
		return
	}
	for f := 0; f <= 5; f++ {
		for i := 0; i <= 1; i++ {
			for zc := 1; zc <= 9; zc++ {
				for zm := 1; zm <= 9; zm++ {
					for zs := 0; zs <= 3; zs++ {
						toolname := fmt.Sprintf("opt f%vi%vzc%vzm%vzs%v", f, i, zc, zm, zs)
						tf, ti, tzc, tzm, tzs := f, i, zc, zm, zs
						tools[toolname] = func(srcfilepath string, srcfiledata []byte, srcfileimg image.Image, printmsg func(...interface{})) []byte {
							return viaOptipng(srcfilepath, printmsg, tf, ti, tzc, tzm, tzs)
						}
					}
				}
			}
		}
	}
}

func viaOptipng(srcFilePath string, printMsg func(...interface{}), f int, i int, zc int, zm int, zs int) []byte {
	return viaCmd(printMsg, nil, "optipng", "-preserve", "-quiet", "-zw32k",
		"-strip", "all",
		fmt.Sprintf("-f%d", f),
		fmt.Sprintf("-i%d", i),
		fmt.Sprintf("-zc%d", zc),
		fmt.Sprintf("-zm%d", zm),
		fmt.Sprintf("-zs%d", zs),
		"-out", "$dstfilepath$",
		srcFilePath)
}
