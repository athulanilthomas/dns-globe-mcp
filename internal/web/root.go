package web

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed welcome.txt
var welcome string

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, welcome)
	})
}
