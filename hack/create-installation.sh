#!/usr/bin/env bash
#
# Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0

set -e

SOURCE_PATH="$(dirname $0)/.."
TMP_DIR="$(mktemp -d)"
INSTALLATION_PATH="${TMP_DIR}/installation.yaml"

echo "Generation an installation"
echo "Registry:          ${REGISTRY}"
echo "Component name:    ${COMPONENT_NAME}"
echo "Effective version: ${EFFECTIVE_VERSION}"

cat << EOF > ${INSTALLATION_PATH}
apiVersion: landscaper.gardener.cloud/v1alpha1
kind: Installation
metadata:
  name: machine-images
spec:
  componentDescriptor:
    ref:
      repositoryContext:
        type: ociRegistry
        baseUrl: ${REGISTRY}
      componentName: github.com/gardener/landscaper-utils/${COMPONENT_NAME}
      version: ${EFFECTIVE_VERSION}

  blueprint:
    ref:
      resourceName: blueprint

  importDataMappings:
    machineImages:
      - name: gardenlinux
        versions:
          - classification: preview
            cri:
              - name: docker
              - containerRuntimes:
                  - type: gvisor
                name: containerd
            version: 318.9.0
          - classification: supported
            cri:
              - name: docker
              - containerRuntimes:
                  - type: gvisor
                name: containerd
            version: 318.8.0
          - classification: deprecated
            cri:
              - name: docker
              - containerRuntimes:
                  - type: gvisor
                name: containerd
            expirationDate: '2022-01-15T23:59:59Z'
            version: 184.0.0
      - name: ubuntu
        versions:
          - classification: supported
            cri:
              - name: docker
              - containerRuntimes:
                  - type: gvisor
                name: containerd
            version: 18.4.20210415
          - classification: deprecated
            cri:
              - name: docker
              - containerRuntimes:
                  - type: gvisor
                name: containerd
            expirationDate: '2021-11-30T23:59:59Z'
            version: 18.4.20200228

    machineImagesLs:
      - name: flatcar
        versions:
          - version: 2765.2.6
            classification:  supported
            cri:
              - name: containerd
                containerRuntimes:
                  - type: gvisor
              - name: docker

    machineImagesProvider:
      - name: gardenlinux
        versions:
          - regions:
              - ami: ami-0a47acd1c0f3b1f57
                name: eu-north-1
              - ami: ami-0d3094fa1189c0ea5
                name: ap-south-1
              - ami: ami-04d610c64abcbb153
                name: eu-west-3
              - ami: ami-05d4646ec6cf0e3e2
                name: eu-west-2
              - ami: ami-07b4af86edec6fc93
                name: eu-west-1
              - ami: ami-0620ba3dc3de56622
                name: ap-northeast-3
              - ami: ami-0a75b30c2c3fdb2d8
                name: ap-northeast-2
              - ami: ami-0a0883043d5a4308a
                name: ap-northeast-1
              - ami: ami-00116e4fa9cac55d3
                name: sa-east-1
              - ami: ami-0f96d8d78708728ff
                name: ca-central-1
              - ami: ami-013e6400712e60c8d
                name: ap-southeast-1
              - ami: ami-0d13dfd2312518f83
                name: ap-southeast-2
              - ami: ami-0f1599beff38ba139
                name: us-east-1
              - ami: ami-08201b4e691239e63
                name: us-east-2
              - ami: ami-0758cf6a37c0b5d18
                name: us-west-1
              - ami: ami-0fa01a63dbe2de931
                name: us-west-2
              - ami: ami-042e1a36006be30e0
                name: eu-central-1
            version: 318.9.0
          - regions:
              - ami: ami-077967ddf49bd7822
                name: eu-north-1
              - ami: ami-07f7b310124c898e1
                name: ap-south-1
              - ami: ami-00e17f38ab1ba55f4
                name: eu-west-3
              - ami: ami-0774a6896e4ac35e0
                name: eu-west-2
              - ami: ami-032de7308a9eb91f5
                name: eu-west-1
              - ami: ami-08e8b1a5ad052751f
                name: ap-northeast-3
              - ami: ami-0fc6badc68b7b7dd1
                name: ap-northeast-2
              - ami: ami-06a5e5183080d4345
                name: ap-northeast-1
              - ami: ami-0f4797d1ea0ffbe21
                name: sa-east-1
              - ami: ami-0320e4ff5c5004a6b
                name: ca-central-1
              - ami: ami-041d6354bb257235d
                name: ap-southeast-1
              - ami: ami-0ef9fbc659adf4a7e
                name: ap-southeast-2
              - ami: ami-066eb78156cb8e30d
                name: us-east-1
              - ami: ami-08782642e97383550
                name: us-east-2
              - ami: ami-088a63ebb9b6d8cca
                name: us-west-1
              - ami: ami-0d3df510f088f6728
                name: us-west-2
              - ami: ami-0b8eda557039b448e
                name: eu-central-1
            version: 318.8.0
          - regions:
              - ami: ami-0a7380651a94f34af
                name: eu-north-1
              - ami: ami-0f32041348763f71c
                name: ap-south-1
              - ami: ami-026f095584d001426
                name: eu-west-3
              - ami: ami-096b4e0e74268868e
                name: eu-west-2
              - ami: ami-06672bb04fdda3653
                name: eu-west-1
              - ami: ami-0f747a4612321e16d
                name: ap-northeast-2
              - ami: ami-064f4f4a6ae883ef1
                name: ap-northeast-1
              - ami: ami-0cce0186c61c3f502
                name: sa-east-1
              - ami: ami-09d6ef5e7a53a9972
                name: ca-central-1
              - ami: ami-0bbb464d2bad84aa4
                name: ap-southeast-1
              - ami: ami-084302575ffe36a40
                name: ap-southeast-2
              - ami: ami-0e5b90e9d0988a82d
                name: eu-central-1
              - ami: ami-03a6f7727025a545d
                name: us-east-1
              - ami: ami-09d1a0dcdd315fd22
                name: us-east-2
              - ami: ami-041b905f7b38e8e9b
                name: us-west-1
              - ami: ami-0d1ee9478903ca81f
                name: us-west-2
              - ami: ami-02a820cd1f31fa728
                name: us-gov-west-1
              - ami: ami-0e05a9dbe093f2152
                name: us-gov-east-1
            version: 184.0.0
      - name: ubuntu
        versions:
          - regions:
              - ami: ami-0afad43e7d620260c
                name: eu-north-1
              - ami: ami-04bde106886a53080
                name: ap-south-1
              - ami: ami-06602da18c878f98d
                name: eu-west-3
              - ami: ami-09a56048b08f94cdf
                name: eu-west-2
              - ami: ami-0943382e114f188e8
                name: eu-west-1
              - ami: ami-092faff259afb9a26
                name: ap-northeast-3
              - ami: ami-0ba5cd124d7a79612
                name: ap-northeast-2
              - ami: ami-0fe22bffdec36361c
                name: ap-northeast-1
              - ami: ami-05aa753c043f1dcd3
                name: sa-east-1
              - ami: ami-0e28822503eeedddc
                name: ca-central-1
              - ami: ami-055147723b7bca09a
                name: ap-southeast-1
              - ami: ami-0f39d06d145e9bb63
                name: ap-southeast-2
              - ami: ami-0b1deee75235aa4bb
                name: eu-central-1
              - ami: ami-0747bdcabd34c712a
                name: us-east-1
              - ami: ami-0b9064170e32bde34
                name: us-east-2
              - ami: ami-07b068f843ec78e72
                name: us-west-1
              - ami: ami-090717c950a5c34d3
                name: us-west-2
              - ami: ami-0448311ded7d81e94
                name: us-gov-east-1
              - ami: ami-0b0e99fc26b846798
                name: us-gov-west-1
            version: 18.4.20210415
          - regions:
              - ami: ami-01f90b0460589991e
                name: ap-northeast-1
              - ami: ami-0cb1ddea3786f6c0d
                name: sa-east-1
              - ami: ami-07042e91d04b1c30d
                name: eu-west-1
              - ami: ami-0d11056c10bfdde69
                name: ap-south-1
              - ami: ami-0367b500fdcac0edc
                name: us-east-2
              - ami: ami-0edf3b95e26a682df
                name: us-west-2
              - ami: ami-046842448f9e74e7d
                name: us-east-1
              - ami: ami-064efdb82ae15e93f
                name: ca-central-1
              - ami: ami-07ce5f60a39f1790e
                name: ap-southeast-1
              - ami: ami-04c7af7de7ad468f0
                name: ap-southeast-2
              - ami: ami-0e850e0e9c20d9deb
                name: eu-north-1
              - ami: ami-0c367ebddcf279dc6
                name: eu-west-3
              - ami: ami-096e3ded41e3bda6a
                name: ap-northeast-2
              - ami: ami-0718a1ae90971ce4d
                name: eu-central-1
              - ami: ami-04cc79dd5df3bffca
                name: eu-west-2
              - ami: ami-0d58800f291760030
                name: us-west-1
              - ami: ami-cc7598bd
                name: us-gov-east-1
              - ami: ami-f2664c93
                name: us-gov-west-1
            version: 18.4.20200228

    machineImagesProviderLs:
      - name: flatcar
        versions:
          - version: 2765.2.6
            regions:
              - name: ap-east-1
                ami: ami-0df249fad4e423ae7
              - name: ap-northeast-1
                ami: ami-028697e8df1ff071b
              - name: ap-northeast-2
                ami: ami-0c5b8f2d07d21da16
              - name: ap-south-1
                ami: ami-055b64c22dbcd61b0
              - name: ap-southeast-1
                ami: ami-045357ea038a43fe7
              - name: ap-southeast-2
                ami: ami-05df81d055054698d
              - name: ca-central-1
                ami: ami-0f639872bfcb49738
              - name: eu-central-1
                ami: ami-055acc5a6e9587b44
              - name: eu-north-1
                ami: ami-04f64f11f4dacda92
              - name: eu-west-1
                ami: ami-019f09de46e4e3f88
              - name: eu-west-2
                ami: ami-0097d8b6241e9cf76
              - name: eu-west-3
                ami: ami-05b8b0131fbb39283
              - name: me-south-1
                ami: ami-0737c661a0881fd94
              - name: sa-east-1
                ami: ami-00bc3ae33287bc81b
              - name: us-east-1
                ami: ami-0fd66875fa1ef8395
              - name: us-east-2
                ami: ami-02eb704ee029f6b9e
              - name: us-west-1
                ami: ami-053fb35697f85574d
              - name: us-west-2
                ami: ami-019657181ea76e880

    includeFilters: []
    excludeFilters: []
    disableMachineImages: []

  exports:
    data:
      - name: machineImages
        dataRef: machine-images-exports
EOF

echo "Installation stored at ${INSTALLATION_PATH}"
