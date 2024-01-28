.SILENT:
.PHONY:

client_build_version = v.0.0.1-beta

ifeq ($(OS), Windows_NT)
	build_date = $(shell date /t)

	client_data_secret = $(shell type .\secrets\data)

	client_binary_path = ./bin/client.exe
	server_binary_path = ./bin/server.exe
else
	build_date = $(shell date -I)

	client_data_secret = $(shell cat secrets/data)

	client_binary_path = ./bin/client
	server_binary_path = ./bin/server.exe
endif

client_ldflags = "-X main.buildDate=$(build_date) -X main.buildVersion=$(client_build_version) -X main.dataSecret=$(client_data_secret)"

build-server:
	go build -o $(server_binary_path) ./cmd/server/main.go

run-server: build-server
	$(server_binary_path)

build-client:
	go build -ldflags $(client_ldflags) -o $(client_binary_path) ./cmd/client/main.go

run-client: build-client
	$(client_binary_path)

test:
	go test ./...

mocks:
	mockgen -source=./internal/services/client/user_stored_data.go -destination=./internal/services/client/mocks/user_stored_data.go
	mockgen -source=./internal/services/server/user_stored_data.go -destination=./internal/services/server/mocks/user_stored_data.go
	mockgen -source=./internal/services/server/user.go -destination=./internal/services/server/mocks/user.go
	mockgen -source="./internal/handlers/user.go" -destination="./internal/handlers/mocks/user.go"
	mockgen -source="./internal/handlers/user_stored_data.go" -destination="./internal/handlers/mocks/user_stored_data.go"

fmt:
	gofmt -s -w .
	goimports -w .