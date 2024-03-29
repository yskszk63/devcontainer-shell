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
	CGO_ENABLED=0 go build -o=$@ ./bin/devcontainer-shell

.PHONY: test cov.out
test: cov.html

cov.out:
	go test -v -coverpkg . -coverprofile $@

cov.html: cov.out
	go tool cover -html=$< -o $@
