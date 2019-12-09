package handler

import (
	"net/http"
)

// RegistALl ...
func RegistALl(registFunc func(string, http.Handler)) {
	registFunc("/query", &Query{})
}
