package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEnvironmentVars(t *testing.T) {
	vars := NewEnvironmentVars()

	assert.Equal(t, EnvironmentVars{
		Port:            "1234",
		WebDir:          "web",
		SiriusURL:       "http://localhost:8080",
		SiriusPublicURL: "",
		Prefix:          "",
	}, vars)
}
