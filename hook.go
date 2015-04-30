package main

import (
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func setSameHost(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.Host = r.URL.Host
	next(rw, r)
}

func main() {

	urlRoot := os.Getenv("JENKINS_ROOT")
	if urlRoot == "" {
		urlRoot = "http://localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverUrlJenkins, _ := url.Parse(urlRoot)

	reverseProxyJenkins := httputil.NewSingleHostReverseProxy(serverUrlJenkins)

	r := httprouter.New()

	for _, url := range []string{"/github-webhook/*all", "/ghprbhook/*all"} {
		for _, val := range []string{"GET", "POST", "PUT"} {
			r.Handler(val, url, reverseProxyJenkins)
		}
	}

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(setSameHost))
	n.UseHandler(r)
	n.Run(":" + port)
}
