package ffprobe

import (
	"encoding/json"
	"os/exec"
)

// MediaInfo is the FFprobe output
type MediaInfo struct {
	Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
	}
}

// Exec runs ffprobe and returns its output as a go struct
func Exec(filepath string) (info MediaInfo, err error) {
	var out []byte

	path, err := GetExecPath()
	if err != nil {
		return
	}
	cmd := exec.Command(
		path,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filepath)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	err = json.Unmarshal(out, &info)
	if err != nil {
		return
	}
	return
}

// GetExecPath returns ffprobe's full path
func GetExecPath() (path string, err error) {
	path, err = exec.LookPath("ffprobe")
	return
}
