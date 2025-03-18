.DEFAULT_GOAL := build

build:
	@go mod tidy
	@go build -o ./bin/microblog .

run: build
	./bin/orbit
