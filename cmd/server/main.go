package main

import (
	"context"
	"fmt"
	"gocore/pkg/framework"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PageData struct {
	Title string
}

func main() {
	framework.LoadEnv(".env")

	isDev := framework.GetEnv("APP_ENV", "development") == "development"

	if isDev {
		fmt.Println("GoCore Framework – Development Mode")
		log.Println("Initializing SCSS Compiler...")
		scss := framework.NewSCSSCompiler(framework.SCSSConfig{
			SourceDir: "assets/scss",
			OutputDir: "public/css",
			Debug:     true,
		})
		if err := scss.Compile(); err != nil {
			log.Println("SCSS compilation failed")
		}

		log.Println("Initializing JS Bundler...")
		js := framework.NewJSCompiler(framework.JSConfig{
			SourceDir: "assets/js",
			OutputDir: "public/js",
		})
		if err := js.Bundle(); err != nil {
			log.Println("JS Bundle failed")
		}

		go scss.Watch()
		go js.Watch()
	}

	r := framework.NewRouter()
	r.Use(framework.Logger)
	r.Use(framework.Recovery)

	r.LoadHTMLGlob("views/*.html")
	r.LoadHTMLGlob("views/components/*.html")
	r.Static("/assets", "./public")

	r.GET("/", func(c *framework.Context) {
		c.HTML(http.StatusOK, "index.html", PageData{Title: "GoCore"})
	})

	r.GET("/api/health", func(c *framework.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	port := framework.GetEnv("PORT", "8080")
	addr := ":" + port

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Server running on http://localhost%s (bind: 0.0.0.0%s)", addr, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
