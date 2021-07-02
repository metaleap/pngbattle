package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	withGoStdLib    = false
	withAdvancecomp = false
	withImagemagick = false
	withOptipng     = false
	withOxipng      = false
	withZopflipng   = true

	timeStarted = time.Now()
	tmpFiles    = map[string]struct{}{}
	tools       = make(map[string]func(string, []byte, image.Image, func(...interface{})) []byte, 8192)
	stats       struct {
		totalSizeSrc int64
		totalSizeDst int64
		tools        map[string][2]int
	}
)

func main() {
	stats.tools = map[string][2]int{}
	if len(tools) == 0 {
		panic("No tools: all disabled")
	}
	defer func() {
		for tmpfilepath := range tmpFiles {
			_ = os.Remove(tmpfilepath)
		}
		if len(os.Args) > 2 {
			println("PNGBATTLE STATS:", strSize64(stats.totalSizeSrc), "down to", strSize64(stats.totalSizeDst), "("+ftoa(100.0/(float64(stats.totalSizeSrc)/float64(stats.totalSizeDst)))+"%)", "in", time.Now().Sub(timeStarted).String())
			totalsavings := stats.totalSizeSrc - stats.totalSizeDst
			for toolname, i2 := range stats.tools {
				percsavings, percfiles := 0.0, 100.0/(float64(len(os.Args)-1)/float64(i2[0]))
				if savings := i2[1]; savings != 0 {
					percsavings = 100.0 / (float64(totalsavings) / float64(savings))
				}
				if toolname == "" {
					toolname = "(no-op)"
				}
				println("\t"+toolname+"\t\twon", ftoa(percfiles)+"% of files,", ftoa(percsavings)+"% of size-savings")
			}
		}
	}()

	if len(os.Args) > 2 {
		println("\n_________\nPNGBATTLE started:", len(os.Args)-1, "×", len(tools), "=", (len(os.Args)-1)*len(tools), "attempts...")
	}
	for _, pngfilepath := range os.Args[1:] {
		if fileinfo, err := os.Stat(pngfilepath); err != nil {
			panic(err)
		} else {
			stats.totalSizeSrc += fileinfo.Size()
		}
	}
	srcdonesize, numremaining := int64(0), len(os.Args)-1
	for _, pngfilepath := range os.Args[1:] {
		pngdata, err := os.ReadFile(pngfilepath)
		if err != nil {
			println("OS ReadFile", err.Error())
			continue
		}
		srcsize := int64(len(pngdata))
		pngdata = pngMin(pngfilepath, pngdata)
		if srcdonesize, numremaining = srcdonesize+srcsize, numremaining-1; numremaining == 0 {
			println()
		} else {
			sizeremaining := stats.totalSizeSrc - srcdonesize
			timetaken := time.Now().Sub(timeStarted)
			timeremaining := (float64(timetaken) / float64(srcdonesize)) * float64(sizeremaining)
			println("\t\tremaining: " + itoa(numremaining) + "/" + itoa(len(os.Args)-1) + " file(s), " + strSize64(sizeremaining) + "/" + strSize64(stats.totalSizeSrc) + ", ~" + time.Duration(timeremaining).String() + "...")
		}
		if len(pngdata) == 0 {
			stats.totalSizeDst += srcsize
		} else {
			stats.totalSizeDst += int64(len(pngdata))
			dstfilepath := pngfilepath + "." + strconv.FormatInt(time.Now().UnixNano(), 36) + ".png"
			if err = writeFile(dstfilepath, pngdata); err != nil {
				_ = os.Remove(dstfilepath)
				println("OS WriteFile", err)
				// } else if true {
				// 	_ = os.Remove(dstfilepath)
			} else if err = os.Rename(dstfilepath, pngfilepath); err != nil {
				_ = os.Remove(dstfilepath)
				println("OS Rename", err)
			}
		}
	}
}

