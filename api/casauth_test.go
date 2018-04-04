package api

import (
	"testing"
	"github.com/go-chi/chi"
)

func TestCasAuth(t *testing.T) {
	r := chi.NewRouter()
	r.Use(casauth)

}
