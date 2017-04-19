.PHONY: fmt
.PHONY: run
.PHONY: build

fmt:
	gofmt -w ./

run:
	go run ./file2parts.go

build:
	if [ ! -e ./build/release ] ; \
	then \
	     mkdir ./build/release ; \
	fi;
	GOARCH=amd64 GOOS=windows go build -o ./build/windows-amd64/file2parts.exe file2parts.go
	zip ./build/release/windows_amd64.zip ./build/windows-amd64/file2parts.exe
	GOARCH=amd64 GOOS=darwin go build -o ./build/darwin-amd64/file2parts file2parts.go
	zip ./build/release/macos_amd64.zip ./build/darwin-amd64/file2parts
	GOARCH=amd64 GOOS=linux go build -o ./build/linux-amd64/file2parts file2parts.go
	zip ./build/release/linux_amd64.zip ./build/linux-amd64/file2parts
	GOARCH=386 GOOS=windows go build -o ./build/windows-386/file2parts.exe file2parts.go
	zip ./build/release/windows_386.zip ./build/windows-386/file2parts.exe
	GOARCH=386 GOOS=darwin go build -o ./build/darwin-386/file2parts file2parts.go
	zip ./build/release/macos_386.zip ./build/darwin-386/file2parts
	GOARCH=386 GOOS=linux go build -o ./build/linux-386/file2parts file2parts.go
	zip ./build/release/linux_386.zip ./build/linux-386/file2parts
