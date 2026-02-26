package framework

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v\n%s", err, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := NewStatusRecorder(w)
		defer func() {
			duration := time.Since(start)
			status := recorder.Status
			if status >= 400 {
				log.Printf("[ERR]  [%s] %s | Status: %d | Duration: %v", r.Method, r.URL.Path, status, duration)
				return
			}
			if duration > 500*time.Millisecond {
				log.Printf("[SLOW] [%s] %s | Status: %d | Duration: %v", r.Method, r.URL.Path, status, duration)
			}
		}()
		next.ServeHTTP(recorder, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
