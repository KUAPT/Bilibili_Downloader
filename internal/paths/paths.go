package paths

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	globalPaths Paths
	initOnce    sync.Once
	initErr     error
)

type Paths struct {
	Root        string
	ConfigDir   string
	DownloadDir string
	CacheDir    string
	LogPath     string
}

func NewPaths(root string) Paths {
	return Paths{
		Root:        root,
		ConfigDir:   filepath.Join(root, "config"),
		DownloadDir: filepath.Join(root, "Download"),
		CacheDir:    filepath.Join(root, "download_cache"),
		LogPath:     filepath.Join(root, "debug.log"),
	}
}

var resolveExecutable = os.Executable

func ResolveRootDir() (string, error) {
	exePath, err := resolveExecutable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func EnsureDirs(p Paths) error {
	dirs := []string{
		p.ConfigDir,
		p.DownloadDir,
		p.CacheDir,
		filepath.Dir(p.LogPath),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (p Paths) ConfigFilePath() string {
	return filepath.Join(p.ConfigDir, "config.json")
}

func (p Paths) CookieFilePath() string {
	return filepath.Join(p.ConfigDir, "cookies.json")
}

func (p Paths) AudioCachePath() string {
	return filepath.Join(p.CacheDir, "audio_cache")
}

func (p Paths) VideoCachePath() string {
	return filepath.Join(p.CacheDir, "video_cache")
}

func (p Paths) OutputPath(filename string) string {
	return filepath.Join(p.DownloadDir, filename)
}

func (p Paths) UpdateTempPath() string {
	return filepath.Join(p.Root, "update_temp")
}

func Init() error {
	initOnce.Do(func() {
		root, err := ResolveRootDir()
		if err != nil {
			initErr = err
			return
		}
		globalPaths = NewPaths(root)
		initErr = EnsureDirs(globalPaths)
	})
	return initErr
}

func Get() Paths {
	return globalPaths
}

func SetForTesting(p Paths) {
	globalPaths = p
}
