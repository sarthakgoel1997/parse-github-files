package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"parse-github-files/model"
	"strings"
	"time"
)

func extractOwnerRepo(url string) (string, string, error) {
	// remove the base URL prefix
	trimmed := strings.TrimPrefix(url, "https://github.com/")

	// split by "/"
	parts := strings.Split(trimmed, "/")

	//ensure we have both owner and repo
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL: %s", url)
	}

	return parts[0], parts[1], nil
}

func prepareGitHubAPIRequest(repository string) (githubReq *http.Request, baseUrl string, err error) {
	owner, repo, err := extractOwnerRepo(repository)
	if err != nil {
		return
	}

	baseUrl = fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/", owner, repo)
	method := "GET"
	githubReq, err = http.NewRequest(method, baseUrl, nil)
	if err != nil {
		err = fmt.Errorf("error generating new request: %s", err.Error())
		return
	}
	addhttpAuthRequestHeaders(githubReq)

	return
}

func addhttpAuthRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PERSONAL_ACCESS_TOKEN"))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
}

func getDataFromGitHub(baseUrl string, file string, githubReq *http.Request) (fileData model.FileData, err error) {
	client := &http.Client{}
	newUrl := baseUrl + file

	u, err := url.Parse(newUrl)
	if err != nil {
		err = fmt.Errorf("error generating new URL: %s", err.Error())
		return
	}
	githubReq.URL = u

	maxRetries := 2
	waitTime := 1 * time.Second // initial backoff time

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := client.Do(githubReq)
		if err != nil {
			if attempt < maxRetries {
				fmt.Printf("%s File: Request failed, retrying... (attempt %d/%d)\n", newUrl, attempt+1, maxRetries)
				time.Sleep(waitTime)
				waitTime *= 2 // exponential backoff
				continue
			}
			return fileData, fmt.Errorf("error making request to GitHub: %w", err)
		}
		if resp.StatusCode >= 400 {
			if attempt < maxRetries {
				fmt.Printf("%s File: Request failed with status %d, retrying... (attempt %d/%d)\n", newUrl, resp.StatusCode, attempt+1, maxRetries)
				time.Sleep(waitTime)
				waitTime *= 2
				continue
			}
			return fileData, fmt.Errorf("GitHub API request failed with status code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		// read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fileData, fmt.Errorf("error reading response: %w", err)
		}

		// unmarshal JSON data
		err = json.Unmarshal(body, &fileData)
		if err != nil {
			return fileData, fmt.Errorf("error unmarshalling JSON: %w", err)
		}

		return fileData, nil
	}

	return fileData, fmt.Errorf("failed to fetch GitHub API data after %d attempts", maxRetries+1)
}
