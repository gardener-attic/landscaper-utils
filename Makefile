# SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors.
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        := $(shell cat $(REPO_ROOT)/VERSION)
EFFECTIVE_VERSION                              := $(VERSION)-$(shell git rev-parse HEAD)
REGISTRY                                       := eu.gcr.io/gardener-project/landscaper-utils
MACHINE_IMAGES_IMAGE_REPOSITORY                := $(REGISTRY)/machine-images
LANDSCAPER_UTILS_IMAGE_REPOSITORY              := $(REGISTRY)/landscaper-utils
COMP_NAME_MACHINE_IMAGES                       := machine-images


.PHONY: install-requirements
install-requirements:
	@curl -sfL "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" | sh -s -- -b $(go env GOPATH)/bin v1.32.2
	@GO111MODULE=off go get golang.org/x/tools/cmd/goimports

.PHONY: revendor
revendor:
	@GO111MODULE=on go mod vendor
	@GO111MODULE=on go mod tidy

.PHONY: format
format:
	@goimports -l -w -local=github.com/gardener/image-vector ./pkg ./cmd

.PHONY: check
check:
	@$(REPO_ROOT)/hack/check.sh --golangci-lint-config=./.golangci.yaml $(REPO_ROOT)/pkg/... $(REPO_ROOT)/cmd/...

.PHONY: test
test:
	@go test ./pkg/... ./cmd/...

.PHONY: verify
verify: check test

.PHONY: install
install:
	@EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) ./hack/install.sh

.PHONY: docker-images
docker-images:
	@echo "Building docker images for version $(EFFECTIVE_VERSION)"
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(MACHINE_IMAGES_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile --target machine-images .
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(LANDSCAPER_UTILS_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-utils .

.PHONY: docker-push
docker-push:
	@echo "Pushing docker images for version $(EFFECTIVE_VERSION) to registry $(REGISTRY)"
	@if ! docker images $(MACHINE_IMAGES_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(MACHINE_IMAGES_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@if ! docker images $(LANDSCAPER_UTILS_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(LANDSCAPER_UTILS_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@docker push $(MACHINE_IMAGES_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)
	@docker push $(LANDSCAPER_UTILS_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)

cnudie: export COMPONENT_NAME=$(COMP_NAME_MACHINE_IMAGES)
.PHONY: cnudie
cnudie:
	@EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) ./hack/generate-cd.sh

tester: export COMPONENT_NAME=$(COMP_NAME_MACHINE_IMAGES)
.PHONY: tester
tester:
	@echo "Building docker images for version $(COMPONENT_NAME) "
