SHELL := bash

MKDIR ?= mkdir
GO    ?= go
YARN  ?= yarn

.EXPORT_ALL_VARIABLES:
.PHONY: deps

deps:
	$(YARN) install

run:
	$(YARN) build
	$(GO) run main.go server

