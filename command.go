package main

import "fmt"

// ProfileFile is the file of profiles to be loaded.
type ProfileFile string

// ProfileName is name of profile.
type ProfileName string

// GitHubAccessToken is access token of github.com with gist scope.
type GitHubAccessToken string

// DestinationDir is destination directory where to clone gist repositories.
type DestinationDir string

// UserHome is home path.
type UserHome string

// EnvValues is environmental values.
type EnvValues struct {
	GitHubAccessToken
	UserHome
}

// ProfileContext is Profiles and ProfileFile which the command will be executed on.
type ProfileContext struct {
	EnvValues
	ProfileFile
	CurrentProfiles []Profile
}

// Token returns github access token for given profile.
func (context *ProfileContext) Token(profileName ProfileName) (GitHubAccessToken, error) {
	for _, profile := range context.CurrentProfiles {
		if profile.Name == profileName {
			token := profile.Token
			if token == "" {
				token = context.EnvValues.GitHubAccessToken
			}
			return token, nil
		}
	}
	return "", fmt.Errorf("no profile found(name = %s)", profileName)
}

// Command represents command being executed by user.
type Command interface {
	Run(ctx *ProfileContext) error
}
