package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestGitHubImpl_GetGist_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skip test because godotenv seems that it is unable to read dynamic env file.")
		return
	}
	err := godotenv.Load("testdata/github-actions-env.env")
	if err != nil {
		assert.Fail(t, "environment not prepared")
		return
	}
	envValues := NewEnvValues()
	context, err := envValues.NewContext(ProfileFile("testdata/github-test.yml"))
	if err != nil {
		assert.Fail(t, "context creation failed", err)
		return
	}
	gitHub := context.NewGitHub()
	gist, err := gitHub.GetGist(GistID("d1b910d36d314b77b057ea66fbb65e81"), ProfileName("default"))
	assert.Nil(t, err)
	assert.Equal(t, "https://api.github.com/gists/d1b910d36d314b77b057ea66fbb65e81", gist.URL)
	assert.Equal(t, "mike-neck", gist.Owner.Login)
}

func TestGitHubImpl_GetGist_404(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skip test because godotenv seems that it is unable to read dynamic env file.")
		return
	}
	err := godotenv.Load("testdata/github-actions-env.env")
	if err != nil {
		assert.Fail(t, "environment not prepared")
		return
	}
	envValues := NewEnvValues()
	context, err := envValues.NewContext(ProfileFile("testdata/github-test.yml"))
	if err != nil {
		assert.Fail(t, "context creation failed", err)
		return
	}
	gitHub := context.NewGitHub()
	gist, err := gitHub.GetGist(GistID("https://api.github.com/gists/aaaaaaaaaaaa"), ProfileName("default"))
	assert.NotNil(t, err)
	assert.Nil(t, gist)
}
