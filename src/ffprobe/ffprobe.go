package ffprobe

import (
	"os/exec"
)

// GetFFprobeOut returns the duration of a media file using ffprobe
func GetFFprobeOut(filepath string) (out []byte, err error) {
	path, err := exec.LookPath("ffprobe")
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
	return
}
