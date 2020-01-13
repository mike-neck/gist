package main

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

// CloneCommand clones gist repository
type CloneCommand struct {
	GistID
	ProfileName
}

// Run command of CloneCommand
func (cc *CloneCommand) Run(ctx ProfileContext) error {
	return nil
}

//git@gist.github.com:0674f0f942295225275c349abbe06675.git
//https://gist.github.com/0674f0f942295225275c349abbe06675.git
