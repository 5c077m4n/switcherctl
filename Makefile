.PHONY: start-dev
start-dev:
	go run .

.PHONY: test-all
test-all:
	go test -v ./...

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: update-deps
update-deps: mod-tidy
	go get -u ./...
