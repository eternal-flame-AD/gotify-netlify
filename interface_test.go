package main

import (
	"testing"

	plugin "github.com/gotify/plugin-api"
	"github.com/stretchr/testify/assert"
)

func TestAPICompatibility(t *testing.T) {
	assert.Implements(t, (*plugin.Plugin)(nil), new(Plugin))
	assert.Implements(t, (*plugin.Webhooker)(nil), new(Plugin))
	assert.Implements(t, (*plugin.Messenger)(nil), new(Plugin))
	assert.Implements(t, (*plugin.Displayer)(nil), new(Plugin))
	assert.Implements(t, (*plugin.Configurer)(nil), new(Plugin))
}
