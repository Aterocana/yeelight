PROJECTNAME=$(shell basename "$(PWD)")
# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)
PKGS=$(shell go list ./... | grep -v /vendor)

test:
	go test $(PKGS) -cover

clean:
	@cd $(GOBASE)/cover && ls | grep -v .gitkeep | xargs rm && cd $(GOBASE)