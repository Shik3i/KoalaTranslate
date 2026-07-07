package github

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		urlStr      string
		expectError bool
		expected    *GitHubFileInfo
	}{
		{
			urlStr:      "https://github.com/Shik3i/KoalaSync/blob/main/extension/_locales/en/messages.json",
			expectError: false,
			expected: &GitHubFileInfo{
				Owner:  "Shik3i",
				Repo:   "KoalaSync",
				Branch: "main",
				Path:   "extension/_locales/en/messages.json",
			},
		},
		{
			urlStr:      "https://raw.githubusercontent.com/Shik3i/KoalaSync/main/extension/_locales/en/messages.json",
			expectError: false,
			expected: &GitHubFileInfo{
				Owner:  "Shik3i",
				Repo:   "KoalaSync",
				Branch: "main",
				Path:   "extension/_locales/en/messages.json",
			},
		},
		{
			urlStr:      "https://malicious.com/Shik3i/KoalaSync/blob/main/messages.json",
			expectError: true,
			expected:    nil,
		},
		{
			urlStr:      "https://github.com/Shik3i/KoalaSync",
			expectError: true,
			expected:    nil,
		},
	}

	for _, tc := range tests {
		actual, err := ParseGitHubURL(tc.urlStr)
		if tc.expectError {
			if err == nil {
				t.Errorf("expected error for URL %q, got nil", tc.urlStr)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for URL %q: %v", tc.urlStr, err)
				continue
			}
			if actual.Owner != tc.expected.Owner ||
				actual.Repo != tc.expected.Repo ||
				actual.Branch != tc.expected.Branch ||
				actual.Path != tc.expected.Path {
				t.Errorf("for URL %q: expected %+v, got %+v", tc.urlStr, tc.expected, actual)
			}
		}
	}
}

func TestToRawURL(t *testing.T) {
	info := &GitHubFileInfo{
		Owner:  "Shik3i",
		Repo:   "KoalaSync",
		Branch: "main",
		Path:   "extension/locales/en.json",
	}

	expected := "https://raw.githubusercontent.com/Shik3i/KoalaSync/main/extension/locales/en.json"
	actual := ToRawURL(info)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}



func TestCreateForkPR(t *testing.T) {
	// Create a test server mock of the GitHub API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1. Get base branch ref
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/git/ref/heads/") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ref":"refs/heads/main","object":{"sha":"abcdef1234567890"}}`))
			return
		}

		// 2. Create branch ref in fork
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/git/refs") {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"ref":"refs/heads/koalatranslate-update-abcdef1"}`))
			return
		}

		// 3. Get existing file in branch
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/contents/") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"sha":"file_blob_sha_123"}`))
			return
		}

		// 4. Commit updated file contents
		if r.Method == "PUT" && strings.Contains(r.URL.Path, "/contents/") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"content":{"name":"messages.json"},"commit":{"sha":"commit_sha_123"}}`))
			return
		}

		// 5. Open a Pull Request
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/pulls") {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"number":42,"html_url":"https://github.com/upstream/repo/pull/42"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	redirectionClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.github.com" {
				req.URL.Scheme = "http"
				req.URL.Host = strings.TrimPrefix(ts.URL, "http://")
			}
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	prNum, prUrl, err := CreateForkPR(
		redirectionClient,
		"mock_token",
		"bot_user",
		"bot_fork_repo",
		"upstream_owner",
		"upstream_repo",
		"main",
		"_locales/de/messages.json",
		"commit msg",
		"pr title",
		"pr body",
		"encoded_content",
	)

	if err != nil {
		t.Fatalf("expected PR creation to succeed, got error: %v", err)
	}

	if prNum != 42 {
		t.Errorf("expected PR number 42, got %d", prNum)
	}

	if prUrl != "https://github.com/upstream/repo/pull/42" {
		t.Errorf("expected PR URL 'https://github.com/upstream/repo/pull/42', got %q", prUrl)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
