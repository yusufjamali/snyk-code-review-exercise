package packagehandler

import (
	"encoding/json"
	"github.com/Masterminds/semver/v3"
	"github.com/gorilla/mux"
	"github.com/snyk/snyk-code-review-exercise/internal/models"
	"net/http"
)

type Handler struct {
	depResSvc DependencyResolverService
}

type DependencyResolverService interface {
	DoWork(name string, version string) (*models.PackageVersion, error)
}

func New(depResSvc DependencyResolverService) *Handler {
	return &Handler{depResSvc: depResSvc}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkgName := vars["package"]
	pkgVersion := vars["version"]

	// validate inputs
	_, err := semver.NewVersion(pkgVersion)
	if err != nil {
		println(err.Error())
		// bad request
		w.WriteHeader(400)
		return
	}

	svcRes, err := h.depResSvc.DoWork(pkgName, pkgVersion)
	if err != nil {
		println(err.Error())
		// internal server error
		w.WriteHeader(500)
		return
	}

	stringified, err := json.Marshal(svcRes)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(stringified)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		return
	}
}
