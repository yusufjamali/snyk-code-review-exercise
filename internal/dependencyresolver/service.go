package dependencyresolver

import (
	"github.com/snyk/snyk-code-review-exercise/internal/models"
)

type Service struct {
	pkgClient PackageClient
}

type PackageClient interface {
	Get(name string, version string) (*models.PackageClientResponse, error)
}

func New(pkgClient PackageClient) *Service {
	return &Service{pkgClient: pkgClient}
}

func (s *Service) DoWork(name string, version string) (*models.PackageVersion, error) {
	resp, err := s.pkgClient.Get(name, version)
	if err != nil {
		return nil, err
	}

	dependencies := make(map[string]*models.PackageVersion, 0)
	for iName, iVersion := range resp.Dependencies {
		iRes, iErr := s.DoWork(iName, iVersion)
		if iErr != nil {
			return nil, iErr
		}
		dependencies[iName] = iRes
	}

	return &models.PackageVersion{
		Name:         resp.Name,
		Version:      resp.Version,
		Dependencies: dependencies,
	}, nil
}
