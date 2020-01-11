package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type Profile struct {
	Name  string `yaml:"profile"`
	Token string `yaml:"github_access_token,omitempty"`
	Dir   string `yaml:"destination_dir,omitempty"`
}

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

	profiles := make([]Profile, 0)
	err = yaml.Unmarshal(bytes, &profiles)
	if err != nil {
		w := fmt.Errorf("LoadProfile_Unmarshal: %w", err)
		return nil, w
	}

	return profiles, nil
}
