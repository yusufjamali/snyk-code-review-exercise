package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/gorilla/mux"
)

func New() http.Handler {
	router := mux.NewRouter()
	router.Handle("/package/{package}/{version}", http.HandlerFunc(packageHandler))
	return router
}

type npmPackageMetaResponse struct {
	Versions map[string]npmPackageResponse `json:"versions"`
}

type npmPackageResponse struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type NpmPackageVersion struct {
	Name         string                        `json:"name"`
	Version      string                        `json:"version"`
	Dependencies map[string]*NpmPackageVersion `json:"dependencies"`
}

func packageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkgName := vars["package"]
	pkgVersion := vars["version"]

	// review: should we consider validating the requested pkg version upfront ? this would help us fail
	// early and free up resources for other requests at the earliest.
	// consider returning 4xx error codes on any validation failures
	// constraint, err := semver.NewConstraint(constraintStr)
	//	if err != nil {
	//		return "", err
	//	}

	// minor review: I wonder if we should explore the possibility of the recursive function not needing to
	// mutate a passed in object, (ie don't pass in &NpmPackageVersion but instead make it return it)
	// imo it would be cleaner and prevent any accidental mutations ?
	rootPkg := &NpmPackageVersion{Name: pkgName, Dependencies: map[string]*NpmPackageVersion{}}
	if err := resolveDependencies(rootPkg, pkgVersion); err != nil {
		println(err.Error())
		w.WriteHeader(500)
		return
	}

	// review: should we use the json.Marshall function instead, it would reduce the payload size
	// and the amount of data being transmitted, improving e2e latency, the clients could then
	// be responsible for formatting the response as desirable
	stringified, err := json.MarshalIndent(rootPkg, "", "  ")
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	// Ignoring ResponseWriter errors
	// review: should we consider logging the error for observability ?
	// also if there is an error should we return a non 2xx ?
	_, _ = w.Write(stringified)
}

func resolveDependencies(pkg *NpmPackageVersion, versionConstraint string) error {
	// review: should we abstract away the retrieval of the highest compatible version into another
	// function, it would reduce the cognitive load when trying to figure out what this function is
	// doing ?
	pkgMeta, err := fetchPackageMeta(pkg.Name)
	if err != nil {
		return err
	}
	concreteVersion, err := highestCompatibleVersion(versionConstraint, pkgMeta)
	if err != nil {
		return err
	}
	pkg.Version = concreteVersion

	// idea: should we introduce an in memory cache that is keyed on name_version and stores the *NpmPackageVersion
	// object, this could reduce the redundant processing when building dependency calls.
	// while the in-memory cache would have disadvantages like
	// a)  memory footprint
	// b)  restart of service would wipe out cache
	// c)  hit rate might not be ideal when multiple instance of the service are running
	// the approach might help us understand if there would be value in investing/managing a persistent cache.
	npmPkg, err := fetchPackage(pkg.Name, pkg.Version)
	if err != nil {
		return err
	}

	for dependencyName, dependencyVersionConstraint := range npmPkg.Dependencies {
		dep := &NpmPackageVersion{Name: dependencyName, Dependencies: map[string]*NpmPackageVersion{}}
		pkg.Dependencies[dependencyName] = dep
		if err := resolveDependencies(dep, dependencyVersionConstraint); err != nil {
			return err
		}
	}
	return nil
}

func highestCompatibleVersion(constraintStr string, versions *npmPackageMetaResponse) (string, error) {
	// review: this check might become obsolete once the validation of the passed in semver happens upfront,
	// could we potentially remove this ?
	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return "", err
	}
	filtered := filterCompatibleVersions(constraint, versions)
	// review: could we remove the need for sorting here, if the above line returns the max version ?
	sort.Sort(filtered)
	if len(filtered) == 0 {
		return "", errors.New("no compatible versions found")
	}
	return filtered[len(filtered)-1].String(), nil
}

func filterCompatibleVersions(constraint *semver.Constraints, pkgMeta *npmPackageMetaResponse) semver.Collection {
	// review: should we consider sending a single concrete version back from this function ?
	// we could keep track of the max version here as we loop through all the versions
	// with something like semVer.GreaterThanEqualTo(maxVersion)
	var compatible semver.Collection
	for version := range pkgMeta.Versions {
		semVer, err := semver.NewVersion(version)
		if err != nil {
			continue
		}
		if constraint.Check(semVer) {
			compatible = append(compatible, semVer)
		}
	}
	return compatible
}

func fetchPackage(name, version string) (*npmPackageResponse, error) {
	// review: should we consider creating a custom http client, to manage transport level properties
	// for concurrency etc.. ?
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// review: should we check status here and return err if not found ?
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed npmPackageResponse
	// review: should check the error and return the error instead ?
	_ = json.Unmarshal(body, &parsed)
	return &parsed, nil
}

func fetchPackageMeta(p string) (*npmPackageMetaResponse, error) {
	// review: should we consider creating a custom http client, to manage transport level properties
	// for concurrency etc.. ?
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/%s", p))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed npmPackageMetaResponse
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

// IDEAS
// Refactor and testing
// 1. I would consider splitting this file further into
// 	 - service layer - would be responsible for computing the dependencies graph, effectively would contain the recursive function
//                     pros: a) can easily test the recursive function w/o any external dependency.
//							 b) can re-use the same service logic for different client interfaces eg: grpc
//							 c) would be easier re-use the service logic for other non-npm registries
//   - client layer - would be responsible for interacting with external npm API's
// 						a) abstract away the integration detail specific to on registry in a single place
// the current file - would be responsible for request/response validation and interfacing with clients

// ARCHITECTURE
// 1. once we have proved the efficacy of the in-memory cache should we consider storing this in a persistent cache like redis ?
//    some considerations
//    a) how do we add keys to the cache - one off process / base it on incoming requests
//    b) how do we update cache - when npm registry update occurs - background worker ?
