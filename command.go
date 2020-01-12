package main

// ProfileName is name of profile.
type ProfileName string

// GitHubAccessToken is access token of github.com with gist scope.
type GitHubAccessToken string

// DestinationDir is destination directory where to clone gist repositories.
type DestinationDir string

// ProfileContext is Profiles and ProfileFile which the command will be executed on.
type ProfileContext struct {
	ProfileFile
	CurrentProfiles []Profile
}

// Command represents command being executed by user.
type Command interface {
	Run(ctx *ProfileContext) error
}
