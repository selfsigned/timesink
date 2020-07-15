package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/selfsigned/timesink/ffprobe"
)

func toTimecode(d float64) string {
	hr := int(d) / 3600
	min := (int(d) - hr*3600) / 60
	sec := int(d) % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", hr, min, sec)
}

func main() {
	var totalDuration float64
	// var recurse = flag.Bool("recursive", false, "Recursively traverse folders")

	// flag.BoolVar(recurse, "R", false, "(shorthand)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		println("Usage: timesink [options] files")
		os.Exit(1)
	}

	// Stop early if ffprobe isn't in PATH
	_, err := ffprobe.GetExecPath()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// WIP refactor this shit
	for _, f := range files {
		println(f)
		fileInfo, err := ffprobe.Exec(f)
		if err == nil {
			if fileInfo.Format.Duration != "" {
				fileDuration, err := strconv.ParseFloat(fileInfo.Format.Duration, 32)
				if err != nil {
					fmt.Printf("%#v", fileInfo)
					println(err.Error())
					os.Exit(1)
				}
				if fileDuration > 1 {
					totalDuration += fileDuration
				}
			}
		}
	}
	if totalDuration < 1 {
		println("Error: No valid media file found")
		os.Exit(1)
	}

	fmt.Printf("\nDate after completion:\t%s\nTotal duration:\t\t%s\n",
		time.Now().Add(time.Second*time.Duration(totalDuration)).Format(time.ANSIC),
		toTimecode(totalDuration))
}
