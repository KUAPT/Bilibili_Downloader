//go:build linux

package update

import (
	"context"
)

func ApplyUpdate(ctx context.Context, info UpdateInfo) (needRestart bool, newProgramName string, err error) {
	return false, "", ErrUpdateAvailable{
		LatestVersion:  info.LatestVersion,
		ReleasePageURL: info.ReleasePageURL,
	}
}

func LaunchNewVersion(newProgramName string, oldProgramPath string) error {
	return nil
}

func CleanupOldVersion(oldVersionPath string) error {
	return nil
}