func pngMin(srcFilePath string, srcFileData []byte) []byte {
	type result struct {
		size    int
		pngData []byte
	}
	srcfilesize := len(srcFileData)
	println("\033[4m" + srcFilePath + "\033[0m\t" + strSize(srcfilesize))
	results := make(map[string]result, len(tools))

	srcfileimg, err := png.Decode(bytes.NewReader(srcFileData))
	if err != nil {
		println("\tNOPNG\t\t" + srcFilePath)
		return nil
	}
	var work sync.WaitGroup
	var mu sync.Mutex
	minsize := uint32(srcfilesize)
	for toolname, fn := range tools {
		printmsg := func(args ...interface{}) {
			mu.Lock()
			print("\tMSGBY '"+toolname+"' "+srcFilePath, "\t\t")
			for _, arg := range args {
				print("", fmt.Sprintf("%v", arg))
			}
			print("\n")
			mu.Unlock()
		}
		if pngdata := fn(srcFilePath, srcFileData, srcfileimg, printmsg); len(pngdata) != 0 {
			if dstimg, err := png.Decode(bytes.NewReader(pngdata)); err != nil {
				mu.Lock()
				println("\tBADBY '" + toolname + "' for " + srcFilePath + ": " + err.Error())
				mu.Unlock()
			} else if !dstimg.Bounds().Eq(srcfileimg.Bounds()) {
				mu.Lock()
				println("\tBUGBY '"+toolname+"' for", srcFilePath, ": src bounds", srcfileimg.Bounds().String(), "BUT dst bounds", dstimg.Bounds().String())
				mu.Unlock()
			} else {
				work.Add(1)
				go func(dstimg image.Image, toolname string, pngdata []byte) {
					allok, res := true, result{size: len(pngdata), pngData: pngdata}
					if ms := atomic.LoadUint32(&minsize); uint32(res.size) >= ms {
						res.pngData = nil
					} else {
						for x := 0; x < dstimg.Bounds().Max.X && allok; x++ {
							for y := 0; y < dstimg.Bounds().Max.Y && allok; y++ {
								dr, dg, db, da := dstimg.At(x, y).RGBA()
								sr, sg, sb, sa := srcfileimg.At(x, y).RGBA()
								if allok = (dr == sr) && (dg == sg) && (db == sb) && (da == sa); !allok {
									mu.Lock()
									println("\tBUGBY '"+toolname+"' for", srcFilePath, ": rgba diff at", x, y)
									mu.Unlock()
								}
							}
						}
					}
					if allok {
						for done, ms := false, atomic.LoadUint32(&minsize); !done; ms = atomic.LoadUint32(&minsize) {
							rs := uint32(res.size)
							done = (rs >= ms) || atomic.CompareAndSwapUint32(&minsize, ms, rs)
						}
						mu.Lock()
						for tn, r := range results {
							if r.size > res.size {
								results[tn] = result{size: r.size, pngData: nil}
							}
							if r.pngData != nil && r.size <= res.size {
								res.pngData = nil
								break
							}
						}
						results[toolname] = res
						mu.Unlock()
					}
					work.Done()
				}(dstimg, toolname, pngdata)
			}
		}
	}
	work.Wait()

	minnames, minresult := []string{""}, result{size: srcfilesize}
	for toolname, result := range results {
		if result.size < minresult.size {
			minnames, minresult.size = []string{toolname}, result.size
		} else if result.size == minresult.size {
			minnames = append(minnames, toolname)
		}
		if result.size <= minresult.size && result.pngData != nil {
			minresult.pngData = result.pngData
		}
	}

	if minnames[0] == "" {
		i2 := stats.tools[""]
		stats.tools[""] = [2]int{i2[0] + 1, 0}
		print("\t(no compression achieved)")
		return nil
	}
	print("\t"+strSize(minresult.size), " ("+ftoa(100.0/(float64(srcfilesize)/float64(minresult.size)))+"%) via '"+strings.Join(minnames, "', '")+"'.")

	commonprefix, maxlen := "", 0
	for _, name := range minnames {
		if maxlen == 0 || len(name) < maxlen {
			maxlen = len(name)
		}
	}
	for sl := maxlen; sl > 0; sl-- {
		alleq, pref := true, minnames[0][:sl]
		for _, name := range minnames {
			if alleq = (name[:sl] == pref); !alleq {
				break
			}
		}
		if alleq {
			commonprefix = pref
			break
		}
	}
	if commonprefix = strings.TrimSpace(commonprefix); commonprefix == "" && len(os.Args) > 2 {
		print("\nNo common prefix: '" + strings.Join(minnames, "', '") + "'")
	} else {
		i2 := stats.tools[commonprefix]
		stats.tools[commonprefix] = [2]int{i2[0] + 1, i2[1] + (srcfilesize - minresult.size)}
	}
	if len(minresult.pngData) == 0 {
		panic("BUG: len(minresult.pngData)==0")
	}
	return minresult.pngData
}
