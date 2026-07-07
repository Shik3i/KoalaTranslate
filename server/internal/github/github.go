package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GitHubFileInfo holds parsed details from a GitHub URL.
type GitHubFileInfo struct {
	Owner  string
	Repo   string
	Branch string
	Path   string
}

// ParseGitHubURL extracts details from supported GitHub URLs.
func ParseGitHubURL(urlStr string) (*GitHubFileInfo, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	hostname := strings.ToLower(u.Hostname())
	if hostname != "github.com" && hostname != "raw.githubusercontent.com" {
		return nil, errors.New("only github.com and raw.githubusercontent.com domains are allowed")
	}

	// Remove leading slash and split
	path := strings.TrimPrefix(u.Path, "/")
	parts := strings.Split(path, "/")

	if hostname == "github.com" {
		// Expect: github.com/{owner}/{repo}/blob/{branch}/{path...}
		if len(parts) < 5 || parts[2] != "blob" {
			return nil, errors.New("invalid github.com URL format; expected /owner/repo/blob/branch/path")
		}
		return &GitHubFileInfo{
			Owner:  parts[0],
			Repo:   parts[1],
			Branch: parts[3],
			Path:   strings.Join(parts[4:], "/"),
		}, nil
	}

	// Hostname is raw.githubusercontent.com
	// Expect: raw.githubusercontent.com/{owner}/{repo}/{branch}/{path...}
	if len(parts) < 4 {
		return nil, errors.New("invalid raw.githubusercontent.com URL format; expected /owner/repo/branch/path")
	}
	return &GitHubFileInfo{
		Owner:  parts[0],
		Repo:   parts[1],
		Branch: parts[2],
		Path:   strings.Join(parts[3:], "/"),
	}, nil
}

// ToRawURL converts a GitHubFileInfo into a raw githubusercontent URL.
func ToRawURL(info *GitHubFileInfo) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", info.Owner, info.Repo, info.Branch, info.Path)
}

// FetchGitHubFile fetches the file content from a GitHub URL after validating it.
func FetchGitHubFile(client *http.Client, urlStr string) ([]byte, error) {
	info, err := ParseGitHubURL(urlStr)
	if err != nil {
		return nil, err
	}

	rawURL := ToRawURL(info)
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch raw content, HTTP status: %d", resp.StatusCode)
	}

	// Limit read to 10MB
	limitReader := io.LimitReader(resp.Body, 10*1024*1024)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}

	return data, nil
}

// CreateForkPR handles creating a Pull Request from a bot fork to upstream.
// It executes the following API calls:
// 1. Create/verify a branch on the bot fork repo.
// 2. Commit the new file content.
// 3. Open a pull request from GITHUB_BOT_OWNER:branch to upstream main.
func CreateForkPR(
	client *http.Client,
	token string,
	botOwner string,
	botForkRepo string,
	upstreamOwner string,
	upstreamRepo string,
	upstreamBranch string,
	filePath string,
	commitMessage string,
	prTitle string,
	prBody string,
	fileContent string,
) (int, string, error) {
	if token == "" || botOwner == "" || botForkRepo == "" {
		return 0, "", errors.New("github credentials not fully configured")
	}

	apiBase := "https://api.github.com"

	// Helper for HTTP requests with auth headers
	sendReq := func(method, endpoint string, bodyData interface{}) ([]byte, int, error) {
		var bodyReader io.Reader
		if bodyData != nil {
			jsonData, err := json.Marshal(bodyData)
			if err != nil {
				return nil, 0, err
			}
			bodyReader = bytes.NewBuffer(jsonData)
		}

		req, err := http.NewRequest(method, apiBase+endpoint, bodyReader)
		if err != nil {
			return nil, 0, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if bodyData != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		return respBody, resp.StatusCode, err
	}

	// 1. Verify upstream repo base branch ref to get latest commit SHA
	refUrl := fmt.Sprintf("/repos/%s/%s/git/ref/heads/%s", upstreamOwner, upstreamRepo, upstreamBranch)
	respBody, status, err := sendReq("GET", refUrl, nil)
	if err != nil || status != http.StatusOK {
		return 0, "", fmt.Errorf("failed to fetch upstream ref (%d): %s", status, string(respBody))
	}

	var refData struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := json.Unmarshal(respBody, &refData); err != nil {
		return 0, "", err
	}
	baseSHA := refData.Object.SHA

	// Generate a unique branch name
	newBranchName := fmt.Sprintf("koalatranslate-update-%s", baseSHA[:7])

	// 2. Create the branch in the bot fork pointing to the base SHA
	createRefUrl := fmt.Sprintf("/repos/%s/%s/git/refs", botOwner, botForkRepo)
	refPayload := map[string]interface{}{
		"ref": "refs/heads/" + newBranchName,
		"sha": baseSHA,
	}
	respBody, status, err = sendReq("POST", createRefUrl, refPayload)
	// If branch already exists (422), we can just proceed to write over it
	if err != nil {
		return 0, "", err
	}
	if status != http.StatusCreated && status != http.StatusUnprocessableEntity {
		return 0, "", fmt.Errorf("failed to create branch in fork (%d): %s", status, string(respBody))
	}

	// 3. Get existing file blob SHA on the branch in the bot fork (if exists)
	fileUrl := fmt.Sprintf("/repos/%s/%s/contents/%s?ref=%s", botOwner, botForkRepo, filePath, newBranchName)
	respBody, status, err = sendReq("GET", fileUrl, nil)
	var fileSHA string
	if err == nil && status == http.StatusOK {
		var fileInfo struct {
			SHA string `json:"sha"`
		}
		_ = json.Unmarshal(respBody, &fileInfo)
		fileSHA = fileInfo.SHA
	}

	// 4. Commit updated messages.json file contents
	commitPayload := map[string]interface{}{
		"message": commitMessage,
		"content": fileContent, // must be base64 encoded by caller
		"branch":  newBranchName,
	}
	if fileSHA != "" {
		commitPayload["sha"] = fileSHA
	}

	contentUrl := fmt.Sprintf("/repos/%s/%s/contents/%s", botOwner, botForkRepo, filePath)
	respBody, status, err = sendReq("PUT", contentUrl, commitPayload)
	if err != nil || (status != http.StatusOK && status != http.StatusCreated) {
		return 0, "", fmt.Errorf("failed to commit file to fork (%d): %s", status, string(respBody))
	}

	// 5. Open a Pull Request from botForkOwner:branch to upstream base branch
	prUrl := fmt.Sprintf("/repos/%s/%s/pulls", upstreamOwner, upstreamRepo)
	prPayload := map[string]interface{}{
		"title": prTitle,
		"body":  prBody,
		"head":  fmt.Sprintf("%s:%s", botOwner, newBranchName),
		"base":  upstreamBranch,
	}
	respBody, status, err = sendReq("POST", prUrl, prPayload)
	if err != nil || status != http.StatusCreated {
		return 0, "", fmt.Errorf("failed to open pull request (%d): %s", status, string(respBody))
	}

	var prResponse struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.Unmarshal(respBody, &prResponse); err != nil {
		return 0, "", err
	}

	return prResponse.Number, prResponse.HTMLURL, nil
}
