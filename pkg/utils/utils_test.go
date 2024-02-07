package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	testCases := []string{"bad.env", "empty.env", "multipleEquals.env", "nonexistent.env"}

	for _, test := range testCases {
		err := LoadEnv("../../mocks/" + test)

		if assert.Nil(t, err) {
			assert.ErrorContains(t, err, "no such file or directory")
		}
	}
}
