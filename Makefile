## Print the help message.
# Parses this Makefile and prints targets that are preceded by "##" comments.
help:
	@echo "" >&2
	@echo "Available targets: " >&2
	@echo "" >&2
	@awk -F : '\
			BEGIN { in_doc = 0; } \
			/^##/ && in_doc == 0 { \
				in_doc = 1; \
				doc_first_line = $$0; \
				sub(/^## */, "", doc_first_line); \
			} \
			$$0 !~ /^#/ && in_doc == 1 { \
				in_doc = 0; \
				if (NF <= 1) { \
					next; \
				} \
				printf "  %-15s %s\n", $$1, doc_first_line; \
			} \
			' <"$(abspath $(lastword $(MAKEFILE_LIST)))" \
		| sort >&2
	@echo "" >&2

build:
	go generate
	go build

## builds the binary and deploys a monolith instance to local machine
monolith.run: build
	weaver single deploy weaver.toml

## builds the binary and deploys each component as a microservice on the local machine
services.run: build
	weaver multi deploy weaver.toml

## starts up the dashboard to view diagnostics and metrics for the services
services.dashboard:
	weaver multi dashboard

## clean built binary and generated service weaver files if the exist
clean:
	rm kafkaless
	find . -name '*_gen.go' -delete