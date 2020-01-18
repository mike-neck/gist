package main

import (
	"bufio"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestRepositoryMetadata_AppendTo_NewFile(t *testing.T) {
	newFile := "build/test/new-metadata.jsonl"
	_, err := os.Stat("build/test")
	if err != nil {
		_ = os.MkdirAll("build", 0644)
	}
	metadata := RepositoryMetadata{
		ID:      "1a2bc3d4ef",
		Name:    "test",
		URL:     "https://api.github.com/gists/1a2bc3d4ef",
		GitURL:  "git@github.com/gists/1a2bc3d4ef.git",
		Owner:   "test-user",
		Created: time.Now().Unix(),
	}
	err = metadata.AppendTo(newFile)
	assert.Nil(t, err)
	file, err := os.Open(newFile)
	if err != nil {
		assert.Fail(t, "unexpected error@open", err)
		return
	}
	defer func() { _ = file.Close() }()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		assert.Fail(t, "unexpected error@read", err)
		return
	}
	var repositoryMetadata RepositoryMetadata
	err = json.Unmarshal(bytes, &repositoryMetadata)
	if err != nil {
		assert.Fail(t, "unexpected error@unmarshal", err)
		return
	}
	assert.Equal(t, metadata, repositoryMetadata)
}

func TestRepositoryMetadata_AppendTo_Existing(t *testing.T) {
	existingFile := "build/test/existing-metadata.jsonl"
	existingData := RepositoryMetadata{
		ID:      "1a2bc3d4ef",
		Name:    "test",
		URL:     "https://api.github.com/gists/1a2bc3d4ef",
		GitURL:  "git@github.com/gists/1a2bc3d4ef.git",
		Owner:   "test-user",
		Created: time.Now().Unix(),
	}
	err := prepareExistingMetadataFile(existingFile, existingData)
	if err != nil {
		assert.Fail(t, "unexpected error@prepare file", err)
		return
	}

	metadata := RepositoryMetadata{
		ID:      "1100aaccb2",
		Name:    "new-repository",
		URL:     "https://api.github.com/gists/1100aaccb2",
		GitURL:  "git@github.com/gists/1100aaccb2.git",
		Owner:   "new-user",
		Created: time.Now().Unix(),
	}

	err = metadata.AppendTo(existingFile)
	assert.Nil(t, err)

	items, err := readExistingMetadataFile(existingFile)
	if err != nil {
		assert.Fail(t, "unexpected error@read file", err)
		return
	}
	assert.Equal(t, 2, len(items))
	assert.Equal(t, existingData, items[0])
	assert.Equal(t, metadata, items[1])
}

func prepareExistingMetadataFile(existingFile string, metadata RepositoryMetadata) error {
	_, err := os.Stat("build/test")
	if err != nil {
		_ = os.MkdirAll("build", 0644)
	}
	f, err := os.OpenFile(existingFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

func readExistingMetadataFile(existingFile string) ([]RepositoryMetadata, error) {
	f, err := os.Open(existingFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	items := make([]RepositoryMetadata, 0)
	for scanner.Scan() {
		line := scanner.Text()
		bytes := []byte(line)
		var md RepositoryMetadata
		err := json.Unmarshal(bytes, &md)
		if err != nil {
			return nil, err
		}
		items = append(items, md)
	}
	return items, nil
}
