package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"os"
	"testing"
)

func TestPreferSSHFromBool(t *testing.T) {
	preferHTTPS := PreferSSHFromBool(false)
	assert.Equal(t, https, preferHTTPS)
	preferSSH := PreferSSHFromBool(true)
	assert.Equal(t, ssh, preferSSH)
}

// example of go-git
func GitCloneExample(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		assert.Fail(t, "failed to get temp dir: %s", err.Error())
	}
	fmt.Println(tempDir)
	err = os.RemoveAll(tempDir)
	if err != nil {
		assert.Fail(t, "remove dir contents: %s", err.Error())
		return
	}
	repository, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: "git@gist.github.com:0674f0f942295225275c349abbe06675.git",
	})
	if err != nil {
		assert.Fail(t, "failed to clone, %s", err.Error())
		return
	}
	branches, err := repository.Branches()
	if err != nil {
		assert.Fail(t, "failed to get branches: %s", err.Error())
		return
	}
	err = branches.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println("branch:", reference.Name(), reference.String(), reference.Hash())
		return nil
	})
	if err != nil {
		assert.Fail(t, "failed to show branches: %s", err.Error())
	}
}

func TestCloneCommand_URL_HTTPS1(t *testing.T) {
	command := CloneCommand{
		GistID:         "11aa22bb33cc",
		ProfileName:    "default",
		PreferSSH:      https,
		RepositoryName: "",
	}
	var url string
	url = command.URL()
	assert.Equal(t, "https://gist.github.com/11aa22bb33cc.git", url)
}

func TestCloneCommand_URL_HTTPS2(t *testing.T) {
	command := CloneCommand{
		GistID:         "10a29b38c",
		ProfileName:    "default",
		PreferSSH:      https,
		RepositoryName: "",
	}
	var url string
	url = command.URL()
	assert.Equal(t, "https://gist.github.com/10a29b38c.git", url)
}

func TestCloneCommand_URL_SSH(t *testing.T) {
	command := CloneCommand{
		GistID:         "11aa22bb33cc",
		ProfileName:    "default",
		PreferSSH:      ssh,
		RepositoryName: "",
	}
	var url string
	url = command.URL()
	assert.Equal(t, "git@gist.github.com:11aa22bb33cc.git", url)
}

func TestCloneCommand_Clone(t *testing.T) {
	command := CloneCommand{
		GistID:         "d1b910d36d314b77b057ea66fbb65e81",
		ProfileName:    "default",
		PreferSSH:      https,
		RepositoryName: "gist-example",
	}
	err := command.Clone("build/clone/test/gist-example")
	assert.Nil(t, err)
	dir, err := os.Open("build/clone/test/gist-example")
	assert.Nil(t, err)
	names, err := dir.Readdirnames(30)
	assert.Nil(t, err)
	assert.True(t, 0 <= indexOf("test.go", names), "test.go", names)
	assert.True(t, 0 <= indexOf("test.md", names), "test.md", names)
	assert.True(t, 0 <= indexOf(".git", names), "git", names)
}

func indexOf(item string, items []string) int {
	for idx, i := range items {
		if item == i {
			return idx
		}
	}
	return -1
}

func TestNewGistID_Success(t *testing.T) {
	var gistID *GistID
	gistID, err := NewGistID("12a34b56c78d90ef")
	assert.Nil(t, err)
	assert.Equal(t, GistID("12a34b56c78d90ef"), *gistID)
}

func TestNewGistID_InvalidID(t *testing.T) {
	gistID, err := NewGistID("12a34b56c78d90efg")
	assert.NotNil(t, err)
	assert.Nil(t, gistID)
}

func TestNewGistID_InvalidID2(t *testing.T) {
	gistID, err := NewGistID("1a2c-9b8e")
	assert.NotNil(t, err)
	assert.Nil(t, gistID)
}
