package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/selfsigned/timesink/ffprobe"
)

type parameters struct {
	recursive bool
	quiet     bool
	json      bool
}

const (
	usageText = `Usage: timesink [OPTION] [FILE]
Estimate how long it'd take to consume the medias in a directory`
)

func toTimecode(d float64) string {
	hr := int(d) / 3600
	min := (int(d) - hr*3600) / 60
	sec := int(d) % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", hr, min, sec)
}

func getFileLength(params parameters, f string) (length float64) {
	fileInfo, err := ffprobe.Exec(f)
	if err == nil {
		if fileInfo.Format.Duration != "" {
			length, err = strconv.ParseFloat(fileInfo.Format.Duration, 64)

			if err == nil && length != 0 {
				if !params.quiet {
					println(f + " ⏳" + toTimecode(length))
				}
				return
			}

			fmt.Printf("%#v", fileInfo)
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
	return
}

func main() {
	var totalDuration float64
	var params parameters

	flag.BoolVar(&params.recursive, "recursive", false, "Recursively traverse folders")
	flag.BoolVar(&params.quiet, "quiet", false, "quieter output")
	flag.BoolVar(&params.json, "json", false, "json output")
	flag.BoolVar(&params.recursive, "R", false, "recursive (shorthand)")
	flag.BoolVar(&params.quiet, "q", false, "quiet (shorthand)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		println(usageText)
		os.Exit(1)
	}

	// Stop early if ffprobe isn't in PATH
	_, err := ffprobe.GetExecPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: can't access path %q: %v\n", path, err)
				return err
			}
			if info.IsDir() {
				if !params.recursive {
					return filepath.SkipDir
				}
			} else {
				totalDuration += getFileLength(params, path)
			}
			return nil
		})

	}

	if totalDuration < 1 {
		fmt.Fprintf(os.Stderr, "Error: No valid media file found")
		os.Exit(1)
	}

	timeNow := time.Now()

	output := struct {
		Duration float64
		TimeNow  time.Time
		TimeEnd  time.Time
	}{
		totalDuration,
		time.Now(),
		timeNow.Add(time.Second * time.Duration(totalDuration)),
	}

	if params.json {
		out, err := json.MarshalIndent(output, "", "	")
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s", string(out))
	} else {
		fmt.Printf("\nDate after completion:\t%s\nTotal duration:\t\t%s\n",
			output.TimeEnd.Format(time.ANSIC),
			toTimecode(output.Duration))
	}
}
