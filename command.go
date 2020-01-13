package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Application interface {
	Start() error
}

// ProfileFile is the file of profiles to be loaded.
type ProfileFile string

//NewWriter creates writer of ProfileFile
func (file *ProfileFile) NewWriter() (io.WriteCloser, error) {
	path := string(*file)
	lastIndex := strings.LastIndex(path, "/")
	if lastIndex > 0 {
		parent := path[0:lastIndex]
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return nil, fmt.Errorf("ProfileFile_NewWriter_MkdirAll(%s): %w", parent, err)
		}
	}
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("ProfileFile_NewWriter_OpenFile: %w", err)
	}
	return writer, nil
}

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

// NewEnvValues loads from environmental variables.
func NewEnvValues() EnvValues {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	userHome := os.Getenv("HOME")
	if userHome == "" {
		userHome = os.Getenv("HOMEPATH")
	}
	return EnvValues{
		GitHubAccessToken: GitHubAccessToken(githubAccessToken),
		UserHome:          UserHome(userHome),
	}
}

func (ev *EnvValues) DefaultProfileFile() ProfileFile {
	return ProfileFile(fmt.Sprintf("%s/.gist.yml", ev.UserHome))
}

func (ev *EnvValues) NewContext(file ProfileFile) (ProfileContext, error) {
	profileFile := file
	if profileFile == "" {
		profileFile = ev.DefaultProfileFile()
	}
	profiles, err := profileFile.LoadProfiles()
	if err != nil {
		return ProfileContext{}, fmt.Errorf("EnvValues_NewContext_LoadProfiles: %w", err)
	}
	return ProfileContext{
		EnvValues:       *ev,
		ProfileFile:     profileFile,
		CurrentProfiles: profiles,
	}, nil
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

// Dir returns DestinationDir of given profile.
func (context *ProfileContext) Dir(profileName ProfileName) (DestinationDir, error) {
	for _, profile := range context.CurrentProfiles {
		if profile.Name == profileName {
			dir := profile.Dir
			if dir == "" {
				dir = DestinationDir(fmt.Sprintf("%s/gist/%s", context.EnvValues.UserHome, profileName))
			}
			return dir, nil
		}
	}
	return "", fmt.Errorf("no profile found(name = %s)", profileName)
}

// Command represents command being executed by user.
type Command interface {
	// Run executes each command.
	Run(ctx ProfileContext) error
}
