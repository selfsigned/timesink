package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/selfsigned/timesink/src/ffprobe"
)

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		println("Usage: timesink [options] files")
		os.Exit(1)
	}
	for _, v := range files {
		out, err := ffprobe.GetFFprobeOut(v)
		if err != nil {
			println("Error: " + err.Error())
			os.Exit(1)
		}
		duration, err := ffprobe.GetFileDuration(out)
		fmt.Printf("%f\n", duration)
	}
}
