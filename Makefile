PROJECT_NAME=crossover
.DEFAULT_GOAL := run

test:
	go test ./...

build:
	CGO_ENABLED=0 go build -o ./target/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)/main.go

run:
	go run ./cmd/$(PROJECT_NAME)/main.go

e2e: build
	./target/$(PROJECT_NAME) --config ./e2e/config.yml & sleep 1 && go test ./e2e/... && killall $(PROJECT_NAME)
