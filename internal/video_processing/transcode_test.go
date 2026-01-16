package video_processing

import (
	"Bilibili_Downloader/internal/downloader"
	"Bilibili_Downloader/internal/paths"
	"context"
	"os"
	"path/filepath"
	"testing"
)

type fakeCleanup struct {
	called bool
}

func TestTranscode_InvalidOutputPath(t *testing.T) {
	p := paths.NewPaths(t.TempDir())
	_, err := Transcode(context.Background(), p, downloader.Result{AudioPath: "a", VideoPath: "v"}, "../escape")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAvoidOverwrite(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.mp4")
	if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := avoidOverwrite(p)
	if got == p {
		t.Fatalf("expected different path")
	}
}
