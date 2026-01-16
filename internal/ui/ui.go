// Package ui defines the UI abstraction layer.
// Only code in main.go and internal/ui/* is allowed to read from stdin.
package ui

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
)

// DownloadTarget represents a single video to download with stable ordering.
// This replaces map[string]map[string]int64 to guarantee iteration order.
type DownloadTarget struct {
	Title string // Video/episode title
	BVID  string // BV ID
	CID   int64  // CID
}

// UI defines the interface for user interaction.
// Implementations must either be interactive (StdIOUI) or non-interactive (NonInteractiveUI).
type UI interface {
	// PromptBVID asks user for BV ID input.
	// Returns empty string if user wants to exit or input is invalid.
	PromptBVID() (string, error)

	// ConfirmVideo shows video info and asks for confirmation.
	// Returns true if user wants to proceed.
	ConfirmVideo(info *data_struct.VideoInfoResponse) bool

	// PromptMultiPart asks if user wants multi-part download mode.
	// Returns true if user wants multi-part selection.
	PromptMultiPart() bool

	// SelectParts lets user select which parts to download.
	// Returns ordered list of DownloadTarget.
	// kind: 1 for UgcSeason (合集), 2 for Pages (分P)
	SelectParts(info *data_struct.VideoInfoResponse, kind int64) ([]DownloadTarget, error)

	// SelectResolution lets user choose video resolution.
	// Returns video index, video code, and resolution description.
	// defaultMode: -1 = ask user each time, 0 = ask once for batch, 1 = auto highest
	SelectResolution(defaultMode int64, title string, downloadInfo *data_struct.DownloadInfoResponse) (videoIndex int, videoCode int, resolutionDesc string)

	// PromptDefaultResolution asks if user wants default highest resolution for batch download.
	// Returns 1 for auto highest, -1 for manual selection each time.
	PromptDefaultResolution() int64

	// PromptHDRSkip asks if user wants to skip HDR when auto-selecting highest resolution.
	// Returns true if HDR should be skipped.
	PromptHDRSkip() bool

	// PromptContinue asks if user wants to continue downloading more videos.
	PromptContinue() bool

	// PromptUpdateDownload asks if user wants to download an available update.
	// Returns true if user wants to download.
	PromptUpdateDownload(currentVersion, latestVersion string) bool

	// ShowMessage displays a message to the user.
	ShowMessage(msg string)

	// ShowError displays an error message to the user.
	ShowError(msg string)

	// WaitForExit waits for user to press Enter before exiting (interactive mode only).
	WaitForExit()

	// ClearScreen clears the terminal screen.
	ClearScreen()

	// ShowNoPermission shows message when user lacks permission to download.
	ShowNoPermission()

	// IsInteractive returns true if this UI supports user interaction.
	IsInteractive() bool
}
