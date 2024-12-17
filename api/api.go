package api

import (
	"github.com/gorilla/mux"
	"github.com/snyk/snyk-code-review-exercise/api/packagehandler"
	"github.com/snyk/snyk-code-review-exercise/internal/dependencyresolver"
	"github.com/snyk/snyk-code-review-exercise/internal/npmpackageclient"
	"net/http"
)

func New() http.Handler {
	client := npmpackageclient.New("https://registry.npmjs.org")
	service := dependencyresolver.New(client)
	handler := packagehandler.New(service)

	router := mux.NewRouter()
	router.Handle("/package/{package}/{version}", http.HandlerFunc(handler.Handle))
	return router
}
