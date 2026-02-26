# GoCore Framework

A lightweight, opinionated Go web framework with a custom Trie router, middleware support, and an integrated **Zero-Config Asset Pipeline** for SCSS and JavaScript.

## Features

- **Custom Trie-Router:** Fast routing with dynamic parameters (`/users/:id`), groups (`/api/v1`), custom 404
- **Asset Pipeline:** SCSS (Dart Sass) and JS (esbuild) – no Webpack/Vite
- **Hot-Reload:** Go via [Air](https://github.com/air-verse/air); SCSS/JS via internal `fsnotify` watchers (no server restart)
- **Production:** Multi-stage Docker build (assets compiled at build time)
- **Context Handlers:** `func(c *framework.Context)` for HTML, JSON, and request binding

## Layout

```
gocore-framework/
├── assets/
│   ├── js/           # Entry: main.js (esbuild)
│   └── scss/         # Entry: main.scss (sass)
├── cmd/server/       # main.go
├── pkg/framework/    # Router, Context, SCSS/JS compilers
├── public/           # Built CSS/JS, served at /assets
├── views/            # Go html/template
├── Dockerfile
├── Makefile
└── .air.toml
```

## Prerequisites

- Go 1.25+
- [Dart Sass](https://sass-lang.com/install) (`npm install -g sass`)
- [esbuild](https://esbuild.github.io/) (`npm install -g esbuild`)
- [Air](https://github.com/air-verse/air) (`go install github.com/air-verse/air@latest`)

## Quick Start

1. Copy `.env.example` to `.env`:
   ```ini
   APP_NAME=gocore
   PORT=8080
   APP_ENV=development
   ```
2. Run:
   ```bash
   make dev
   ```
3. Open `http://localhost:8080`

## Make Targets

| Command          | Description                                      |
|------------------|--------------------------------------------------|
| `make dev`       | Air + asset watchers (development)               |
| `make run`       | Single run, assets compiled once                 |
| `make build`     | Binary to `bin/server`                           |
| `make docker-build` | Production Docker image (builds assets in image) |
| `make docker-run`   | Run container on port 8080                       |
| `make clean`     | Remove binary and image                          |

## Usage Example

```go
r := framework.NewRouter()
r.Use(framework.Logger)
r.Use(framework.Recovery)

r.LoadHTMLGlob("views/*.html")
r.Static("/assets", "./public")

r.GET("/", func(c *framework.Context) {
    c.HTML(200, "index.html", PageData{Title: "Home"})
})

r.POST("/api/data", func(c *framework.Context) {
    var input MyStruct
    if err := c.BindJSONStrict(&input); err != nil {
        c.JSON(400, map[string]string{"error": err.Error()})
        return
    }
    c.JSON(200, map[string]string{"status": "ok"})
})

r.SetNotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}))
```

## Docker

- **Build:** `make docker-build`
- **Run:** `make docker-run` (or `docker run -p 8080:8080 -e APP_ENV=production gocore`)

The Dockerfile compiles SCSS and JS in the builder stage; the final image only contains the binary, `views`, and `public`.

## License

MIT
