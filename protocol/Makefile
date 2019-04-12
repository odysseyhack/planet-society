.PHONY: all clean-generated graphql-generate test notification-server requester responder

all: graphql-generate binaries

graphql-generate: clean-generated
	@gqlgen

clean-generated:
	@rm -f generated.go
	@rm -f generated_resolver.go
	@rm -f models/generated.go

clean: clean-generated

test:
	go test -v ./database
	go test -v ./cryptography

# binaries
binaries: requester responder notification-server

requester:
	@go build ./cmd/requester

responder:
	@go build ./cmd/responder

notification-server:
	@GOOS=linux go build ./cmd/notification-server

