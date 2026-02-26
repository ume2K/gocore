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

func (r *StatusRecorder) WriteHeader(code int) {
	r.Status = code
	r.Written = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *StatusRecorder) Write(b []byte) (int, error) {
	if !r.Written {
		r.Status = http.StatusOK
		r.Written = true
	}
	return r.ResponseWriter.Write(b)
}

func (r *StatusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
