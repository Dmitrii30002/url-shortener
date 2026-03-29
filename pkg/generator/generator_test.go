package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate_ShouldReturnCorrectString(t *testing.T) {
	gen := New()

	result := gen.Generate()

	assert.Len(t, result, 10)
	for _, char := range result {
		assert.Contains(t, alphabet, string(char))
	}
}

func TestGenerate_ShouldReturnDifferentValues(t *testing.T) {
	gen := New()

	result1 := gen.Generate()
	result2 := gen.Generate()

	assert.NotEqual(t, result1, result2)
}
