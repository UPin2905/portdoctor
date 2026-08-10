.PHONY: build test fmt vet clean

build:
	go build -o portdoctor.exe ./cmd/portdoctor

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test build

clean:
	del /f portdoctor.exe 2>nul || true
