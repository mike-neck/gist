package main

import "fmt"

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

// RepositoryName is name for gist directory.
type RepositoryName string

// CloneCommand clones gist repository
type CloneCommand struct {
	GistID
	ProfileName
	PreferSSH
	RepositoryName
}

// Run command of CloneCommand
func (cc *CloneCommand) Run(ctx ProfileContext) error {
	// determine destination dir
	destinationDir, err := ctx.Dir(cc.ProfileName)
	if err != nil {
		return fmt.Errorf("CloneCommand_Run_ProfileContext_Dir: %w", err)
	}
	fmt.Println(destinationDir)
	// resolve destination dir
	// determine url
	// execute git clone
	// get info on gist
	// write info into repository file under destination dir
	return nil
}

//git@gist.github.com:0674f0f942295225275c349abbe06675.git
//https://gist.github.com/0674f0f942295225275c349abbe06675.git
