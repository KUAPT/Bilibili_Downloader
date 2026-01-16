//go:build linux

package video_processing

import (
	"os/exec"
)

func ffmpegExecutable() (string, func() error, error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", nil, ErrFFmpegNotFound{Hint: "ffmpeg 未找到，请安装 ffmpeg"}
	}
	return path, nil, nil
}
