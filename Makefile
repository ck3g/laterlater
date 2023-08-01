.PHONY: run
run_cli:
	go run ./cmd/cli

.PHONY: test
test:
	go test ./...

