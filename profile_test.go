package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestDefaultProfileFile(t *testing.T) {
	profileFile := DefaultProfileFile()
	assert.True(t, len(string(profileFile)) > 0)
	assert.True(t, strings.HasSuffix(string(profileFile), ".gist.yml"))
}

func TestProfileFile_LoadProfiles(t *testing.T) {
	profileFile := ProfileFile("testdata/profile.yml")
	var profiles []Profile
	profiles, err := profileFile.LoadProfiles()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(profiles))
	profileNames := make([]ProfileName, 2)
	for i, p := range profiles {
		profileNames[i] = p.Name
	}
	assert.Equal(t, []ProfileName{"default", "privates"}, profileNames)
}

func TestProfileFile_LoadProfiles_on_FileNotFound(t *testing.T) {
	profileFile := ProfileFile("testdata/not-exist.yml")
	profiles, err := profileFile.LoadProfiles()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(profiles))
}

func TestProfileFile_LoadProfiles_that_FailsForLoadingDirectory(t *testing.T) {
	profileFile := ProfileFile("testdata")
	_, err := profileFile.LoadProfiles()
	assert.NotNil(t, err)
}

func TestLoadFromReader(t *testing.T) {
	ps, err := LoadFromReader(validReader)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 2, len(ps))

	firstExpected := Profile{
		Name: "default",
		Dir:  "my-gists",
	}
	assert.Equal(t, firstExpected, ps[0])

	secondExpected := Profile{
		Name:  "privates",
		Token: "5f4e3d2c1b0a",
	}
	assert.Equal(t, secondExpected, ps[1])
}

var validReader io.Reader = strings.NewReader(`
- profile: default
  destination_dir: my-gists
  # GITHUB_ACCESS_TOKEN will be used for the profile "default".
- profile: privates
  github_access_token: 5f4e3d2c1b0a
  # $HOME/gist/privates will be used for the profile "privates".
`)

func TestLoadFromReader_that_FailsForNotSlice(t *testing.T) {
	_, err := LoadFromReader(invalidReader)
	assert.NotNil(t, err)
}

var invalidReader io.Reader = strings.NewReader(`
profile: default
github_access_token: aa0bb1cc2
destination_dir: dest
`)

func TestLoadFromReader_that_FailsForInvalidType(t *testing.T) {
	_, err := LoadFromReader(invalidFormatReader)
	assert.NotNil(t, err)
}

var invalidFormatReader io.Reader = strings.NewReader(`
- name: default
  title: my-gists
  rules:
    - name: rule1
      spec: do something
    - name: rule2
      spec: do another
`)

func TestLoadFromReader_that_FailsForEof(t *testing.T) {
	_, err := LoadFromReader(&ErrReader{})
	assert.NotNil(t, err)
}

type ErrReader struct{}

func (*ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestLoadProfileFromFile(t *testing.T) {
	profiles, err := LoadProfileFromFile("testdata/profile.yml")
	if err != nil {
		t.Fail()
	}
	firstExpected := Profile{
		Name: "default",
		Dir:  "my-gists",
	}
	secondExpected := Profile{
		Name:  "privates",
		Token: "5f4e3d2c1b0a",
	}
	assert.Equal(t, []Profile{firstExpected, secondExpected}, profiles)
}

func TestLoadProfileFromFile_that_FailsForNotExisting(t *testing.T) {
	_, err := LoadProfileFromFile("testdata/not-existing.yaml")
	assert.NotNil(t, err)
}
