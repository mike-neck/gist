package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

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
  github_access_token: ""
  destination_dir: my-gists
  # GITHUB_ACCESS_TOKEN will be used for the profile "default".
- profile: privates
  github_access_token: 5f4e3d2c1b0a
  destination_dir: ""
  # $HOME/gist/privates will be used for the profile "privates".
`)

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
