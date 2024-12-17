package models

type PackageVersion struct {
	Name         string                     `json:"name"`
	Version      string                     `json:"version"`
	Dependencies map[string]*PackageVersion `json:"dependencies"`
}
