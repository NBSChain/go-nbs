SHELL=PATH='$(PATH)' /bin/sh

PLATFORM := $(shell uname -o)


EXTEND := .exe
ifeq ($(PLATFORM), Msys)
	EXTEND := .exe
else ifeq ($(PLATFORM), Cygwin)
	EXTEND := .exe
else
	EXTEND :=
endif

EXENAME := natTest$(EXTEND)

all:
	go build -race -o $(EXENAME)