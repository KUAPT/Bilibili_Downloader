package video_processing

import (
	"Bilibili_Downloader/internal/downloader"
	"Bilibili_Downloader/internal/paths"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Transcode(ctx context.Context, p paths.Paths, input downloader.Result, outputRel string) (string, error) {
	if outputRel == "" {
		return "", fmt.Errorf("output path is empty")
	}

	cleanRel := filepath.Clean(outputRel)
	if filepath.IsAbs(cleanRel) || cleanRel == "." || strings.HasPrefix(cleanRel, "..") {
		return "", fmt.Errorf("invalid output path")
	}

	if !strings.HasSuffix(strings.ToLower(cleanRel), ".mp4") {
		cleanRel += ".mp4"
	}

	outputPath := filepath.Join(p.DownloadDir, cleanRel)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	outputPath = avoidOverwrite(outputPath)

	ffmpegPath, cleanup, err := ffmpegExecutable()
	if err != nil {
		return "", err
	}
	if cleanup != nil {
		defer cleanup()
	}

	args := []string{
		"-f", "mp4",
		"-i", input.VideoPath,
		"-i", input.AudioPath,
		"-c:a", "copy",
		"-c:v", "copy",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}

	return outputPath, nil
}

func avoidOverwrite(path string) string {
	if _, err := os.Stat(path); err != nil {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(candidate); err != nil {
			return candidate
		}
	}

	return path
}
