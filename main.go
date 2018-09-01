package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/selfsigned/timesink/src/ffprobe"
)

func toTimecode(d float64) string {
	hr := int(d) / 3600
	min := (int(d) - hr*3600) / 60
	sec := int(d) % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", hr, min, sec)
}

func main() {
	var duration float64
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		println("Usage: timesink [options] files")
		os.Exit(1)
	}
	for _, v := range files {
		out, err := ffprobe.GetFFprobeOut(v)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		length, err := ffprobe.GetFileDuration(out)
		if err == nil {
			duration += length
		}
	}
	println(toTimecode(duration))
}
