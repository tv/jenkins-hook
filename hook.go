package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

func main() {

	urlRoot := os.Getenv("URL_ROOT")
	if urlRoot == "" {
		urlRoot = "http://localhost/github-webhook"
	}

	serverUrl, _ := url.Parse(urlRoot)

	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	api := sameHost(reverseProxy)

	http.ListenAndServe(":"+port, api)
}
