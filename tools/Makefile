TOOLS := $(PWD)/tools/bin
PATH  := $(TOOLS):$(PATH)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	find . -type f \( -name "*.go" ! -path "$(PWD)/vendor/*" \) \
		-exec goimports -w -local sprapp {} \;

.PHONY: tools
tools:
	go mod init tools || true
	GOBIN=$(TOOLS) go generate $(PWD)/tools

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify
