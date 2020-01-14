package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GitHub offers access to github.com
type GitHub interface {
	GetGist(gistID GistID, profileName ProfileName) (*Gist, error)
}

// Gist represents gist API response, some of them are omitted.
type Gist struct {
	URL         string     `json:"url"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"created_at"`
	Owner       GitHubUser `json:"owner"`
}

// GitHubUser is github user.
type GitHubUser struct {
	Login string `json:"login"`
}

var githubAPIBaseURL = "https://api.github.com"
var acceptHeader string = "application/vnd.github.v3+json"

// NewGitHub create github data from current ProfileContext.
func (context *ProfileContext) NewGitHub() GitHub {
	return &gitHubImpl{*context}
}

type gitHubImpl struct {
	ProfileContext
}

func (gh *gitHubImpl) GetGist(gistID GistID, profileName ProfileName) (*Gist, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/gists/%s", githubAPIBaseURL, gistID), nil)
	if err != nil {
		return nil, fmt.Errorf("GitHub_GetGist_NewRequest: %w", err)
	}
	accessToken, err := gh.Token(profileName)
	if err != nil {
		return nil, fmt.Errorf("GitHub_GetGist_Token: %w", err)
	}
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Add("accept", acceptHeader)
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GitHub_GetGist_DoRequest: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	sc := response.StatusCode
	if sc < 200 || 300 <= sc {
		return nil, fmt.Errorf("failed to get gist info(%s, http status:%s)", gistID, response.Status)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("GitHub_GetGist_ReadAll: %w", err)
	}

	var gist Gist
	err = json.Unmarshal(bytes, &gist)
	if err != nil {
		return nil, fmt.Errorf("GitHub_GetGist_JsonUnmarshal: %w", err)
	}

	return &gist, nil
}
