
LDFLAGS =

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
VERSION    = $(shell echo "${GITHUB_REF}" | cut -d / -f 3)

ifdef GITHUB_ACTIONS
		GIT_DIRTY = clean
endif

LDFLAGS += -X github.com/qrtt1/friendly-yaml/internal/flatyaml.tagVersion=${VERSION}
LDFLAGS += -X github.com/qrtt1/friendly-yaml/internal/flatyaml.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/qrtt1/friendly-yaml/internal/flatyaml.gitTreeState=${GIT_DIRTY}
LDFLAGS += $(EXT_LDFLAGS)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


all: app


# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Build friendly-yaml binary
app: fmt vet
	go build -o helm-values -ldflags '$(LDFLAGS)' cmd/main.go

# Run usage-agnet
run: fmt vet
	go run ./cmd/main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...
