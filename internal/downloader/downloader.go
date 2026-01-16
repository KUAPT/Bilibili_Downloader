package downloader

import (
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/pkg/toolkit"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

type Result struct {
	AudioPath string
	VideoPath string
}

type Options struct {
	ProgressWriter io.Writer
}

func Download(ctx context.Context, client *http.Client, p paths.Paths, audioURL, videoURL string, opts Options) (Result, error) {
	if client == nil {
		return Result{}, fmt.Errorf("client is nil")
	}

	if err := os.MkdirAll(p.CacheDir, 0755); err != nil {
		return Result{}, err
	}

	audioPath := p.AudioCachePath()
	videoPath := p.VideoCachePath()

	audioLen, err := headContentLength(ctx, client, audioURL)
	if err != nil {
		return Result{}, err
	}

	videoLen, err := headContentLength(ctx, client, videoURL)
	if err != nil {
		return Result{}, err
	}

	totalSize := audioLen + videoLen
	if totalSize < 0 {
		totalSize = 0
	}

	bar := pb.New64(totalSize)
	bar.Set(pb.SIBytesPrefix, true)
	if opts.ProgressWriter != nil {
		bar.SetWriter(opts.ProgressWriter)
	}
	bar.Start()

	if err := downloadSingleFile(ctx, client, audioURL, audioPath, bar); err != nil {
		bar.Finish()
		return Result{}, err
	}
	if err := downloadSingleFile(ctx, client, videoURL, videoPath, bar); err != nil {
		bar.Finish()
		return Result{}, err
	}

	bar.Finish()

	return Result{AudioPath: audioPath, VideoPath: videoPath}, nil
}

func headContentLength(ctx context.Context, client *http.Client, url string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}
	toolkit.SetVideoHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("head request status error: %s", resp.Status)
	}

	return resp.ContentLength, nil
}

func downloadSingleFile(ctx context.Context, client *http.Client, url string, filePath string, bar *pb.ProgressBar) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	toolkit.SetVideoHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download request status error: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	var reader io.Reader = resp.Body
	if bar != nil {
		reader = bar.NewProxyReader(resp.Body)
	}

	_, copyErr := io.Copy(out, reader)
	if copyErr != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if errors.Is(copyErr, context.Canceled) {
			return context.Canceled
		}
		return copyErr
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
