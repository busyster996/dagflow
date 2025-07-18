PACKAGE_NAME	:= github.com/busyster996/dagflow
GOLANG_VERSION	?= latest

.PHONY: all
all: binary copy-binary
	@sha256sum bin/dagflow* > bin/latest.sha256sum

.PHONY: dev
dev: # generate
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" --tags "codec.notfastpath" -o bin/dagflow-base_linux_amd64_v1 cmd/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" --tags "codec.notfastpath,dagflow.all_runner_features" -o bin/dagflow_linux_amd64_v1 cmd/main.go

swag:
	@swag init --exclude pkg --parseDependencyLevel 3 --dir internal/server/router --outputTypes json -g router.go

# Run code generation
generate:
	@echo "Tidying up Go modules..."
	@go mod tidy
	@echo "Running go generate..."
	@go generate ./...

.PHONY: binary
binary:
	@echo "Building the binary..."
	@rm -fr $(CURDIR)/dist
	@docker run \
		--rm \
		--privileged \
		--network host \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(CURDIR):/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GOLANG_VERSION} \
		release -p 3 --timeout 3h --clean --auto-snapshot --snapshot --skip=chocolatey,docker,homebrew,publish,scoop,validate,winget

.PHONY: copy-binary
copy-binary:
	@echo "Copying binaries..."
	@rm -fr $(CURDIR)/bin
	@mkdir -p $(CURDIR)/bin
	@find $(CURDIR)/dist/dagflow* -type f -not -path "*checksums*" -exec bash -c 'cp -f {} $(CURDIR)/bin/`echo {}|sed "s|$(CURDIR)/dist/||g"|sed "s|/dagflow||g"`' \;
	@rm -fr dist