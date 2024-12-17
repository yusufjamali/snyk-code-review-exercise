package dependencyresolver

import (
	"errors"
	"fmt"
	"github.com/snyk/snyk-code-review-exercise/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_DoWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		mPkgClient := newMockPkgClient(t)
		mPkgClient.Set("root", "1.0.0", &models.PackageClientResponse{Name: "root", Version: "1.0.0", Dependencies: map[string]string{"dep_1": "1.0.0", "dep_2": "1.0.0"}})
		mPkgClient.Set("dep_1", "1.0.0", &models.PackageClientResponse{Name: "dep_1", Version: "1.0.0"})
		mPkgClient.Set("dep_2", "1.0.0", &models.PackageClientResponse{Name: "dep_2", Version: "1.0.0", Dependencies: map[string]string{"dep_3": "1.0.0"}})
		mPkgClient.Set("dep_3", "1.0.0", &models.PackageClientResponse{Name: "dep_3", Version: "1.0.0"})

		// act
		sut := New(mPkgClient)
		got, gotErr := sut.DoWork("root", "1.0.0")

		// assert
		want := &models.PackageVersion{
			Name:    "root",
			Version: "1.0.0",
			Dependencies: map[string]*models.PackageVersion{
				"dep_1": &models.PackageVersion{
					Name:         "dep_1",
					Version:      "1.0.0",
					Dependencies: map[string]*models.PackageVersion{},
				},
				"dep_2": &models.PackageVersion{
					Name:    "dep_2",
					Version: "1.0.0",
					Dependencies: map[string]*models.PackageVersion{
						"dep_3": &models.PackageVersion{
							Name:         "dep_3",
							Version:      "1.0.0",
							Dependencies: map[string]*models.PackageVersion{},
						},
					},
				},
			},
		}

		assert.Nil(t, gotErr)
		assert.Equal(t, want, got)
	})

	t.Run("error", func(t *testing.T) {
		mPkgClient := newMockPkgClient(t)
		mPkgClient.SetErr("root", "1.0.0", errors.New("test error"))

		// act
		sut := New(mPkgClient)
		got, gotErr := sut.DoWork("root", "1.0.0")

		// assert
		want := errors.New("test error")
		assert.Nil(t, got)
		assert.Equal(t, gotErr, want)
	})
}

type mockPkgClient struct {
	t                 *testing.T
	invocation        map[string]*models.PackageClientResponse
	invocationWithErr map[string]error
}

func newMockPkgClient(t *testing.T) *mockPkgClient {
	t.Helper()
	invocation := make(map[string]*models.PackageClientResponse, 0)
	invocationWithErr := make(map[string]error, 0)
	return &mockPkgClient{t, invocation, invocationWithErr}
}

func (m *mockPkgClient) Get(name string, version string) (*models.PackageClientResponse, error) {
	if v, ok := m.invocation[m.key(name, version)]; ok {
		return v, nil
	}
	if v, ok := m.invocationWithErr[m.key(name, version)]; ok {
		return nil, v
	}
	m.t.Errorf("no mock set up for: %s and %s", name, version)
	return nil, nil
}

func (m *mockPkgClient) Set(name string, version string, want *models.PackageClientResponse) {
	m.t.Helper()
	m.invocation[m.key(name, version)] = want
}

func (m *mockPkgClient) SetErr(name string, version string, want error) {
	m.t.Helper()
	m.invocationWithErr[m.key(name, version)] = want
}

func (m *mockPkgClient) key(name string, version string) string {
	return fmt.Sprintf("%s_%s", name, version)
}
