package models

type PackageClientResponse struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type PackageClientMetaResponse struct {
	Versions map[string]PackageClientResponse `json:"versions"`
}
