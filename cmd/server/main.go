package main

import (
	"context"
	"fmt"
	"gocore/pkg/framework"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type Product struct {
	Name        string
	Description string
	Price       float64
	ImageURL    string
}

type OpeningHour struct {
	DayLabel  string
	DaySchema string
	Open      string
	Close     string
}

type PageData struct {
	Title           string
	MetaDescription string
	CurrentYear     int
	CurrentDate     time.Time
	Products        []Product
	OpeningHours    []OpeningHour
}

func NewPageData(title string) PageData {
	return PageData{
		Title:           title,
		MetaDescription: "Willkommen auf dem Naturhof Rinisgmüende – Frische Produkte, Tiere und Events direkt vom Hof.",
		CurrentYear:     time.Now().Year(),
		CurrentDate:     time.Now(),
		OpeningHours: []OpeningHour{
			{
				DayLabel:  "Montag - Samstag",
				DaySchema: "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday",
				Open:      "08:00",
				Close:     "20:00",
			},
			{
				DayLabel:  "Sonntag",
				DaySchema: "Sunday",
				Open:      "Geschlossen",
				Close:     "",
			},
		},
	}
}

func main() {
	framework.LoadEnv(".env")

	isDev := framework.GetEnv("APP_ENV", "development") == "development"

	if isDev {
		banner := `
██████╗ ███████╗██╗   ██╗    ███╗   ███╗ ██████╗ ██████╗ ███████╗
██╔══██╗██╔════╝██║   ██║    ████╗ ████║██╔═══██╗██╔══██╗██╔════╝
██║  ██║█████╗  ██║   ██║    ██╔████╔██║██║   ██║██║  ██║█████╗  
██║  ██║██╔══╝  ╚██╗ ██╔╝    ██║╚██╔╝██║██║   ██║██║  ██║██╔══╝  
██████╔╝███████╗ ╚████╔╝     ██║ ╚═╝ ██║╚██████╔╝██████╔╝███████╗
╚═════╝ ╚══════╝  ╚═══╝      ╚═╝     ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝
`

		fmt.Println(banner)
		log.Println("🔧 Initializing SCSS Compiler...")

		scss := framework.NewSCSSCompiler(framework.SCSSConfig{
			SourceDir: "assets/scss",
			OutputDir: "public/css",
			Debug:     true,
		})

		if err := scss.Compile(); err != nil {
			log.Println("⚠️ SCSS compilation failed")
		}

		log.Println("🔧 Initializing JS Bundler...")

		js := framework.NewJSCompiler(framework.JSConfig{
			SourceDir: "assets/js",
			OutputDir: "public/js",
		})

		if err := js.Bundle(); err != nil {
			log.Println("⚠️ JS Bundle failed")
		}

		go scss.Watch()
		go js.Watch()
	}

	r := framework.NewRouter()

	funcMap := template.FuncMap{
		"currency": func(amount float64) string {
			return fmt.Sprintf("CHF %.2f", amount)
		},

		"formatDate": func(t time.Time) string {
			return t.Format("02.01.2006")
		},

		"lastItem": func(index int, slice any) bool {
			v := reflect.ValueOf(slice)
			if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
				return true
			}
			return index == v.Len()-1
		},

		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	r.SetFuncMap(funcMap)

	r.Use(framework.Logger)
	r.Use(framework.Recovery)

	// r.SetNotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// }))

	r.LoadHTMLGlob("views/*.html")
	r.LoadHTMLGlob("views/components/*.html")

	r.Static("/assets", "./public")

	// api := r.Group("/api")
	// api.GET("/test", func(c *framework.Context) {
	// 	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	// })

	// r.POST("/", func(c *framework.Context) {
	// 	var dump any
	// 	if err := c.BindJSONStrict(&dump); err != nil {
	// 		if httpErr, ok := err.(*framework.HTTPError); ok {
	// 			c.JSON(httpErr.Code, map[string]string{
	// 				"error":   "Invalid Request",
	// 				"details": httpErr.Message,
	// 			})
	// 		} else {
	// 			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Error"})
	// 		}
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, map[string]string{"status": "accepted"})
	// })

	r.GET("/", func(c *framework.Context) {
		data := NewPageData("Naturhof Rinisgmüende")
		data.MetaDescription = "Willkommen auf dem Naturhof Rinisgmüende - Frische Produkte, Tiere und Events direkt vom Hof."
		data.Products = []Product{
			{Name: "Saison-Korb 'Klein'", Description: "Gemischtes Gemüse, 3kg", Price: 25.00, ImageURL: "/assets/img/basket.jpg"},
			{Name: "Hof-Honig", Description: "Waldhonig", Price: 12.50, ImageURL: "/assets/img/honey.jpg"},
			{Name: "Kartoffeln 'Agria'", Description: "Mehligkochend, 5kg", Price: 8.00, ImageURL: "/assets/img/potatoes.jpg"},
		}
		c.HTML(http.StatusOK, "index.html", data)
	})

	r.GET("/tiere", func(c *framework.Context) {
		data := NewPageData("Tiere | Naturhof Rinisgmüende")
		data.MetaDescription = "Lernen Sie unsere Kühe, Ziegen und Hühner kennen. Artgerechte Tierhaltung auf dem Naturhof."
		c.HTML(http.StatusOK, "tiere.html", data)
	})

	r.GET("/eventraum", func(c *framework.Context) {
		data := NewPageData("Eventraum | Naturhof Rinisgmüende")
		data.MetaDescription = "Lernen Sie unsere Kühe, Ziegen und Hühner kennen. Artgerechte Tierhaltung auf dem Naturhof."
		c.HTML(http.StatusOK, "eventraum.html", data)
	})

	r.GET("/besuch", func(c *framework.Context) {
		data := NewPageData("Besuch | Naturhof Rinisgmüende")
		data.MetaDescription = "Lernen Sie unsere Kühe, Ziegen und Hühner kennen. Artgerechte Tierhaltung auf dem Naturhof."
		c.HTML(http.StatusOK, "besuch.html", data)
	})

	// r.GET("/oops", func(c *framework.Context) {
	// 	c.HTML(http.StatusOK, "placeholder.html", nil)
	// })

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
		log.Printf("Server running on http://0.0.0.0%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
