.PHONY: all clean-generated graphql-generate test notification-server

all: graphql-generate notification-server

graphql-generate: clean-generated
	@gqlgen

clean-generated:
	@rm -f generated.go
	@rm -f generated_resolver.go
	@rm -f models/generated.go

clean: clean-generated

test:
	@go test -v ./database

notification-server:
	@GOOS=linux go build ./cmd/notification-server

