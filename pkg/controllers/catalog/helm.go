package catalog

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	"github.com/cnrancher/octopus-api-server/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	httpTimeout = time.Second * 300
	httpClient  = &http.Client{
		Timeout: httpTimeout,
	}
)

type Helm struct {
	catalogName string
	Hash        string
	url         string
	username    string
	password    string
}

func (h *Helm) request(pathURL string) (*http.Response, error) {
	baseEndpoint, err := url.Parse(pathURL)
	if err != nil {
		return nil, err
	}
	if !baseEndpoint.IsAbs() {
		helmURLstring := h.url
		if !strings.HasSuffix(helmURLstring, "/") {
			helmURLstring = helmURLstring + "/"
		}
		helmURL, err := url.Parse(helmURLstring)
		if err != nil {
			return nil, err
		}
		baseEndpoint = helmURL.ResolveReference(baseEndpoint)
	}

	if err := util.ValidateURL(baseEndpoint.String()); err != nil {
		return nil, err
	}

	if len(h.username) > 0 && len(h.password) > 0 {
		baseEndpoint.User = url.UserPassword(h.username, h.password)
	}
	req, err := http.NewRequest(http.MethodGet, baseEndpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

func (h *Helm) downloadIndex(indexURL string) (*v1alpha1.IndexFile, error) {
	indexURL = strings.TrimSuffix(indexURL, "/")
	indexURL = indexURL + "/index.yaml"
	resp, err := h.request(indexURL)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil, errors.Errorf("Timeout in HTTP GET to [%s], did not respond in %s", indexURL, httpTimeout)
		}
		return nil, errors.Errorf("Error in HTTP GET to [%s], error: %s", indexURL, err)
	}
	defer resp.Body.Close()

	// only return forgot error if status code is unauthorized.
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.Errorf("Unexpected HTTP status code %d from [%s], expected 200", resp.StatusCode, indexURL)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Unexpected HTTP status code %d from [%s], expected 200", resp.StatusCode, indexURL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error while reading response from [%s], error: %s", indexURL, err)
	}
	file := &v1alpha1.IndexFile{}
	err = yaml.Unmarshal(body, file)
	if err != nil {
		return nil, errors.Errorf("error unmarshalling response from [%s]", indexURL)
	}
	return file, nil
}
