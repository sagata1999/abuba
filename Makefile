.PONY: all check arm64 clean

NAME = joker

RM = $(shell which rm)
VERSION  = $(shell git describe)
BTIME    = $(shell env LANG=C date -u +'%b_%d_%Y_%H:%M:%S')
BSYSTEM  = $(shell uname -srv | sed 's/ /_/g')
COMPILER = $(shell go version | sed 's/ /_/g')

MODULE = $(shell awk '$$1 ~ /module/ {print $$2}' go.mod)
VERSIONPATH = /internal/pkg/version

FLAGS = -trimpath

LDPARAM = -s -v
LDPARAM += -X $(MODULE)$(VERSIONPATH).Version=$(VERSION)
LDPARAM += -X $(MODULE)$(VERSIONPATH).BuildTime=$(BTIME)
LDPARAM += -X $(MODULE)$(VERSIONPATH).System=$(BSYSTEM)
LDPARAM += -X $(MODULE)$(VERSIONPATH).Compiler=$(COMPILER)

LDFLAGS =-ldflags "$(LDPARAM)"

TARGET = -o bin/$(NAME)

SRC = ./cmd/$(NAME)

all: clean check arm64

check:
	@echo "*** check stage ***"
	@go-consistent -v ./...
	@echo "*** check complete ***\n"

arm64:
	@echo "*** compile ***"
	GOOS=linux GOARCH=arm64 go build $(FLAGS) $(LDFLAGS) $(TARGET) $(SRC)
	@echo "*** compile complete ***\n"

clean:
	@echo "*** clean stage ***"
	rm -rf bin/*
	@echo "*** clean complete ***\n"
