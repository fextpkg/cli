ifndef $(GOPATH)
	export GOPATH=$(shell go env GOPATH)
endif
ifndef $(GOBIN)
	export GOBIN=$(shell go env GOBIN)
endif

build:
	go build -i upip.go

deps:
	go get

install:
ifeq ($(OS),Windows_NT)
	mv upip.exe "$(USERPROFILE)\\"
else
ifeq ($(shell uname),Linux)
	mv upip ~/.local/bin/
endif
ifeq ($(shell uname),Darwin)
	mv upip usr/local/bin/
endif
endif

all: deps build install
