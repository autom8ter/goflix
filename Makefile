.PHONY: proto
proto: ## generate protobufs
	@docker run -v `pwd`:/defs colemanword/prototool:1.17_0 generate

.PHONY: help
help:	## show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'