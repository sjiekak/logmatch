package main

import (
	"testing"

	"github.com/sjiekak/logen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logen(t *testing.T) {
	st, err := logen.NewSanitizer()
	require.NoError(t, err)
	assert.NotNil(t, st)
}
