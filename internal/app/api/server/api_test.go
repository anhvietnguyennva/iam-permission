package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPublicAPIServer_Successfully(t *testing.T) {
	name := "TestBuildPublicAPIServer_Successfully"
	t.Log(name)

	server := NewAPIServer(false)
	assert.NotNil(t, server)
}

func TestBuildAdminAPIServer_Successfully(t *testing.T) {
	name := "TestBuildAdminAPIServer_Successfully"
	t.Log(name)

	server := NewAPIServer(true)
	assert.NotNil(t, server)
}
