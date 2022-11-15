DEVCONTAINER = devcontainer
GO = go
SRCS = main.go
BIN = devcontainer-shell

.PHONY: nop
nop:

.PHONY:
run:
	$(GO) run .

$(BIN): $(SRCS)
	$(DEVCONTAINER) exec --workspace-folder . env CGO_ENABLED=0 go build -o $@

.PHONY: devcontainer-build
devcontainer-build:
	$(DEVCONTAINER) --workspace-folder . build

.PHONY: devcontainer-up
devcontainer-up:
	$(DEVCONTAINER) --workspace-folder . up
