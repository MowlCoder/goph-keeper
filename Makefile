.SILENT:
.PHONY:

client_build_version = v.0.0.1-beta
ifeq ($(OS), Windows_NT)
	build_date = $(shell date /t)
else
	build_date = $(shell date -I)
endif
client_logpass_secret = $(shell cat secrets/logpass)
client_card_secret = $(shell cat secrets/card)
client_ldflags = "-X main.buildDate=$(build_date) -X main.buildVersion=$(client_build_version) -X main.logPassSecret=$(client_logpass_secret) -X main.cardSecret=$(client_card_secret)"

build-server:
	go build -o ./bin/server ./cmd/server/main.go

run-server: build-server
	./bin/server

build-client:
	go build -ldflags $(client_ldflags) -o ./bin/client ./cmd/client/main.go

run-client: build-client
	./bin/client