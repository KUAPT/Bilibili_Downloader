package downloader

import (
	"Bilibili_Downloader/internal/paths"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDownload_Success(t *testing.T) {
	const audioSize = 128 * 1024
	const videoSize = 256 * 1024

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var size int
		switch r.URL.Path {
		case "/audio":
			size = audioSize
		case "/video":
			size = videoSize
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(size))
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(make([]byte, size))
	}))
	defer server.Close()

	root := t.TempDir()
	p := paths.NewPaths(root)

	client := &http.Client{}
	result, err := Download(context.Background(), client, p, server.URL+"/audio", server.URL+"/video", Options{ProgressWriter: io.Discard})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkFileSize(t, result.AudioPath, audioSize)
	checkFileSize(t, result.VideoPath, videoSize)
}

func TestDownload_Canceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(1024*1024))
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)

		buf := make([]byte, 32*1024)
		for {
			select {
			case <-r.Context().Done():
				return
			default:
			}
			_, err := w.Write(buf)
			if err != nil {
				return
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(5 * time.Millisecond)
		}
	}))
	defer server.Close()

	root := t.TempDir()
	p := paths.NewPaths(root)

	ctx, cancel := context.WithCancel(context.Background())

	client := &http.Client{}
	errCh := make(chan error, 1)
	go func() {
		_, err := Download(ctx, client, p, server.URL+"/audio", server.URL+"/video", Options{ProgressWriter: io.Discard})
		errCh <- err
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	err := <-errCh
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	if _, statErr := os.Stat(p.CacheDir); statErr != nil {
		t.Fatalf("expected cache dir to exist: %v", statErr)
	}
}

func checkFileSize(t *testing.T, path string, expected int) {
	t.Helper()

	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s failed: %v", path, err)
	}
	if st.Size() != int64(expected) {
		t.Fatalf("expected %d bytes for %s, got %d", expected, path, st.Size())
	}
}
