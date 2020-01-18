package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// RepositoryMetadata is metadata for each gist.
type RepositoryMetadata struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	URL     string `json:"url"`
	GitURL  string `json:"git_url"`
	Owner   string `json:"owner"`
	Created int64  `json:"created"`
}

// NewMetadataFromGist converts Gist into metadata.
func NewMetadataFromGist(repositoryName RepositoryName, gist Gist) (*RepositoryMetadata, error) {
	createdAt, err := time.Parse("2006-01-02T15:04:05Z", gist.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &RepositoryMetadata{
		ID:      gist.ID,
		Name:    string(repositoryName),
		URL:     gist.URL,
		GitURL:  gist.GitURL,
		Owner:   gist.Owner.Login,
		Created: createdAt.Unix(),
	}, nil
}

// AppendTo appends RepositoryMetadata to file. If file is not existing, the file will be created.
func (md *RepositoryMetadata) AppendTo(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("RepositoryMetadata_AppendTo_OpenFile: %w", err)
	}
	defer func() { _ = file.Close() }()

	bytes, err := json.Marshal(*md)
	if err != nil {
		return fmt.Errorf("Repositorymetadata_AppendTo_MarshalJson: %w", err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("Repositorymetadata_AppendTo_Write: %w", err)
	}
	_, err = file.WriteString("\n")
	if err != nil {
		return fmt.Errorf("Repositorymetadata_AppendTo_WriteNewLine: %w", err)
	}
	return nil
}
