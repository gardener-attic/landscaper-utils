# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

#### BUILDER ####
FROM golang:1.16.5 AS builder

WORKDIR /go/src/github.com/gardener/landscaper-utils
COPY . .

ARG EFFECTIVE_VERSION

RUN make install EFFECTIVE_VERSION=$EFFECTIVE_VERSION

#### BASE ####
FROM eu.gcr.io/gardenlinux/gardenlinux:184.0 AS base

#### Landscaper Utils ####
FROM base as landscaper-utils

COPY --from=builder /go/bin/landscaper-utils /landscaper-utils

WORKDIR /

ENTRYPOINT ["/landscaper-utils"]

#### Machine Images ####
FROM base as machine-images

COPY --from=builder /go/bin/machine-images /machine-images

WORKDIR /

ENTRYPOINT ["/machine-images"]