package httpmiddleware

import (
	"fmt"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/log"
	"net/http"
	"runtime/debug"
)

var (
	ErrPanic = pkgErr.NewCustomError("unexpected panic error", "PANIC_ERROR", http.StatusInternalServerError)
)

// Recoverer wraps an HTTP handler and recovers from any panics that occur during
// its execution.
func PanicRecoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				ctx := r.Context()
				if r.Header.Get("Connection") != "Upgrade" {
					debugStack := debug.Stack()
					err, ok := rvr.(error)
					if !ok {
						err = fmt.Errorf("panic occurred : %v", rvr)
					}
					log.Error(ctx, "unexpected panic error occured", err, log.LogField{Key: "stack", Value: debugStack})

					httpserver.WriteJsonError(ctx, w, ErrPanic)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
