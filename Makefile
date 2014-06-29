GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

all: run

run:
	$(GO) run api.go
