package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitHubImpl_GetGist_404(t *testing.T) {
	err := godotenv.Load("testdata/github-actions-env.env")
	if err != nil {
		assert.Fail(t, "environment not prepared")
		return
	}
	envValues := NewEnvValues()
	assert.Nil(t, envValues, "envValues", envValues)
}
