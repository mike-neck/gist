package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

func NewProfileCommand(name, token, dir *string) Command {
	return &AppendOrOverrideProfilesCommand{
		ProfileName:       ProfileName(*name),
		GitHubAccessToken: GitHubAccessToken(*token),
		DestinationDir:    DestinationDir(*dir),
	}
}

// AppendOrOverrideProfilesCommand represents command profile
type AppendOrOverrideProfilesCommand struct {
	ProfileName
	GitHubAccessToken
	DestinationDir
}

// Run profile command.
func (command *AppendOrOverrideProfilesCommand) Run(ctx ProfileContext) error {
	executor := command.executor(ctx)
	profileLists := executor.Invoke(ctx.CurrentProfiles)

	writeCloser, err := ctx.ProfileFile.NewWriter()
	if err != nil {
		return fmt.Errorf("AppendOrOverrideProfilesCommand_Run_NewWriter: %w", err)
	}
	defer func() { _ = writeCloser.Close() }()

	err = profileLists.saveTo(writeCloser)
	if err != nil {
		return fmt.Errorf("AppendOrOverrideProfilesCommand_Run_SaveTo: %w", err)
	}
	return nil
}

////////
// determine profileCommandExecutor
func (command *AppendOrOverrideProfilesCommand) executor(ctx ProfileContext) profileCommandExecutor {
	profileName := command.ProfileName
	for _, p := range ctx.CurrentProfiles {
		if p.Name == profileName {
			executor := overrideExecutor{Profile{
				Name:  profileName,
				Token: command.GitHubAccessToken,
				Dir:   command.DestinationDir,
			}}
			return &executor
		}
	}
	return &appendExecutor{Profile{
		Name:  profileName,
		Token: command.GitHubAccessToken,
		Dir:   command.DestinationDir,
	}}
}

////////
// profileCommandExecutor
type profileCommandExecutor interface {
	Invoke(currentProfiles []Profile) profileList
}

type profileList []Profile

type appendExecutor struct {
	Profile
}

func (ae *appendExecutor) Invoke(currentProfiles []Profile) profileList {
	profiles := make([]Profile, len(currentProfiles)+1)
	profiles[0] = ae.Profile
	for i, p := range currentProfiles {
		profiles[i+1] = p
	}
	return profiles
}

type overrideExecutor struct {
	Profile
}

func (oe *overrideExecutor) profileName() ProfileName {
	return oe.Profile.Name
}

func (oe *overrideExecutor) Invoke(currentProfiles []Profile) profileList {
	profiles := make([]Profile, len(currentProfiles))
	for i, p := range currentProfiles {
		if p.Name == oe.profileName() {
			profiles[i] = oe.Profile
		} else {
			profiles[i] = p
		}
	}
	return profiles
}

////////
// write profiles
func (pl *profileList) saveTo(writer io.Writer) error {
	bytes, err := yaml.Marshal(*pl)
	if err != nil {
		return fmt.Errorf("profileList_SaveTo_Marshal: %w", err)
	}
	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("profileList_SaveTo_Write: %w", err)
	}
	return nil
}
