
.PHONY: build
build:
	gofumpt -l -w .
	go mod tidy
	./build.sh

.PHONY: run
run:
	./build.sh
	./run.sh

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: test
test:
	go test -gcflags=all=-l ./...