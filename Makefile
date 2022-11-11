BINARY_NAME=cinnox-homework

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

run:
	./${BINARY_NAME}-linux

build_and_run: build run

test:
	go test ./...

dep:
 go mod download

dev_env_up:
	docker compose up -d

dev_env_down:
	docker compose down

clean:
	go clean
	rm ${BINARY_NAME}-linux