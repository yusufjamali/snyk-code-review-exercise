package npmpackageclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Get(t *testing.T) {
	t.Run("when equal to version", func(t *testing.T) {
		// arrange
		sut := New("https://registry.npmjs.org")
		got, err := sut.Get("react", "16.13.0")

		// assert
		assert.Nil(t, err)
		assert.NotNil(t, got)
	})
	t.Run("when ^ (care) version", func(t *testing.T) {
		// arrange
		sut := New("https://registry.npmjs.org")
		got, err := sut.Get("prop-types", "^15.6.0")

		// assert
		want := "15.8.1"
		assert.Nil(t, err)
		assert.Equal(t, want, got.Version)
	})
	t.Run("when ~ (care) version", func(t *testing.T) {
		// arrange
		sut := New("https://registry.npmjs.org")
		got, err := sut.Get("prop-types", "~15.6.0")

		// assert
		want := "15.6.2"
		assert.Nil(t, err)
		assert.Equal(t, want, got.Version)
	})

	t.Run("when unknown package", func(t *testing.T) {
		// arrange
		sut := New("https://registry.npmjs.org")
		got, err := sut.Get("unknown", "0.0.1")

		// assert
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})

	t.Run("when unknown version", func(t *testing.T) {
		// arrange
		sut := New("https://registry.npmjs.org")
		got, err := sut.Get("prop-types", "15.6.3")

		// assert
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})
}
