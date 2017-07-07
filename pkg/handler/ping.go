package handler

import (
	"net/http"

	"github.com/unrolled/render"
)

// Ping handles simple healthcheck requests.
func Ping(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Text(
			w,
			http.StatusOK,
			"OK",
		)
	}
}
