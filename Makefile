APP_NAME=gocore
PORT=8080

css-dev:
	@echo "Compiling SCSS (Dev)..."
	sass assets/scss/main.scss public/css/main.css

watch-css:
	@echo "Watching SCSS..."
	sass --watch assets/scss/main.scss:public/css/main.css

css-prod:
	@echo "Compiling SCSS (Production)..."
	sass --no-source-map --style=compressed assets/scss/main.scss public/css/main.css

run: css-dev
	@echo "Running locally..."
	go run cmd/server/main.go

dev:
	@echo "Starting Hot-Reload Server..."
	set APP_ENV=development && air

build: css-dev
	@echo "Building binary..."
	if not exist bin mkdir bin
	go build -o bin/server cmd/server/main.go

docker-build: css-prod
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

docker-run:
	@echo "Running Container on port $(PORT)..."
	docker run --rm -p $(PORT):$(PORT) -e APP_ENV=production --env-file .env $(APP_NAME)

clean:
	@echo "Cleaning up..."
	if exist bin rmdir /s /q bin
	docker rmi $(APP_NAME)
