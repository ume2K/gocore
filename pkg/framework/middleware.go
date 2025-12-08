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
				// Log the stack trace
				log.Printf("[PANIC] %v\n%s", err, string(debug.Stack()))

				// If the header hasn't been written yet, send 500
				// (Check logic depends on your ResponseWriter wrapper, simplified here)
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
		recorder := NewStatusRecorder(w) // Dein korrigierter Recorder aus response.go

		defer func() {
			duration := time.Since(start)
			status := recorder.Status

			// 1. Alles was Fehler ist (4xx, 5xx) IMMER loggen
			if status >= 400 {
				log.Printf("[ERR]  [%s] %s | Status: %d | Duration: %v", r.Method, r.URL.Path, status, duration)
				return
			}

			// 2. Alles was langsam ist (> 500ms) IMMER loggen
			if duration > 500*time.Millisecond {
				log.Printf("[SLOW] [%s] %s | Status: %d | Duration: %v", r.Method, r.URL.Path, status, duration)
				return
			}

			// 3. Status 200 ignorieren wir (Stille im Log = Alles gut)
		}()

		next.ServeHTTP(recorder, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In Production: Hier nur die echte Domain erlauben (via Env Var)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		// Preflight Requests (Browser fragt: "Darf ich?")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
