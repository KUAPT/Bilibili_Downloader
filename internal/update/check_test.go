package update

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckForUpdate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ReleaseResponse{
			TagName: "v2.0.0",
			HTMLURL: "https://github.com/test/release",
			Assets: []Asset{
				{Name: "Bilibili_Downloader.exe", BrowserDownloadURL: "https://example.com/dl.exe"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	info, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !info.HasUpdate {
		t.Error("expected HasUpdate to be true")
	}
	if info.LatestVersion != "v2.0.0" {
		t.Errorf("expected LatestVersion v2.0.0, got %s", info.LatestVersion)
	}
	if info.WindowsAssetURL != "https://example.com/dl.exe" {
		t.Errorf("expected WindowsAssetURL, got %s", info.WindowsAssetURL)
	}
}

func TestCheckForUpdate_NoUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ReleaseResponse{
			TagName: "v1.0.0",
			HTMLURL: "https://github.com/test/release",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	info, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.HasUpdate {
		t.Error("expected HasUpdate to be false")
	}
}

func TestCheckForUpdate_EmptyAssets_NoPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ReleaseResponse{
			TagName: "v2.0.0",
			HTMLURL: "https://github.com/test/release",
			Assets:  []Asset{},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	info, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.WindowsAssetURL != "" {
		t.Errorf("expected empty WindowsAssetURL for empty assets")
	}
}

func TestCheckForUpdate_RateLimitError_NoPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "API rate limit exceeded",
		})
	}))
	defer server.Close()

	_, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err == nil {
		t.Error("expected error for rate limit response")
	}

	var checkErr ErrUpdateCheckFailed
	if _, ok := err.(ErrUpdateCheckFailed); !ok {
		t.Errorf("expected ErrUpdateCheckFailed, got %T", err)
	} else {
		checkErr = err.(ErrUpdateCheckFailed)
		if checkErr.Reason != "API rate limit exceeded" {
			t.Errorf("expected rate limit message, got %s", checkErr.Reason)
		}
	}
}

func TestCheckForUpdate_Non200_NoPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err == nil {
		t.Error("expected error for non-200 response")
	}
}

func TestCheckForUpdate_InvalidJSON_NoPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	_, err := CheckForUpdate(context.Background(), server.URL, "v1.0.0")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSelectWindowsAsset(t *testing.T) {
	tests := []struct {
		name     string
		assets   []Asset
		expected string
	}{
		{
			name: "exact match",
			assets: []Asset{
				{Name: "other.exe", BrowserDownloadURL: "https://other.exe"},
				{Name: "Bilibili_Downloader.exe", BrowserDownloadURL: "https://exact.exe"},
			},
			expected: "https://exact.exe",
		},
		{
			name: "contains bilibili_downloader",
			assets: []Asset{
				{Name: "Bilibili_Downloader_v2.0.0.exe", BrowserDownloadURL: "https://versioned.exe"},
			},
			expected: "https://versioned.exe",
		},
		{
			name: "fallback to any exe",
			assets: []Asset{
				{Name: "SomeOther.exe", BrowserDownloadURL: "https://fallback.exe"},
			},
			expected: "https://fallback.exe",
		},
		{
			name:     "no exe files",
			assets:   []Asset{{Name: "readme.md", BrowserDownloadURL: "https://readme"}},
			expected: "",
		},
		{
			name:     "empty assets",
			assets:   []Asset{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectWindowsAsset(tt.assets)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
