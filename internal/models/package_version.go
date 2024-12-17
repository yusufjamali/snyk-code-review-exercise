package models

type PackageVersion struct {
	Name         string
	Version      string
	Dependencies map[string]*PackageVersion
}
