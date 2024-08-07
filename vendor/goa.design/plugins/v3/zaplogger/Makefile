#! /usr/bin/make
#
# Makefile for Goa zaplogger plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
PLUGIN_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include ../plugins.mk

gen: clean
	goa gen goa.design/plugins/v3/zaplogger/examples/calc/design -o "$(PLUGIN_DIR)/examples/calc" && \
	make example

example:
	@ rm -rf "$(PLUGIN_DIR)/examples/calc/cmd" && \
	goa example goa.design/plugins/v3/zaplogger/examples/calc/design -o "$(PLUGIN_DIR)/examples/calc"

build-examples:
	@cd "$(PLUGIN_DIR)/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli

clean:
	@cd "$(PLUGIN_DIR)/examples/calc" && \
		rm -rf cmd calc.go
