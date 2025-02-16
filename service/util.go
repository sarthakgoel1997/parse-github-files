package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"parse-github-files/model"
	"strings"
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

func getGitHubFileData(baseUrl string, file string, githubReq *http.Request) (fileData model.FileData, err error) {
	client := &http.Client{}
	newUrl := baseUrl + file

	u, err := url.Parse(newUrl)
	if err != nil {
		err = fmt.Errorf("error generating new URL: %s", err.Error())
		return
	}
	githubReq.URL = u

	resp, err := client.Do(githubReq)
	if err != nil {
		err = fmt.Errorf("error making request to GitHub: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response: %s", err.Error())
		return
	}

	// unmarshal JSON data into commits variable
	err = json.Unmarshal(body, &fileData)
	if err != nil {
		err = fmt.Errorf("error unmarshalling JSON: %s", err.Error())
	}

	return
}

func decodeAndParseBase64Data(encoded string) ([]model.ScanResult, error) {
	// decode Base64
	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %v", err)
	}

	// unmarshal JSON into struct array
	var scanResultsWrappers []model.ScanResultsWrapper
	err = json.Unmarshal(jsonData, &scanResultsWrappers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// extract ScanResults from wrapper
	var scanResults []model.ScanResult
	for _, wrapper := range scanResultsWrappers {
		scanResults = append(scanResults, wrapper.ScanResults)
	}

	return scanResults, nil
}

func convertArrayToJson(req []string) string {
	jsonData, _ := json.Marshal(req)
	return string(jsonData)
}
