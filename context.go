package main

import "fmt"

// NewContext returns ProfileContext created by the Environmental variables.
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
