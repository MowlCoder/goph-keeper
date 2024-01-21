.SILENT:
.PHONY:

client_build_version = v.0.0.1-beta

ifeq ($(OS), Windows_NT)
	build_date = $(shell date /t)

	client_logpass_secret = $(shell type .\secrets\logpass)
	client_card_secret = $(shell type .\secrets\card)

	client_binary_path = ./bin/client.exe
	server_binary_path = ./bin/server.exe
else
	build_date = $(shell date -I)

	client_logpass_secret = $(shell cat secrets/logpass)
	client_card_secret = $(shell cat secrets/card)

	client_binary_path = ./bin/client
	server_binary_path = ./bin/server.exe
endif

client_ldflags = "-X main.buildDate=$(build_date) -X main.buildVersion=$(client_build_version) -X main.logPassSecret=$(client_logpass_secret) -X main.cardSecret=$(client_card_secret)"

build-server:
	go build -o $(server_binary_path) ./cmd/server/main.go

run-server: build-server
	$(server_binary_path)

build-client:
	go build -ldflags $(client_ldflags) -o $(client_binary_path) ./cmd/client/main.go

run-client: build-client
	$(client_binary_path)

fmt:
	gofmt -s -w .
	goimports -w .