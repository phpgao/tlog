package handler

import (
	"net/http"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	mux := http.NewServeMux()
	RegisterHandler(mux)
}

func TestRegisterHandlerWithPath(t *testing.T) {
	mux := http.NewServeMux()
	RegisterHandlerWithPath(mux, "/")
	http.ListenAndServe(":8080", mux)
}
