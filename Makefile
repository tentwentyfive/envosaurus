TAG?=""

.DEFAULT_GOAL := test

# Run all tests
.PHONY: test
test: fmt lint vet test-unit go-mod-tidy

# Download dependencies
.PHONY: mod-download
mod-download:
	go mod download

# Run unit tests
.PHONY: test-unit
test-unit: mod-download
	gotestsum --format testname -- -race ./...

# Clean go.mod
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy
	git diff --exit-code go.sum

# Check formatting
.PHONY: fmt
fmt:
	test -z "$(shell gofmt -l .)"

# Run linter
.PHONY: lint
lint:
	golint -set_exit_status ./...

# Run vet
.PHONY: vet
vet:
	go vet ./...

# Run a test release with goreleaser
.PHONY: test-release
test-release:
	goreleaser --snapshot --skip-publish --rm-dist

# Clean up any cruft left over from old builds
.PHONY: clean
clean:
	rm -rf envosaurus dist/

# Build the application
.PHONY: build
build: clean
	CGO_ENABLED=0 go build ./cmd/envosaurus/envosaurus.go

# Create a git tag
.PHONY: tag
tag:
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

# Requires GITHUB_TOKEN environment variable to be set
.PHONY: release
release: clean
	goreleaser
