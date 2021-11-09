export NAME=rego-example

all: build
clean:
	rm -f ecal

mod:
	go mod init || true
	go mod tidy

fmt:
	gofmt -l -w -s .

vet:
	go vet ./...

build: clean mod fmt vet
	go build -ldflags "-s -w" -o $(NAME) *.go

