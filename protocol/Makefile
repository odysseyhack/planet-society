.PHONY: all clean-generated graphql-generate tesd

all: graphql-generate

graphql-generate: clean-generated
	@gqlgen

clean-generated:
	@rm -f generated.go
	@rm -f generated_resolver.go
	@rm -f models/generated.go

clean: clean-generated

test:
	@go test -v ./database

