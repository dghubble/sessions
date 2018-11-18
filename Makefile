.PHONY: all
all: test vet lint fmtcheck

.PHONY: test
test:
	@go test ./...

.PHONY: vet
vet:
	@go vet -all -shadow ./...

.PHONY: lint
lint:
	@golint -set_exit_status ./...

.PHONY: fmt
fmt:
	@gofmt -l -s -w .

.PHONY: fmtcheck
fmtcheck:
	@gofmt -l -s .
	@test -z $$(gofmt -l -s .)
