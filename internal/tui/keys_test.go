package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyMap_ShortHelp(t *testing.T) {
	shortHelp := keys.ShortHelp()

	assert.Len(t, shortHelp, 2)
	assert.Equal(t, keys.Help, shortHelp[0])
	assert.Equal(t, keys.Quit, shortHelp[1])
}

func TestKeyMap_FullHelp(t *testing.T) {
	fullHelp := keys.FullHelp()

	assert.Len(t, fullHelp, 4)

	assert.Len(t, fullHelp[0], 4)
	assert.Contains(t, fullHelp[0], keys.Up)
	assert.Contains(t, fullHelp[0], keys.Down)
	assert.Contains(t, fullHelp[0], keys.Left)
	assert.Contains(t, fullHelp[0], keys.Right)

	assert.Len(t, fullHelp[1], 4)
	assert.Contains(t, fullHelp[1], keys.InProgress)
	assert.Contains(t, fullHelp[1], keys.ToDo)
	assert.Contains(t, fullHelp[1], keys.Done)
	assert.Contains(t, fullHelp[1], keys.Cancelled)

	assert.Len(t, fullHelp[2], 4)
	assert.Contains(t, fullHelp[2], keys.Show)
	assert.Contains(t, fullHelp[2], keys.Edit)
	assert.Contains(t, fullHelp[2], keys.Add)
	assert.Contains(t, fullHelp[2], keys.Remove)

	assert.Len(t, fullHelp[3], 2)
	assert.Contains(t, fullHelp[3], keys.Help)
	assert.Contains(t, fullHelp[3], keys.Quit)
}
