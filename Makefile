.PHONY: all
all: test vet lint checkfmt

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

.PHONY: checkfmt
checkfmt:
	@gofmt -l -s .
	@test -z $$(gofmt -l -s .)
