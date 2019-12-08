package handler

import (
	"fmt"
	"net/http"
)

// Query :
type Query struct {
}

func (b *Query) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(w, string("success"))
}
