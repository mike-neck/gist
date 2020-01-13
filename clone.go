package main

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"os"
	"strings"
)

// GistID is an id of a gist.
type GistID string

// PreferSSH determines protocol (https/ssh).
type PreferSSH int

const (
	https PreferSSH = iota
	ssh
)

// PreferSSHFromBool converts bool to PreferSSH
func PreferSSHFromBool(sshFlag bool) PreferSSH {
	if sshFlag {
		return ssh
	}
	return https
}

// String for PreferSSH
func (s *PreferSSH) String() string {
	switch *s {
	case https:
		return "https"
	case ssh:
		return "ssh"
	}
	panic(fmt.Sprintf("unknown ssh value: %d", s))
}

// BaseURL is gist base url
func (s *PreferSSH) BaseURL() string {
	switch *s {
	case https:
		return "https://gist.github.com/"
	case ssh:
		return "git@gist.github.com:"
	}
	panic(fmt.Sprintf("unknown ssh value: %d", s))
}

// RepositoryName is name for gist directory.
type RepositoryName string

// CloneCommand clones gist repository
type CloneCommand struct {
	GistID
	ProfileName
	PreferSSH
	RepositoryName
}

// DirName is directory name for clone command.
func (cc *CloneCommand) DirName() string {
	if cc.RepositoryName == "" {
		return string(cc.GistID)
	}
	return string(cc.RepositoryName)
}

// Run command of CloneCommand
func (cc *CloneCommand) Run(ctx ProfileContext) error {
	// determine destination dir
	destinationDir, err := ctx.Dir(cc.ProfileName)
	if err != nil {
		return fmt.Errorf("CloneCommand_Run_ProfileContext_Dir: %w", err)
	}
	// resolve destination dir
	targetDirectory, err := destinationDir.Resolve(cc.DirName())
	if err != nil {
		return fmt.Errorf("CloneCommand_Run_Resolve: %w", err)
	}
	fmt.Println(targetDirectory)
	// test directory is empty or not existing
	result, err := testDestinationDir(targetDirectory)
	if err != nil {
		return fmt.Errorf("CloneCommand_Run_TestDir: %w", err)
	}
	if result != resultEmptyDir && result != resultNotExistingDir {
		return fmt.Errorf("CloneCommand_Run_TestDir: %s is not empty directory", targetDirectory)
	} else if result == resultNotExistingDir {
		err := createParentDirectory(targetDirectory)
		if err != nil {
			return fmt.Errorf("CloneCommand_Run_CreateParentDir(%s): %w", targetDirectory, err)
		}
	}
	// determine url
	err = cc.Clone(targetDirectory)
	if err != nil {
		return fmt.Errorf("CloneCommand_Run_Clone: %w", err)
	}
	// execute git clone
	// get info on gist
	// write info into repository file under destination dir
	return nil
}

type testDirResult int

const (
	resultNotExistingDir testDirResult = iota
	resultEmptyDir
	resultHasContentsDir
	resultError
)

func testDestinationDir(directory string) (testDirResult, error) {
	file, err := os.Open(directory)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return resultNotExistingDir, nil
		}
		return resultError, err
	}
	defer func() { _ = file.Close() }()
	_, err = file.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return resultEmptyDir, nil
	}
	if err != nil {
		return resultError, err
	}
	return resultHasContentsDir, nil
}

func createParentDirectory(directory string) error {
	index := strings.LastIndex(directory, "/")
	if index < 0 {
		return nil
	} else if index == 0 {
		return os.Mkdir(directory, 0755)
	}
	dir := directory[:index]
	return os.MkdirAll(dir, 0755)
}

// URL is gist git url
func (cc *CloneCommand) URL() string {
	return fmt.Sprintf("%s%s.git", cc.BaseURL(), cc.GistID)
}

// Clone clones gist repository.
func (cc *CloneCommand) Clone(directory string) error {
	url := cc.URL()
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return fmt.Errorf("GitClone(%s): %w", url, err)
	}
	return nil
}

//git@gist.github.com:0674f0f942295225275c349abbe06675.git
//https://gist.github.com/0674f0f942295225275c349abbe06675.git
