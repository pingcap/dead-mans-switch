.PHONY: build-linux build docker clean

EXE = $(GOPATH)/bin/dms

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-s -w -static"' -o $(EXE) .

build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-s -w -static"' -o $(EXE) .

docker:
	docker build -t deadmansswitch .

clean:
	rm -f $(EXE)
