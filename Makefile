DEVCONTAINER = devcontainer
GO = go
SRCS = main.go
BIN = devcontainer-shell

.PHONY: nop
nop:

.PHONY:
run:
	$(GO) run .

.PHONY: devcontainer-build
devcontainer-build:
	$(DEVCONTAINER) --workspace-folder . build

.PHONY: devcontainer-up
devcontainer-up:
	$(DEVCONTAINER) --workspace-folder . up



.PHONY: devcontainer-shell
devcontainer-shell:
	go build -a -tags netgo -installsuffix netgo -ldflags='-s -w -extldflags "-static"' -o=$@ ./bin/devcontainer-shell

.PHONY: devcontainer-shell-agent
devcontainer-shell-agent:
	go build -a -tags netgo -installsuffix netgo -ldflags='-s -w -extldflags "-static"' -o=$@ ./bin/devcontainer-shell-agent

.PHONY: build-devcontainer-shell-agent-container
build-devcontainer-shell-agent-container:
	docker buildx build -f docker/Dockerfile.agent . -t ghcr.io/yskszk63/devcontainer-shell-agent
