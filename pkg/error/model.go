package error

import "net/http"

var (
	ErrFailedProcessPassword = NewCustomError("Failed Process Password", "FAILED_PROCESS_PASSWORD", http.StatusBadGateway)
)
