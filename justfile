# Git Branch Grouper — Development Recipes
# Run `just` to see all available commands

set shell := ["powershell", "-NoProfile", "-Command"]

# Module and binary config
binary := "git-branch-grouper"
module := "git-branch-grouper-plugin"
entry  := "./cmd/git-branch-grouper"

# Get version from latest git tag (without 'v' prefix), fallback to "dev"
version := `git describe --tags --always 2>$null | ForEach-Object { $_ -replace '^v', '' }; if (-not $?) { echo dev }`

# ─── Build ────────────────────────────────────────────────────────────────────

# Build binary for current OS
[group("Build")]
build:
    go build -trimpath -ldflags "-s -w -X main.version={{version}}" -o {{binary}}.exe {{entry}}

# Build and run
[group("Build")]
run *ARGS: build
    .\{{binary}}.exe {{ARGS}}

# Install binary to $GOPATH/bin
[group("Build")]
install *ARGS:
    go install -trimpath -ldflags "-s -w -X main.version={{version}}" {{entry}}@latest
    @echo Installed {{binary}} to $env:GOPATH\bin

# Remove installed binary from $GOPATH/bin
[group("Build")]
uninstall:
    go clean -i {{entry}}

# Remove build artifacts
[group("Build")]
clean:
    go clean
    rm -Force {{binary}}.exe -ErrorAction SilentlyContinue
    rm -Recurse -Force dist/ -ErrorAction SilentlyContinue
    rm -Force coverage.html, coverage.out, profile.cov, bench-new.txt -ErrorAction SilentlyContinue

# ─── Quality ──────────────────────────────────────────────────────────────────

# Run all checks: format, vet, lint, test
[group("Quality")]
check: fmt vet lint test

# Format all Go files
[group("Quality")]
fmt:
    gofumpt -l -w .

# Run go vet
[group("Quality")]
vet:
    go vet ./...

# Run golangci-lint (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
[group("Quality")]
lint:
    golangci-lint run ./...

# ─── Test ─────────────────────────────────────────────────────────────────────

# Run all tests
[group("Test")]
test:
    go test ./...

# Run tests with race detector
[group("Test")]
test-race:
    go test -race ./...

# Run tests with coverage report (opens HTML in browser)
[group("Test")]
cover:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report saved to coverage.html"
    start coverage.html

# Run benchmarks (10s each)
[group("Test")]
bench:
    go test -bench=. -benchmem -run=^$ ./...

# Run benchmarks with multiple iterations for benchstat comparison
[group("Test")]
bench-compare:
    go test -bench=. -benchmem -count=6 -run=^$ ./... > bench-new.txt
    @echo Results saved to bench-new.txt. Run: benchstat bench-new.txt

# ─── Version ──────────────────────────────────────────────────────────────────

# Print current version (from git tag)
[group("Version")]
version:
    @echo {{version}}

# Bump patch version (v1.0.0 → v1.0.1), create tag, print result
[group("Version")]
patch:
    powershell -NoProfile -Command "& { \
        $tag = (git describe --tags --abbrev=0 2>$null); \
        if (-not $tag) { $tag = 'v0.0.0' }; \
        $parts = $tag -split '\.'; \
        $major = [int]$parts[0].TrimStart('v'); \
        $minor = [int]$parts[1]; \
        $patch = [int]$parts[2] + 1; \
        $new = \"v$major.$minor.$patch\"; \
        git tag $new; \
        Write-Host \"Created tag: $new\" \
    }"

# Bump minor version (v1.0.0 → v1.1.0), create tag, print result
[group("Version")]
minor:
    powershell -NoProfile -Command "& { \
        $tag = (git describe --tags --abbrev=0 2>$null); \
        if (-not $tag) { $tag = 'v0.0.0' }; \
        $parts = $tag -split '\.'; \
        $major = [int]$parts[0].TrimStart('v'); \
        $minor = [int]$parts[1] + 1; \
        $new = \"v$major.$minor.0\"; \
        git tag $new; \
        Write-Host \"Created tag: $new\" \
    }"

# Bump major version (v1.0.0 → v2.0.0), create tag, print result
[group("Version")]
major:
    powershell -NoProfile -Command "& { \
        $tag = (git describe --tags --abbrev=0 2>$null); \
        if (-not $tag) { $tag = 'v0.0.0' }; \
        $parts = $tag -split '\.'; \
        $major = [int]$parts[0].TrimStart('v') + 1; \
        $new = \"v$major.0.0\"; \
        git tag $new; \
        Write-Host \"Created tag: $new\" \
    }"

# Push all tags to remote
[group("Version")]
push-tags:
    git push --tags

# ─── Release ──────────────────────────────────────────────────────────────────

# Run GoReleaser in snapshot mode (local build, no publish)
[group("Release")]
release-snapshot:
    goreleaser release --snapshot --clean

# Run GoReleaser full release (typically triggered by CI on tag push)
[group("Release")]
release:
    goreleaser release --clean

# ─── Dev ──────────────────────────────────────────────────────────────────────

# Update Go module dependencies
[group("Dev")]
tidy:
    go mod tidy

# Verify dependency checksums
[group("Dev")]
verify:
    go mod verify

# Update all dependencies to latest minor/patch versions
[group("Dev")]
update-deps:
    go get -u ./...
    go mod tidy

# Show available recipes
[group("Dev")]
default:
    @just --list
