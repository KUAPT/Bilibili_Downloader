package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPaths_AllPathsHaveRootPrefix(t *testing.T) {
	root := filepath.FromSlash("/tmp/test_root")
	p := NewPaths(root)

	if p.Root != root {
		t.Errorf("Root = %q, want %q", p.Root, root)
	}

	mustBeUnderRoot(t, root, p.ConfigDir)
	mustBeUnderRoot(t, root, p.DownloadDir)
	mustBeUnderRoot(t, root, p.CacheDir)
	mustBeUnderRoot(t, root, p.LogPath)
}

func mustBeUnderRoot(t *testing.T, root string, path string) {
	t.Helper()

	rel, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatalf("filepath.Rel(%q, %q) error: %v", root, path, err)
	}
	if strings.HasPrefix(rel, "..") {
		t.Fatalf("%q is not under %q", path, root)
	}
}

func TestNewPaths_UsesFilepathJoin(t *testing.T) {
	root := "/tmp/my_root"
	p := NewPaths(root)

	expected := filepath.Join(root, "config")
	if p.ConfigDir != expected {
		t.Errorf("ConfigDir = %q, want %q", p.ConfigDir, expected)
	}

	expected = filepath.Join(root, "Download")
	if p.DownloadDir != expected {
		t.Errorf("DownloadDir = %q, want %q", p.DownloadDir, expected)
	}

	expected = filepath.Join(root, "download_cache")
	if p.CacheDir != expected {
		t.Errorf("CacheDir = %q, want %q", p.CacheDir, expected)
	}
}

func TestResolveRootDir_CanBeInjected(t *testing.T) {
	original := resolveExecutable
	defer func() { resolveExecutable = original }()

	resolveExecutable = func() (string, error) {
		return filepath.FromSlash("/fake/path/to/binary"), nil
	}

	root, err := ResolveRootDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.FromSlash("/fake/path/to")
	if root != expected {
		t.Errorf("ResolveRootDir() = %q, want %q", root, expected)
	}
}

func TestEnsureDirs_CreatesAllDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "paths_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	p := NewPaths(tempDir)
	if err := EnsureDirs(p); err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	dirsToCheck := []string{p.ConfigDir, p.DownloadDir, p.CacheDir}
	for _, dir := range dirsToCheck {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("directory %q does not exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%q is not a directory", dir)
		}
	}
}

func TestPaths_HelperMethods(t *testing.T) {
	root := "/test/root"
	p := NewPaths(root)

	if got := p.ConfigFilePath(); got != filepath.Join(root, "config", "config.json") {
		t.Errorf("ConfigFilePath() = %q", got)
	}

	if got := p.CookieFilePath(); got != filepath.Join(root, "config", "cookies.json") {
		t.Errorf("CookieFilePath() = %q", got)
	}

	if got := p.AudioCachePath(); got != filepath.Join(root, "download_cache", "audio_cache") {
		t.Errorf("AudioCachePath() = %q", got)
	}

	if got := p.VideoCachePath(); got != filepath.Join(root, "download_cache", "video_cache") {
		t.Errorf("VideoCachePath() = %q", got)
	}

	if got := p.OutputPath("video.mp4"); got != filepath.Join(root, "Download", "video.mp4") {
		t.Errorf("OutputPath() = %q", got)
	}
}
