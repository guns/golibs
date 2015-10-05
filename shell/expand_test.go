package shell

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandUserDir(t *testing.T) {
	u, err := ExpandUserDir("")
	assert.Equal(t, "", u)
	assert.Nil(t, err)

	u, err = ExpandUserDir("/")
	assert.Equal(t, "/", u)
	assert.Nil(t, err)

	u, err = ExpandUserDir("~")
	assert.Equal(t, os.ExpandEnv("${HOME}"), u)
	assert.Nil(t, err)

	u, err = ExpandUserDir("~/Mail")
	assert.Equal(t, os.ExpandEnv("${HOME}/Mail"), u)
	assert.Nil(t, err)

	// FIXME: Not portable
	u, err = ExpandUserDir("~root/Mail")
	assert.Equal(t, os.ExpandEnv("/root/Mail"), u)
	assert.Nil(t, err)
}
