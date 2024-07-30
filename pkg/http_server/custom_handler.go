package httpserver

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// TODO: send alert

		WriteJsonError(r.Context(), w, err)
	}
}
