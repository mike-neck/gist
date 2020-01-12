package main

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestProfileCommandExecutor_AppendExecutor(t *testing.T) {
	var executor profileCommandExecutor
	profile := Profile{
		Name:  "default",
		Token: "aa00bb11cc22",
		Dir:   "/users/ec2-user/gists/repositories",
	}
	executor = &appendExecutor{Profile: profile}
	var current []Profile
	profiles := executor.Invoke(current)
	assert.Equal(t, profileList{profile}, profiles)
}

func TestProfileCommandExecutor_AppendExecutor_InsertedAtFirst(t *testing.T) {
	var executor profileCommandExecutor
	profile := Profile{
		Name: "default",
		Dir:  "/users/ec2-user/gists/repositories",
	}
	executor = &appendExecutor{Profile: profile}
	current := []Profile{
		{
			Name:  "privates",
			Token: "a0b1c2d3e4f5",
		},
	}
	profiles := executor.Invoke(current)
	assert.Equal(t, 1, len(current))
	assert.Equal(t, 2, len(profiles))
	assert.Equal(t, profile, profiles[0])
}

func TestProfileCommandExecutor_OverrideExecutor(t *testing.T) {
	var executor profileCommandExecutor
	profile := Profile{
		Name: "default",
		Dir:  "/users/ec2-user/gists/repositories",
	}
	executor = &overrideExecutor{Profile: profile}
	current := []Profile{
		{
			Name:  "default",
			Token: "aa00bb11cc22",
			Dir:   "/users/ec2-user/destination",
		},
	}
	profiles := executor.Invoke(current)
	assert.Equal(t, 1, len(profiles))
	assert.Equal(t, profile, profiles[0])
}

func TestProfileCommandExecutor_OverrideExecutor_KeepingAnotherProfile(t *testing.T) {
	var executor profileCommandExecutor
	profile := Profile{
		Name: "default",
		Dir:  "/users/ec2-user/gists/repositories",
	}
	executor = &overrideExecutor{Profile: profile}
	another := Profile{
		Name:  "privates",
		Token: "000ccc111",
	}
	current := []Profile{
		{
			Name:  "default",
			Token: "aa00bb11cc22",
			Dir:   "/users/ec2-user/destination",
		},
		another,
	}
	profiles := executor.Invoke(current)
	assert.Equal(t, 2, len(profiles))
	assert.Equal(t, profileList{profile, another}, profiles)
}

func TestProfileList_WriteTo(t *testing.T) {
	profiles := profileList{
		{
			Name:  "default",
			Token: "00ff11ee22dd",
		},
		{
			Name: "privates",
			Dir:  "/users/ec2-user/items",
		},
	}
	var writer io.Writer
	buffer := new(bytes.Buffer)
	writer = buffer
	err := profiles.saveTo(writer)
	assert.Nil(t, err)
	expected := []byte(`- profile: default
  github_access_token: 00ff11ee22dd
- profile: privates
  destination_dir: /users/ec2-user/items
`)
	assert.Equal(t, expected, buffer.Bytes())
}

func TestProfileList_WriteTo_that_FailsForErrorWriter(t *testing.T) {
	profiles := profileList{
		{
			Name:  "default",
			Token: "00ff11ee22dd",
		},
		{
			Name: "privates",
			Dir:  "/users/ec2-user/items",
		},
	}
	writer := &ErrorWriter{}
	err := profiles.saveTo(writer)
	assert.NotNil(t, err)
}

type ErrorWriter struct {
}

func (*ErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("test writer error")
}

func TestAppendOrOverrideProfilesCommand_Executor_OverrideExecutor(t *testing.T) {
	command := AppendOrOverrideProfilesCommand{
		ProfileName:       "default",
		GitHubAccessToken: "",
		DestinationDir:    "",
	}
	context := ProfileContext{
		EnvValues:   EnvValues{},
		ProfileFile: "/users/ec2-user/.gist.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/users/ec2-user/gist/default",
			},
		},
	}
	var executor profileCommandExecutor
	executor = command.executor(context)
	e, ok := executor.(*overrideExecutor)
	assert.True(t, ok)
	assert.Equal(t, ProfileName("default"), e.Name)
	assert.Equal(t, GitHubAccessToken(""), e.Token)
	assert.Equal(t, DestinationDir(""), e.Dir)
}

func TestAppendOrOverrideProfilesCommand_Executor_AppendExecutor(t *testing.T) {
	command := AppendOrOverrideProfilesCommand{
		ProfileName:       "privates",
		GitHubAccessToken: "",
		DestinationDir:    "",
	}
	context := ProfileContext{
		EnvValues:   EnvValues{},
		ProfileFile: "/users/ec2-user/.gist.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/users/ec2-user/gist/default",
			},
		},
	}
	executor := command.executor(context)
	e, ok := executor.(*appendExecutor)
	assert.True(t, ok)
	assert.Equal(t, ProfileName("privates"), e.Name)
	assert.Equal(t, GitHubAccessToken(""), e.Token)
	assert.Equal(t, DestinationDir(""), e.Dir)
}
