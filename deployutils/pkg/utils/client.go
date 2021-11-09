// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// NewClientFromTarget creates a new kubernetes client for the kubeconfig in the given target.
func NewClientFromTarget(target lsv1alpha1.Target, clientOptions client.Options) (client.Client, error) {
	targetConfig := target.Spec.Configuration.RawMessage
	targetConfigMap := make(map[string]string)

	err := yaml.Unmarshal(targetConfig, &targetConfigMap)
	if err != nil {
		return nil, err
	}

	kubeconfig, ok := targetConfigMap["kubeconfig"]
	if !ok {
		return nil, fmt.Errorf("Imported target does not contain a kubeconfig")
	}

	return NewClientFromKubeconfig([]byte(kubeconfig), clientOptions)
}

// NewClientFromKubeconfig creates a new kubernetes client for the given kubeconfig.
func NewClientFromKubeconfig(kubeconfig []byte, clientOptions client.Options) (client.Client, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return client.New(restConfig, clientOptions)
}
