NAME=example

.PHONY: build
## build: Compile the packages.
build:
	go build -o ./bin/$(NAME)

.PHONY: run
## run: Run the application
run:
	@go run main.go --config.file=config.yaml

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f ./bin

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
.PHONY:lint
lint: 
	golangci-lint run

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo