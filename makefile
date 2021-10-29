# Get version from git hash
git_hash := $(shell git rev-parse --short HEAD || echo 'development')

# project version
# version = $(shell git describe --tags --abbrev=0 || echo 'development')
# version = $(shell git tag | sort -V | tail -1 || echo 'development')
version = "dev"
# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

name:="wallectl"

# Add linker flags
linker_flags = '-s -X github.com/piqba/wallertme/cmd/cli/commands.buildTime=${current_time} -X github.com/piqba/wallertme/cmd/cli/commands.versionHash=${git_hash} -X github.com/piqba/wallertme/cmd/cli/commands.version=${version}'

# Build binaries for current OS and Linux
.PHONY:
compile:
	@echo "Building binaries..."

	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o ./build/${name}-${version}-linux-amd64 cmd/cli/main.go


compress:
	./upx -9 -q ./build/${name}-${version}-linux-amd64

# Only work if you have goreleaser instaled
releaser:
	goreleaser release --snapshot --rm-dist

# Run redis for local env
redis:
	podman run --name redis --rm -e ALLOW_EMPTY_PASSWORD=yes -p 6379:6379 quay.io/bitnami/redis:latest