package main

// GistID is an id of a gist.
type GistID string

// UserName is owner of gist.
type UserName string

// CloneCommand clones gist repository
type CloneCommand struct {
	GistID
	ProfileName
	UserName
}

// Run command of CloneCommand
func (cc *CloneCommand) Run(ctx ProfileContext) error {
	return nil
}
