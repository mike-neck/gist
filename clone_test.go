package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPreferSSHFromBool(t *testing.T) {
	preferHTTPS := PreferSSHFromBool(false)
	assert.Equal(t, https, preferHTTPS)
	preferSSH := PreferSSHFromBool(true)
	assert.Equal(t, ssh, preferSSH)
}
