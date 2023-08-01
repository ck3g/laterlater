
.PHONY: run
run:
	go run ./cmd/web

.PHONY: run_cli
run_cli:
	go run ./cmd/cli

.PHONY: test
test:
	go test ./...

