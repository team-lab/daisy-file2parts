.PHONY: fmt
.PHONY: run
.PHONY: build

fmt:
	gofmt -w ./

run:
	go run ./file2parts.go

build:
	GOARTCH=amd64 GOOS=windows go build -o ./build/windows-amd64/file2parts.exe file2parts.go
	GOARTCH=amd64 GOOS=darwin go build -o ./build/darwin-amd64/file2parts file2parts.go
	GOARTCH=amd64 GOOS=linux go build -o ./build/linux-amd64/file2parts file2parts.go
