package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
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
