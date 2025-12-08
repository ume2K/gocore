package framework

import "net/http"

type StatusRecorder struct {
	http.ResponseWriter
	Status  int
	Written bool
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

// WriteHeader captures the status code explicitly.
func (r *StatusRecorder) WriteHeader(code int) {
	r.Status = code
	r.Written = true
	r.ResponseWriter.WriteHeader(code)
}

// Write captures "Implicit 200 OK".
// If data is written without a header, Go defaults to 200. We must record that.
func (r *StatusRecorder) Write(b []byte) (int, error) {
	if !r.Written {
		r.Status = http.StatusOK
		r.Written = true
	}
	return r.ResponseWriter.Write(b)
}

// Unwrap provides compatibility with http.ResponseController (Go 1.20+)
func (r *StatusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
