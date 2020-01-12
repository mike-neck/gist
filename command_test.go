package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfileContext_Token(t *testing.T) {
	context := ProfileContext{
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("privates")
	var githubAccessToken GitHubAccessToken
	githubAccessToken, err := context.Token(profileName)
	assert.Nil(t, err)
	assert.Equal(t, GitHubAccessToken("a0b1c2d3e4f5"), githubAccessToken)
}

func TestProfileContext_Token_failure(t *testing.T) {
	context := ProfileContext{
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("app")
	_, err := context.Token(profileName)
	assert.NotNil(t, err)
}

func TestProfileContext_Token_From_Env(t *testing.T) {
	context := ProfileContext{
		EnvValues: EnvValues{
			GitHubAccessToken: "aa00bb11cc22",
		},
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("default")
	githubAccessToken, err := context.Token(profileName)
	assert.Nil(t, err)
	assert.Equal(t, GitHubAccessToken("aa00bb11cc22"), githubAccessToken)
}

func TestProfileContext_Dir(t *testing.T) {
	context := ProfileContext{
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("default")
	var destinationDir DestinationDir
	destinationDir, err := context.Dir(profileName)
	assert.Nil(t, err)
	assert.Equal(t, DestinationDir("/user/name/gists"), destinationDir)
}

func TestProfileContext_Dir_failure(t *testing.T) {
	context := ProfileContext{
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("app")
	_, err := context.Dir(profileName)
	assert.NotNil(t, err)
}

func TestProfileContext_Dir_From_Env(t *testing.T) {
	context := ProfileContext{
		EnvValues: EnvValues{
			UserHome: "/users/ec2-user",
		},
		ProfileFile: "test.yml",
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/user/name/gists",
			},
			{
				Name:  "privates",
				Token: "a0b1c2d3e4f5",
			},
		},
	}
	profileName := ProfileName("privates")
	destinationDir, err := context.Dir(profileName)
	assert.Nil(t, err)
	assert.Equal(t, DestinationDir("/users/ec2-user/gist/privates"), destinationDir)
}

func TestNewEnvValues(t *testing.T) {
	_ = godotenv.Load("testdata/test.env")
	expected := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
	}
	envValues := NewEnvValues()
	assert.Equal(t, expected.GitHubAccessToken, envValues.GitHubAccessToken)
	assert.NotEmpty(t, envValues.UserHome)
}
