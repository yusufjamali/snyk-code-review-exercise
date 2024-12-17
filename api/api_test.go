package api_test

import (
	"encoding/json"
	"github.com/snyk/snyk-code-review-exercise/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/snyk/snyk-code-review-exercise/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// comment: this test case might start failing when react/16.13.0 package updates its underlying dependency
// would you consider changing the asserts in a way that makes it not rely on version matches for successful test results?
// OR can use an older version of the package for testing here?

// consider adding tests for other validation scenarios, eg: incorrect input version, unknown package name and version
func TestPackageHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := api.New()
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := server.Client().Get(server.URL + "/package/react/16.13.0")
		require.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.Nil(t, err)

		var data models.PackageVersion
		err = json.Unmarshal(body, &data)
		require.Nil(t, err)

		assert.Equal(t, "react", data.Name)
		assert.Equal(t, "16.13.0", data.Version)

		fixture, err := os.Open(filepath.Join("testdata", "react-16.13.0.json"))
		require.Nil(t, err)
		var fixtureObj models.PackageVersion
		require.Nil(t, json.NewDecoder(fixture).Decode(&fixtureObj))

		assert.Equal(t, fixtureObj, data)
	})

	t.Run("incorrect package", func(t *testing.T) {
		handler := api.New()
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := server.Client().Get(server.URL + "/package/unknown/16.13.0")
		require.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

}
