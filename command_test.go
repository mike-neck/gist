package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewEnvValues(t *testing.T) {
	_ = godotenv.Load("testdata/test.env")
	expected := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
	}
	envValues := NewEnvValues()
	assert.Equal(t, expected.GitHubAccessToken, envValues.GitHubAccessToken)
	assert.NotEmpty(t, envValues.UserHome)
}

func TestProfileFile_NewWriter(t *testing.T) {
	profileFile := ProfileFile("build/test/new-file")
	writer, err := profileFile.NewWriter()
	assert.Nil(t, err)
	_ = writer.Close()
	info, err := os.Stat(string(profileFile))
	assert.Nil(t, err)
	assert.False(t, info.IsDir())
}

func TestProfileFile_NewWriter_SameDirectory(t *testing.T) {
	profileFile := ProfileFile("build-test")
	writeCloser, err := profileFile.NewWriter()
	assert.Nil(t, err, "call of NewWriter")
	_ = writeCloser.Close()
	info, err := os.Stat(string(profileFile))
	assert.Nil(t, err, "call of Stat")
	assert.False(t, info.IsDir())
	err = os.Remove(string(profileFile))
	assert.Nil(t, err, "call of Remove")
}

func TestProfileFile_NewWriter_ExistingFile(t *testing.T) {
	prepareExistingFile(t)
	profileFile := ProfileFile("build/existing/conf.yml")
	writeCloser, err := profileFile.NewWriter()
	assert.Nil(t, err)
	_ = writeCloser.Close()
}

func prepareExistingFile(t *testing.T) {
	err := os.MkdirAll("build/existing", 0755)
	if err != nil {
		t.Fail()
	}
	file, err := os.Create("build/existing/conf.yml")
	if err != nil {
		t.Fail()
	}
	defer func() { _ = file.Close() }()
	_, err = file.Write([]byte("- profile: default"))
	if err != nil {
		t.Fail()
	}
}

func TestEnvValues_NewContext(t *testing.T) {
	envValues := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
		UserHome:          "/users/ec2-user",
	}
	context, err := envValues.NewContext("")
	assert.Nil(t, err)
	expected := ProfileContext{
		EnvValues:       envValues,
		ProfileFile:     ProfileFile("/users/ec2-user/.gist.yml"),
		CurrentProfiles: []Profile{},
	}
	assert.Equal(t, expected, context)
}

func TestEnvValues_NewContext_WithUserSpecProfileFile(t *testing.T) {
	envValues := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
		UserHome:          "/users/ec2-user",
	}
	context, err := envValues.NewContext("testdata/test.yml")
	assert.Nil(t, err)
	expected := ProfileContext{
		EnvValues:       envValues,
		ProfileFile:     ProfileFile("testdata/test.yml"),
		CurrentProfiles: []Profile{},
	}
	assert.Equal(t, expected, context)
}

func TestEnvValues_NewContext_WithUserSpecProfileFileWithContents(t *testing.T) {
	envValues := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
		UserHome:          "/users/ec2-user",
	}
	context, err := envValues.NewContext("testdata/profile.yml")
	assert.Nil(t, err)
	expected := ProfileContext{
		EnvValues:   envValues,
		ProfileFile: ProfileFile("testdata/profile.yml"),
		CurrentProfiles: []Profile{
			{
				Name: "default",
				Dir:  "/users/foo/my-gists",
			},
			{
				Name:  "privates",
				Token: "5f4e3d2c1b0a",
			},
		},
	}
	assert.Equal(t, expected, context)
}

func TestEnvValues_NewContext_NotExistingProfileFile(t *testing.T) {
	envValues := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
		UserHome:          "/users/ec2-user",
	}
	_, err := envValues.NewContext("testdata/not-exists.yml")
	assert.Nil(t, err)
}

func TestEnvValues_NewContext_FileIsDirectory(t *testing.T) {
	envValues := EnvValues{
		GitHubAccessToken: "aa00bb11cc22",
		UserHome:          "/users/ec2-user",
	}
	_, err := envValues.NewContext("testdata")
	assert.NotNil(t, err)
}

func TestDestinationDir_Resolve(t *testing.T) {
	destinationDir := DestinationDir("build")
	var dir string
	dir, err := destinationDir.Resolve("sub")
	assert.Nil(t, err)
	assert.Equal(t, "build/sub", dir)
}

func TestDestinationDir_Resolve_WithSuffixSlush(t *testing.T) {
	destinationDir := DestinationDir("build/")
	dir, err := destinationDir.Resolve("sub")
	assert.Nil(t, err)
	assert.Equal(t, "build/sub", dir)
}

func TestDestinationDir_Resolve_WithPrefixOnParam(t *testing.T) {
	destinationDir := DestinationDir("build")
	dir, err := destinationDir.Resolve("/sub")
	assert.Nil(t, err)
	assert.Equal(t, "build/sub", dir)
}

func TestDestinationDir_Resolve_ParamEmpty(t *testing.T) {
	destinationDir := DestinationDir("build")
	_, err := destinationDir.Resolve("")
	assert.NotNil(t, err)
}
