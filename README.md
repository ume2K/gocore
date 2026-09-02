# GoCore Framework

A lightweight, opinionated Go web framework with a custom radix-trie router, middleware pipeline, recursive HTML template loading, and an integrated **Zero-Config Asset Pipeline** for SCSS and modern JavaScript.

## Features

* **Custom Trie Router:** Fast route matching with path parameters (`:param`), route grouping, and configurable 404 handlers.

* **Recursive Template Discovery:** Automatic directory walking that registers templates by both relative subpath (e.g., `episodes/vietnam/vietnam.html`) and flat base filename.

* **Modern Asset Pipeline:** Development-time SCSS (Dart Sass) and JS (esbuild) bundling with multi-entrypoint support, ESM format, and code splitting.

* **WebSocket & Protocol Upgrade Support:** `StatusRecorder` implements `http.Hijacker` so request logging and recovery middleware do not disrupt WebSockets or SSE connections.

* **Live Reload:** Instant rebuilds in development via Air, paired with background `fsnotify` file watchers for frontend assets.

* **Production Ready:** Multi-stage Docker build producing a minimal container with pre-compiled, optimized static assets.

* **Context Handlers:** Unified request context and convenience responders via `func(c *framework.Context)`.

## Layout

```text
gocore-framework/
├── assets/
│   ├── js/             # JS sources; configurable entry points (e.g., main.js, ...)
│   └── scss/           # SCSS sources; main.scss entrypoint and partials
├── cmd/server/         # Application entrypoint (main.go)
├── pkg/framework/      # Router, Context, compilers, and middleware
├── public/             # Compiled CSS/JS and static assets, served at /assets
│   ├── css/
│   ├── js/
│   └── images/
├── views/              # HTML templates (recursively scanned, subfolders supported)
├── Dockerfile
├── Makefile
└── .air.toml

```

## Assets (Images, SCSS & JS)

### Images

* **Storage:** Place raw images in `public/images/`.

* **Routing:** Access images via `/assets/images/...` (the router maps `/assets` to `./public`).

  * HTML: `<img src="/assets/images/logo.png" alt="Logo">`

  * SCSS: `background-image: url('/assets/images/bg.jpg');`

### SCSS Compilation

* **Entrypoint:** The compiler processes `assets/scss/main.scss`.

* **Partials:** Organize modular partials using an underscore prefix (e.g., `assets/scss/components/_nav.scss`).

* **Imports:** Reference partials without the leading underscore or `.scss` extension:

  ```
  // assets/scss/main.scss
  @import 'variables';
  @import 'utils';
  @import 'components/nav';
  @import 'pages/home';
  
  ```

### JavaScript Bundler

* **Configurable Entry Points:** Specify one or more entry points using `JSConfig.EntryPoints`. Defaults to `main.js` if left empty.

* **ESM & Code Splitting:** Bundles code using `--format=esm`, `--splitting`, minification, and source maps into the configured output folder.

* **Configuration:**

  ```
  js := framework.NewJSCompiler(framework.JSConfig{
      SourceDir:   "assets/js",
      OutputDir:   "public/js",
      EntryPoints: []string{
          "main.js",
          "components/whiteboard.js",
      },
  })
  
  ```

## HTML Templates & Components

Call `r.LoadHTMLGlob("views")` once during initialization. The engine recursively walks the directory tree and registers every `.html` file under two keys:

1. **Relative Subpath:** The forward-slash relative path from the root view directory (e.g., `episodes/vietnam/vietnam.html` or `components/card.html`). Use this to isolate templates and prevent naming collisions across nested modules.

2. **Base Filename:** The flat filename (e.g., `vietnam.html`, `index.html`), preserved for backward-compatible lookups.

### Rendering Views

Render templates using either their relative path or base name:

```
r.GET("/episoden/vietnam", func(c *framework.Context) {
    c.HTML(http.StatusOK, "episodes/vietnam/vietnam.html", PageData{Title: "Vietnam"})
})

```

### Component Partials

Include reusable partials within templates using standard Go template actions:

```
{{ template "components/nav.html" . }}
<!-- or by base name -->
{{ template "nav.html" . }}

```

*Note: The trailing dot (`.`) passes the active data context into the partial.*

## Prerequisites

* [Go](https://go.dev/) 1.25+

* [Dart Sass](https://sass-lang.com/install) (`npm install -g sass`)

* [esbuild](https://esbuild.github.io/) (`npm install -g esbuild`)

* [Air](https://github.com/air-verse/air) (`go install github.com/air-verse/air@latest`)

## Quick Start

1. Create your environment configuration:

   ```
   cp .env.example .env
   
   ```

2. Set configuration values in `.env`:

   ```
   APP_NAME=gocore
   PORT=8080
   APP_ENV=development
   
   ```

3. Start the development server with live reload:

   ```
   make dev
   
   ```

4. Navigate to `http://localhost:8080`.

## Make Targets

| Command | Description | 
 | ----- | ----- | 
| `make dev` | Starts Air live reload alongside SCSS/JS file watchers | 
| `make run` | Performs a single asset build and runs the Go server | 
| `make build` | Compiles the production Go binary to `bin/server` | 
| `make docker-build` | Builds a self-contained multi-stage production Docker image | 
| `make docker-run` | Runs the production container on port `8080` | 
| `make clean` | Removes compiled binaries and Docker build artifacts | 

## Usage Example

```
package main

import (
    "gocore/pkg/framework"
    "net/http"
)

type PageData struct {
    Title string
}

func main() {
    framework.LoadEnv(".env")

    r := framework.NewRouter()
    r.Use(framework.Logger)
    r.Use(framework.Recovery)

    // Recursively scans the views/ directory
    r.LoadHTMLGlob("views")
    r.Static("/assets", "./public")

    r.GET("/", func(c *framework.Context) {
        c.HTML(http.StatusOK, "index.html", PageData{Title: "GoCore"})
    })

    // Nested view lookup by relative path
    r.GET("/episoden/vietnam", func(c *framework.Context) {
        c.HTML(http.StatusOK, "episodes/vietnam/vietnam.html", PageData{Title: "Vietnam"})
    })

    r.POST("/api/data", func(c *framework.Context) {
        var input struct {
            Name string `json:"name"`
        }
        if err := c.BindJSONStrict(&input); err != nil {
            c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, map[string]string{"status": "ok"})
    })

    r.SetNotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
    }))

    http.ListenAndServe(":8080", r)
}

```

## Docker Deployment

Build and run using the provided production target:

```
make docker-build
make docker-run

```

The multi-stage `Dockerfile` compiles SCSS and JavaScript during the build stage; the final deployment image contains only the compiled Go binary, the `views/` folder, and the generated `public/` directory.

## License

MIT