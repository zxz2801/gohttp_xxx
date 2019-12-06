package main

import (
	"fmt"
	"net/http"
)

// query :
type query struct {
}

func (b *query) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(w, string("success"))
}
