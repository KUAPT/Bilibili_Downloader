//go:build windows

package video_processing

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed ffmpeg/ffmpeg
var embeddedFFmpeg embed.FS

func ffmpegExecutable() (string, func() error, error) {
	data, err := embeddedFFmpeg.ReadFile("ffmpeg/ffmpeg")
	if err != nil {
		return "", nil, err
	}

	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		return "", nil, err
	}

	ffmpegPath := filepath.Join(tempDir, "ffmpeg.exe")
	if err := os.WriteFile(ffmpegPath, data, 0755); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", nil, err
	}

	cleanup := func() error {
		return os.RemoveAll(tempDir)
	}

	return ffmpegPath, cleanup, nil
}
