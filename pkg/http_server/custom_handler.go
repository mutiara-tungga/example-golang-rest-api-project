package httpserver

import "net/http"

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (h HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// TODO: send alert

		WriteJsonError(r.Context(), w, err)
	}
}
