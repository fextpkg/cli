ifndef $(GOPATH)
	export GOPATH=$(shell go env GOPATH)
endif
ifndef $(GOBIN)
	export GOBIN=$(shell go env GOBIN)
endif

build:
	go build -i -o fext-bin ./fext/fext.go

deps:
	go get ./fext

install:
ifeq ($(OS),Windows_NT)
	mv fext-bin "$(USERPROFILE)\\fext.exe"
else
ifeq ($(shell uname),Linux)
	mv fext-bin ~/.local/bin/fext
endif
ifeq ($(shell uname),Darwin)
	mv fext-bin usr/local/bin/fext
endif
endif

all: deps build install
