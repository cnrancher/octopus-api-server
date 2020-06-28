package ui

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cnrancher/edge-api-server/pkg/settings"
	"github.com/rancher/apiserver/pkg/parse"
	"github.com/sirupsen/logrus"
)

var (
	insecureClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func Content() http.Handler {
	return http.FileServer(http.Dir(settings.UIPath.Get()))
}

func UI(next http.Handler) http.Handler {
	_, err := os.Stat(indexHTML())
	if err != nil {
		logrus.Warnf("failed to get local index html: %s, error: %s", indexHTML(), err.Error())
	}
	local := err == nil
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if parse.IsBrowser(req, true) {
			if local && settings.UIIndex.Get() == "local" {
				http.ServeFile(resp, req, indexHTML())
			} else {
				ui(resp, req)
			}
		} else {
			next.ServeHTTP(resp, req)
		}
	})
}

func indexHTML() string {
	return filepath.Join(settings.UIPath.Get(), "index.html")
}

func ui(resp http.ResponseWriter, req *http.Request) {
	if err := serveIndex(resp, req); err != nil {
		logrus.Errorf("failed to serve UI: %v", err)
		resp.WriteHeader(500)
	}
}

func serveIndex(resp http.ResponseWriter, req *http.Request) error {
	r, err := insecureClient.Get(settings.UIIndex.Get())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	_, err = io.Copy(resp, r.Body)
	return err
}

func JSURLGetter() string {
	if settings.UIIndex.Get() == "local" {
		return "/api-ui/ui.min.js"
	}
	return ""
}

func CSSURLGetter() string {
	if settings.UIIndex.Get() == "local" {
		return "/api-ui/ui.min.css"
	}
	return ""
}
