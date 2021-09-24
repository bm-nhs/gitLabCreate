package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreatePullRequestPayload struct {
	Title string `json:"title"`
	Head string `json:"head"`
	Base string `json:"base"`
}

// PullRequest creates a GitHub pull request taking the CreatePullRequestPayload and the repo owner and repository name
func PullRequest(payload CreatePullRequestPayload, owner string, repository string) error {
	payloadBytes , err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repository)
	req, err := http.NewRequest("Post", url, body)
	req.Header.Set("Content-Type", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)
	return err
}
