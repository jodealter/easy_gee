package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	gee := gee.New()
	gee.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.WriteHeader(404)
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	gee.Run(":9999")

}
