package main

import (
	"fmt"
	"net/http"
)

// CreateFileServer to serve file in the `dir` directory as static files
func CreateFileServer(host string, port int64, baseurl, dir string) *http.Server {
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: http.StripPrefix(
			baseurl,
			http.FileServer(http.Dir(dir)),
		),
	}

	return srv
}
