# Image URL to use all building/pushing image targets
IMG ?= controller:dev

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Build manager binary
.PHONY: build
build: generate
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate
	go run ./main.go

# Build the docker image
docker: generate build
	docker rmi -f ${IMG}
	docker build -f Dockerfile -t ${IMG} .
	docker save -o ./bin/${IMG}.tar ${IMG}

# Install CRDs into a cluster
install: manifests bin/kustomize
	./bin/kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests bin/kustomize
	./bin/kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: docker manifests bin/kustomize
	cd config/manager && ../../bin/kustomize edit set image controller=${IMG}
	./bin/kustomize build config/default | kubectl apply -f -

# Generate code
generate: bin/controller-gen
	./bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Generate manifests e.g. CRD, RBAC etc.
manifests: bin/controller-gen
	./bin/controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Generate docs
.PHONY: docs
docs:
	./hack/plantuml.sh docs/*/*.plantuml

# Run tests
test: generate gofmt golint govet gosec unit manifests

# Run go fmt against code
gofmt:
	./hack/gofmt.sh

# Run go lint against code
golint: bin/golangci-lint
	./bin/golangci-lint run ./... --timeout=10m

# Run go vet against code
govet:
	go vet ./...

# Run go sec against code
gosec: bin/gosec
	./bin/gosec -quiet ./...

# Run go test against code
unit:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o coverage.html

# Clean files generated by scripts
clean:
	rm -f ./cover.out
	rm -f ./coverage.html
	rm -rf ./bin/*

# Clean up go module settings
mod:
	go mod tidy
	go mod verify

# Install kustomize
bin/kustomize:
	./hack/install_kustomize.sh

# Install controller-gen
bin/controller-gen:
	./hack/install_controller-gen.sh

# Install golangci-lint
bin/golangci-lint:
	./hack/install_golangci-lint.sh

# Install gosec
bin/gosec:
	./hack/install_gosec.sh
