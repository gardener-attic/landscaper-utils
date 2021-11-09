# SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors.
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        := $(shell cat $(REPO_ROOT)/VERSION)
EFFECTIVE_VERSION                              := $(VERSION)-$(shell git rev-parse HEAD)

DISABLE_CLEANUP := false

.PHONY: install-requirements
install-requirements:
	@go install -mod=vendor $(REPO_ROOT)/vendor/github.com/ahmetb/gen-crd-api-reference-docs
	@go install -mod=vendor $(REPO_ROOT)/vendor/github.com/golang/mock/mockgen
	@go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
	@$(REPO_ROOT)/hack/install-requirements.sh

.PHONY: revendor
revendor:
	@$(REPO_ROOT)/hack/revendor.sh
	@echo "revendor machineimages"
	@cd $(REPO_ROOT)/machineimages && make revendor

.PHONY: check
check:
	@echo "check machineimages"
	@cd $(REPO_ROOT)/machineimages && make check

.PHONY: verify
verify: check

.PHONY: test
test:
	@echo "test machineimages"
	@cd $(REPO_ROOT)/machineimages && make test


