package ffprobe

import (
	"encoding/json"
	"strconv"
)

type FileInfo struct {
	Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
	}
}

// GetFileDuration returns a file's duration
// from a JSON object outputted by ffprobe
func GetFileDuration(data []byte) (duration float64, err error) {
	var d FileInfo

	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}
	duration, err = strconv.ParseFloat(d.Format.Duration, 32)
	return
}
