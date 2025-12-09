# GoCore Framework

A lightweight, opinionated Go web framework built for speed and developer experience. It features a custom high-performance router, middleware support, and a fully integrated **Zero-Config Asset Pipeline** for SCSS and JavaScript.

## 🚀 Key Features

* **Custom Trie-Router:** Fast, efficient routing with support for dynamic parameters (`/users/:id`), groups (`/api/v1`), and custom 404 handling.
* **Integrated Asset Pipeline:** No Webpack/Vite required. The Go server orchestrates asset compilation directly using `sass` and `esbuild`.
* **Smart Hot-Reloading:**
    * **Go Code:** Watched by [Air](https://github.com/air-verse/air). Restarts the binary on change.
    * **Assets (SCSS/JS):** Watched by internal Go routines using `fsnotify` (recursive). Recompiles instantly *without* restarting the server.
* **Production Ready:** Docker support with Multi-Stage builds (assets compiled at build time).
* **Context-Based Handlers:** Clean API `func(c *framework.Context)` for handling JSON, HTML, and Request binding.

## 📂 Architecture Overview

### Directory Structure

```text
gocore/
├── assets/             # Source files for assets
│   ├── js/             # Modern ES6+ JavaScript (bundled via esbuild)
│   └── scss/           # SCSS styles (compiled via dart-sass)
├── cmd/
│   └── server/         # Entry point (main.go)
├── pkg/
│   └── framework/      # The core framework logic (Router, Context, Compiler)
├── public/             # Served statically at /assets
│   ├── css/            # Generated CSS (do not edit)
│   ├── js/             # Generated JS (do not edit)
│   └── img/            # Static images (WebP recommended)
├── views/              # HTML Templates (Go html/template)
│   └── templates/      # Partials (nav, footer, head)
├── Dockerfile          # Multi-stage production build
├── Makefile            # Command shortcuts
└── .air.toml           # Hot-reload configuration
```

Hier ist der komplette Inhalt für deine README.md, fertig formatiert zum Kopieren.Markdown# GoCore Framework

A lightweight, opinionated Go web framework built for speed and developer experience. It features a custom high-performance router, middleware support, and a fully integrated **Zero-Config Asset Pipeline** for SCSS and JavaScript.

## 🚀 Key Features

* **Custom Trie-Router:** Fast, efficient routing with support for dynamic parameters (`/users/:id`), groups (`/api/v1`), and custom 404 handling.
* **Integrated Asset Pipeline:** No Webpack/Vite required. The Go server orchestrates asset compilation directly using `sass` and `esbuild`.
* **Smart Hot-Reloading:**
    * **Go Code:** Watched by [Air](https://github.com/air-verse/air). Restarts the binary on change.
    * **Assets (SCSS/JS):** Watched by internal Go routines using `fsnotify` (recursive). Recompiles instantly *without* restarting the server.
* **Production Ready:** Docker support with Multi-Stage builds (assets compiled at build time).
* **Context-Based Handlers:** Clean API `func(c *framework.Context)` for handling JSON, HTML, and Request binding.

## 📂 Architecture Overview

### Directory Structure

```text
gocore/
├── assets/             # Source files for assets
│   ├── js/             # Modern ES6+ JavaScript (bundled via esbuild)
│   └── scss/           # SCSS styles (compiled via dart-sass)
├── cmd/
│   └── server/         # Entry point (main.go)
├── pkg/
│   └── framework/      # The core framework logic (Router, Context, Compiler)
├── public/             # Served statically at /assets
│   ├── css/            # Generated CSS (do not edit)
│   ├── js/             # Generated JS (do not edit)
│   └── img/            # Static images (WebP recommended)
├── views/              # HTML Templates (Go html/template)
│   └── templates/      # Partials (nav, footer, head)
├── Dockerfile          # Multi-stage production build
├── Makefile            # Command shortcuts
└── .air.toml           # Hot-reload configuration
```

### The Asset Pipeline Logic
Unlike traditional setups where a Node.js server runs alongside Go, **GoCore controls the compilers directly**:

1.  **Dev Mode (`APP_ENV=development`):** The app starts background goroutines that watch `assets/scss` and `assets/js`. Changes trigger a recompile via CLI tools (`sass`, `esbuild`) immediately.
2.  **Prod Mode (`APP_ENV=production`):** The watchers are disabled. The app expects assets to be pre-compiled (handled by the Dockerfile during the build stage).

## 🛠 Prerequisites

Ensure you have the following tools installed globally on your machine:

1.  **Go** (1.25+)
2.  **Sass** (Dart Sass)
    * Mac: `brew install sass/sass/sass`
    * Windows: `npm install -g sass` or via Chocolatey
3.  **ESBuild**
    * Mac: `brew install esbuild`
    * Windows: `npm install -g esbuild`
4.  **Air** (for Go hot-reloading)
    * Run: `go install github.com/air-verse/air@latest`

## ⚡ Getting Started

1.  **Clone the repository**
2.  **Setup Environment**
    Copy `.env.example` to `.env`:
    ```ini
    APP_NAME=gocore
    PORT=8080
    APP_ENV=development
    ```
3.  **Run Development Server**
    ```bash
    make dev
    ```
    Access the app at `http://localhost:8080`.

## 📜 Make Commands

The project uses a `Makefile` to simplify common tasks.

| Command | Description |
| :--- | :--- |
| `make dev` | **The main dev command.** Starts `air` for Go and internal watchers for Assets. |
| `make run` | Runs `go run main.go` (Compiles assets once, no hot-reload). |
| `make build` | Compiles the binary to `bin/server` (or `server.exe` on Windows). |
| `make docker-build` | Builds the optimized production Docker image (compiles assets inside Docker). |
| `make docker-run` | Runs the Docker container locally (Production simulation). |
| `make clean` | Removes binary and temporary files. |

## 🐳 Deployment (Docker)

The `Dockerfile` uses a **Multi-Stage Build** to keep the final image minimal (`scratch` or `alpine`).

1.  **Builder Stage:** Installs Node/Sass, compiles SCSS/JS, and builds the Go binary.
2.  **Final Stage:** Copies only the binary, the `views`, and the compiled `public` assets.

**To deploy:**
```bash
# 1. Build
make docker-build

# 2. Run (Simulate production)
make docker-run
# OR manually:
docker run -p 8080:8080 -e APP_ENV=production gocore
```

## 🧩 Framework Usage

### Router & Handlers
The router uses a Radix Tree (Trie) for O(n) lookup performance.

```go
r := framework.NewRouter()

// Static files (mapped to ./public)
r.Static("/assets", "./public")

// HTML Views
r.GET("/", func(c *framework.Context) {
    data := PageData{Title: "Home"}
    c.HTML(200, "index.html", data)
})

// JSON API
r.POST("/api/data", func(c *framework.Context) {
    var input MyStruct
    if err := c.BindJSONStrict(&input); err != nil {
         c.JSON(400, map[string]string{"error": err.Error()})
         return
    }
    c.JSON(200, map[string]string{"status": "ok"})
})

// Custom 404 Redirect
r.SetNotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}))
```

### Templates
Templates are loaded globally from views/*.html.
Define Partials: Use `{{ define "nav" }} ... {{ end }}` inside views/templates/.
Include Partials: Use `{{ template "nav" . }}` (Don't forget the dot to pass context!).

License MIT
