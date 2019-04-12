.PHONY: all clean-generated graphql-generate test notification-server requester

all: graphql-generate notification-server requester

graphql-generate: clean-generated
	@gqlgen

clean-generated:
	@rm -f generated.go
	@rm -f generated_resolver.go
	@rm -f models/generated.go

clean: clean-generated

test:
	@go test -v ./database

# binaries

requester:
	@go build ./cmd/requester

notification-server:
	@GOOS=linux go build ./cmd/notification-server

