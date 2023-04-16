package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func EnvironmentCommit(repository, environment, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%v/deployments?environment=%v&per_page=1", repository, environment)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request to GitHub: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get response from GitHub: %w", err)
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read GitHub response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("request to GitHub failed with status code %v: %v", res.Status, string(rawBody))
	}

	resBody := []environmentResponseModel{}
	if err := json.Unmarshal(rawBody, &resBody); err != nil {
		return "", fmt.Errorf("failed to decode GitHub environment response body: %w", err)
	}

	// When there are no deployments, then we do a git diff on ..{END_COMMIT},
	// which should give all changes.
	if len(resBody) < 1 {
		return "", nil
	}

	return resBody[0].Sha, nil
}

type environmentResponseModel struct {
	Sha string `json:"sha"`
}
