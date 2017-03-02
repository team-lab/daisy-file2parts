.PHONY: fmt
.PHONY: run

fmt:
	gofmt -w ./

run:
	go run ./file2parts.go
