package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

// LoadProfiles loads profile from ProfileFile
func (file *ProfileFile) LoadProfiles() ([]Profile, error) {
	fileName := string(*file)
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return []Profile{}, nil
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return nil, fmt.Errorf("%s is directory", fileName)
	}
	profiles, err := LoadProfileFromFile(fileName)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

// DefaultProfileFile returns default profile file.
func DefaultProfileFile() ProfileFile {
	home := os.Getenv("HOME")
	file := fmt.Sprintf("%s/.gist.yml", home)
	return ProfileFile(file)
}

// ProfileYaml is the Profile data structure.
// This is raw type of Profile that is not validated.
type ProfileYaml struct {
	Name  ProfileName       `yaml:"profile"`
	Token GitHubAccessToken `yaml:"github_access_token,omitempty"`
	Dir   DestinationDir    `yaml:"destination_dir,omitempty"`
}

// Profile is validated ProfileYaml.
type Profile ProfileYaml

// LoadProfileFromFile loads profiles from a given file.
func LoadProfileFromFile(file string) ([]Profile, error) {
	reader, err := os.Open(file)
	if err != nil {
		w := fmt.Errorf("LoadProfile_Open: %w", err)
		return nil, w
	}
	defer func() {
		_ = reader.Close()
	}()
	return LoadFromReader(reader)
}

// LoadFromReader loads profiles from given reader.
func LoadFromReader(reader io.Reader) ([]Profile, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		w := fmt.Errorf("LoadProfile_Read: %w", err)
		return nil, w
	}

	profileYamls := make([]ProfileYaml, 0)
	err = yaml.Unmarshal(bytes, &profileYamls)
	if err != nil {
		w := fmt.Errorf("LoadProfile_Unmarshal: %w", err)
		return nil, w
	}

	profiles := make([]Profile, 0)
	for _, p := range profileYamls {
		if p.Name != "" {
			profiles = append(profiles, Profile(p))
		} else {
			return nil, errors.New("invalid YAML format")
		}
	}

	return profiles, nil
}

//////////////

// AppendOrOverrideProfilesCommand represents command profile
type AppendOrOverrideProfilesCommand struct {
	ProfileName
	GitHubAccessToken
	DestinationDir
}

// Run profile command.
func (command *AppendOrOverrideProfilesCommand) Run(ctx ProfileContext) error {
	// check current profiles and determine profileCommandExecutor
	// invoke profileCommandExecutor
	// write profiles to configuration file
	return nil
}

////////
// determine profileCommandExecutor

////////
// profileCommandExecutor
type profileCommandExecutor interface {
	Invoke(currentProfiles []Profile) []Profile
}

type appendExecutor struct {
	Profile
}

func (ae *appendExecutor) Invoke(currentProfiles []Profile) []Profile {
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

func (oe *overrideExecutor) Invoke(currentProfiles []Profile) []Profile {
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
