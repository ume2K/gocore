# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git nodejs npm
RUN npm install -g sass esbuild

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN sass assets/scss/main.scss public/css/main.css --style=compressed
RUN esbuild assets/js/main.js --bundle --minify --outfile=public/js/main.js

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# --- Stage 2: Runtime ---
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/server /server
COPY --from=builder /app/views /views
COPY --from=builder /app/public /public

USER 10001

EXPOSE 8080

ENTRYPOINT ["/server"]
