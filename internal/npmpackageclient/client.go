package npmpackageclient

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/snyk/snyk-code-review-exercise/internal/models"
	"io"
	"net/http"
)

type Client struct {
	hClient http.Client
	url     string
}

func New(url string) *Client {
	return &Client{hClient: http.Client{}, url: url}
}

func (c *Client) Get(name string, version string) (_ *models.PackageClientResponse, err error) {
	version, err = c.getHighestCompatibleVersion(name, version)
	if err != nil {
		return nil, err
	}

	pkgResp, err := c.hClient.Get(fmt.Sprintf("%s/%s/%s", c.url, name, version))
	if err != nil {
		return nil, err
	}
	defer pkgResp.Body.Close()

	if pkgResp.StatusCode != 200 {
		return nil, fmt.Errorf("unsuccessful response: %s", pkgResp.Status)
	}

	body, err := io.ReadAll(pkgResp.Body)
	if err != nil {
		return nil, err
	}

	var parsed models.PackageClientResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (c *Client) getHighestCompatibleVersion(name string, version string) (string, error) {
	resp, err := c.hClient.Get(fmt.Sprintf("%s/%s", c.url, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed models.PackageClientMetaResponse
	if err = json.Unmarshal([]byte(body), &parsed); err != nil {
		return "", err
	}

	// We could ignore errors here as the versions are returned by the registry
	// and unlikely to be incorrect
	constraint, _ := semver.NewConstraint(version)

	// let's find the highest compatible version
	maxVersion := &semver.Version{}
	for _, v := range parsed.Versions {
		semVer, iErr := semver.NewVersion(v.Version)
		if iErr != nil {
			continue
		}
		if constraint.Check(semVer) && semVer.GreaterThan(maxVersion) {
			maxVersion = semVer
		}
	}
	return maxVersion.String(), nil
}
