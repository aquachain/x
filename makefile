BUILD_DIR = ${PWD}/build/bin
COMMAND_DIR = ./cmd/...
.PHONY: all

all:
	GOBIN=${BUILD_DIR} go install ${COMMAND_DIR}

clean:
	@test -f ${BUILD_DIR}/README.md && echo "bad $$BUILD_DIR variable" && exit 111 || true
	@test -f ${BUILD_DIR}/bin && echo "bad $$BUILD_DIR variable" && exit 111 || true
	@test -d ${BUILD_DIR} && rm -fvr ${BUILD_DIR} || true

