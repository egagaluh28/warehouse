package unit

import (
	"testing"
	"warehouse-api/utils"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	t.Run("Success - Generate code with prefix", func(t *testing.T) {
		code := utils.GenerateCode("INV")
		
		assert.NotEmpty(t, code)
		assert.Contains(t, code, "INV-")
		assert.Greater(t, len(code), 10) // Should have prefix + timestamp + random
	})

	t.Run("Success - Different prefixes", func(t *testing.T) {
		code1 := utils.GenerateCode("FK")
		code2 := utils.GenerateCode("INV")

		assert.Contains(t, code1, "FK-")
		assert.Contains(t, code2, "INV-")
		assert.NotEqual(t, code1, code2)
	})

	t.Run("Success - Generate unique codes", func(t *testing.T) {
		code1 := utils.GenerateCode("TEST")
		code2 := utils.GenerateCode("TEST")

		// Should be different due to random component
		assert.NotEqual(t, code1, code2)
	})

	t.Run("Success - Empty prefix", func(t *testing.T) {
		code := utils.GenerateCode("")
		
		assert.NotEmpty(t, code)
		// Should start with dash
		assert.Contains(t, code, "-")
	})
}

func TestRandomString(t *testing.T) {
	t.Run("Success - Generate random string with specified length", func(t *testing.T) {
		length := 10
		str := utils.RandomString(length)

		assert.NotEmpty(t, str)
		assert.Equal(t, length, len(str))
	})

	t.Run("Success - Different calls produce different strings", func(t *testing.T) {
		str1 := utils.RandomString(10)
		str2 := utils.RandomString(10)

		assert.NotEqual(t, str1, str2)
	})

	t.Run("Success - Zero length", func(t *testing.T) {
		str := utils.RandomString(0)
		
		assert.Empty(t, str)
		assert.Equal(t, 0, len(str))
	})

	t.Run("Success - Large length", func(t *testing.T) {
		length := 100
		str := utils.RandomString(length)

		assert.Equal(t, length, len(str))
	})
}
