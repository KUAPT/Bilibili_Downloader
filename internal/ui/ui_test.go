package ui

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"testing"
)

func TestParsePartSelection_AllParts_OrderPreserved(t *testing.T) {
	info := createTestVideoInfo()

	targets, err := parsePartSelection(info, 2, "0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}

	expectedOrder := []string{"Part1", "Part2", "Part3"}
	for i, target := range targets {
		if target.Title != expectedOrder[i] {
			t.Errorf("target[%d].Title = %q, want %q", i, target.Title, expectedOrder[i])
		}
	}
}

func TestParsePartSelection_SpecificParts_OrderPreserved(t *testing.T) {
	info := createTestVideoInfo()

	targets, err := parsePartSelection(info, 2, "3,1,2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}

	expectedOrder := []string{"Part1", "Part2", "Part3"}
	for i, target := range targets {
		if target.Title != expectedOrder[i] {
			t.Errorf("target[%d].Title = %q, want %q (order should follow source, not input)", i, target.Title, expectedOrder[i])
		}
	}
}

func TestParsePartSelection_SinglePart(t *testing.T) {
	info := createTestVideoInfo()

	targets, err := parsePartSelection(info, 2, "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}

	if targets[0].Title != "Part2" {
		t.Errorf("target.Title = %q, want %q", targets[0].Title, "Part2")
	}
}

func TestSelectBestResolution_SkipsHDR(t *testing.T) {
	available := []int{125, 80, 64}

	best := SelectBestResolution(available, true)

	if best != 80 {
		t.Errorf("SelectBestResolution(skipHDR=true) = %d, want 80", best)
	}
}

func TestSelectBestResolution_UsesHDRWhenOnlyOption(t *testing.T) {
	available := []int{125}

	best := SelectBestResolution(available, true)

	if best != 125 {
		t.Errorf("SelectBestResolution with only HDR = %d, want 125", best)
	}
}

func TestSelectBestResolution_SelectsHighest(t *testing.T) {
	available := []int{32, 64, 80, 120}

	best := SelectBestResolution(available, true)

	if best != 120 {
		t.Errorf("SelectBestResolution = %d, want 120 (4K)", best)
	}
}

func createTestVideoInfo() *data_struct.VideoInfoResponse {
	return &data_struct.VideoInfoResponse{
		Data: data_struct.VideoData{
			Title: "Test Video",
			Bvid:  "BV1234567890",
			Cid:   12345,
			Pages: []data_struct.Page{
				{Part: "Part1", Cid: 1001},
				{Part: "Part2", Cid: 1002},
				{Part: "Part3", Cid: 1003},
			},
		},
	}
}
