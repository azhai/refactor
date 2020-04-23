BINNAME=reverse
RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=1
GOBUILD=$(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"

.PHONY: all build clean pre upx upxx

all: pre clean build
build:
	@echo "Compile $(BINNAME) ..."
	GOOS=$(GOOS) $(GOBUILD) -mod=vendor -o $(BINNAME) ./cmd/
	@echo "Build success."
clean:
	rm -f $(BINNAME)
	@echo "Clean all."
pre:
	$(GOBIN) mod tidy && $(GOBIN) mod vendor
upx: build command
	$(UPXBIN) $(BINNAME)
upxx: build command
	$(UPXBIN) --ultra-brute $(BINNAME)
