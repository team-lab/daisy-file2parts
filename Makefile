.PHONY: fmt
.PHONY: run
.PHONY: build

fmt:
	gofmt -w ./

run:
	go run ./file2parts.go

build:
	GOARCH=amd64 GOOS=windows go build -o ./build/windows-amd64/file2parts.exe file2parts.go
	GOARCH=amd64 GOOS=darwin go build -o ./build/darwin-amd64/file2parts file2parts.go
	GOARCH=amd64 GOOS=linux go build -o ./build/linux-amd64/file2parts file2parts.go
	GOARCH=386 GOOS=windows go build -o ./build/windows-386/file2parts.exe file2parts.go
	GOARCH=386 GOOS=darwin go build -o ./build/darwin-386/file2parts file2parts.go
	GOARCH=386 GOOS=linux go build -o ./build/linux-386/file2parts file2parts.go
