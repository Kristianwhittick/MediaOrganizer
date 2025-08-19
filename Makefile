.PHONY: build clean all

all: build

build: build-windows build-linux build-arm

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/media-organizer-windows-x64.exe ./cmd

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/media-organizer-linux-x64 ./cmd

build-arm:
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/media-organizer-linux-arm ./cmd

clean:
	rm -rf bin/

deps:
	go mod tidy