.PHONY: all ahex aqua-explorer

all:
	GOBIN=${PWD} go install ./cmd/...

ahex:
	GOBIN=${PWD} go install ./cmd/ahex

aqua-explorer:
	GOBIN=${PWD} go install ./cmd/aqua-explorer

