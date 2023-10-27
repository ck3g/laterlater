.PHONY: build_docker
build_docker:
	docker build --tag later-later-web .

.PHONY: run_docker
run_docker:
	docker run --rm -it -p 4000:4000/tcp later-later-web:latest

 
.PHONY: run
run:
	go run ./cmd/web

.PHONY: run_cli
run_cli:
	go run ./cmd/cli

.PHONY: test
test:
	go test ./...
