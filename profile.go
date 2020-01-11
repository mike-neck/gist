package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type ProfileYaml struct {
	Name  string `yaml:"profile"`
	Token string `yaml:"github_access_token,omitempty"`
	Dir   string `yaml:"destination_dir,omitempty"`
}

// Profile is validated ProfileYaml
type Profile ProfileYaml

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
