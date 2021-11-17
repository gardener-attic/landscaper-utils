// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gardener/component-cli/ociclient/oci"
	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/codec"
	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"sigs.k8s.io/yaml"

	"github.com/gardener/landscaper-utils/deployutils/pkg/logger"
)

const (
	EnvVarOperation               = "OPERATION"
	EnvVarImportsPath             = "IMPORTS_PATH"
	EnvVarExportsPath             = "EXPORTS_PATH"
	EnvVarComponentDescriptorPath = "COMPONENT_DESCRIPTOR_PATH"
)

const (
	// OperationReconcile is a  constant for the RECONCILE operation.
	OperationReconcile = "RECONCILE"
	// OperationDelete is a constant for the DELETE operation.
	OperationDelete = "DELETE"
)

type Options struct {
	Log logr.Logger
	// Operation is the operation to be executed.
	Operation string
	// ImportsPath is the path to the imports file.
	ImportsPath string
	// ExportsPath is the path to the exports file.
	ExportsPath string
	// ComponentDescriptorPath is the path to the component descriptor file.
	ComponentDescriptorPath string
}

func NewOptions() *Options {
	log, err := logger.NewCliLogger()
	if err != nil {
		fmt.Println("unable to setup logger")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	logger.SetLogger(log)

	return &Options{
		Log: log,
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Operation, "operation", "o", "", "The operation. Possible values: RECONCILE, DELETE")
	fs.StringVarP(&o.ImportsPath, "imports-path", "i", "", "The path to the imports file")
	fs.StringVarP(&o.ExportsPath, "exports-path", "e", "", "The path to the exports file")
	fs.StringVarP(&o.ComponentDescriptorPath, "component-desciptor-path", "c", "", "The path to the component descriptor file")
}

// complete parses all options and flags and initializes the basic functions
func (o *Options) Complete() error {
	if len(o.Operation) == 0 {
		if op := os.Getenv(EnvVarOperation); len(op) > 0 {
			o.Operation = op
		}
	}

	if len(o.ImportsPath) == 0 {
		o.ImportsPath = os.Getenv(EnvVarImportsPath)
	}

	if len(o.ExportsPath) == 0 {
		o.ExportsPath = os.Getenv(EnvVarExportsPath)
	}

	if len(o.ComponentDescriptorPath) == 0 {
		o.ComponentDescriptorPath = os.Getenv(EnvVarComponentDescriptorPath)
	}

	return o.validate()
}

func (o *Options) validate() error {
	if o.Operation != OperationReconcile && o.Operation != OperationDelete {
		return fmt.Errorf("the operation must be %q or %q", OperationReconcile, OperationDelete)
	}

	if len(o.ImportsPath) == 0 {
		return errors.New("an imports path must be provided. ")
	}

	if len(o.ExportsPath) == 0 {
		return errors.New("an exports path must be provided. ")
	}

	if len(o.ComponentDescriptorPath) == 0 {
		return errors.New("a component descriptor path must be provided. ")
	}

	return nil
}

func (o *Options) ReadImports(imports interface{}) error {
	o.Log.Info("Reading imports", "imports-path", o.ImportsPath)

	data, err := ioutil.ReadFile(o.ImportsPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, imports); err != nil {
		return err
	}

	return nil
}

func (o *Options) WriteExports(exports interface{}) error {
	o.Log.Info("Writing exports", "exports-path", o.ExportsPath)

	b, err := yaml.Marshal(exports)
	if err != nil {
		return err
	}

	parentPath := filepath.Dir(o.ExportsPath)
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		if err := os.MkdirAll(parentPath, 0700); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(o.ExportsPath, b, os.ModePerm)
}

// GetComponentDescriptors returns the ComponentDescriptorList.
func (o *Options) GetComponentDescriptors() (*cdv2.ComponentDescriptorList, error) {
	o.Log.Info("Reading component descriptor", "component-descriptor-path", o.ComponentDescriptorPath)

	data, err := ioutil.ReadFile(o.ComponentDescriptorPath)
	if err != nil {
		return nil, fmt.Errorf("error reading component descriptor file: %w ", err)
	}

	cdList := &cdv2.ComponentDescriptorList{}
	if err := codec.Decode(data, cdList); err != nil {
		return nil, fmt.Errorf("error decoding component descriptor list: %w ", err)
	}

	return cdList, nil
}

// GetComponentDescriptorByName returns the component descriptor with the given name.
// Returns an error if there is no such component descriptor, or more than one.
func (o *Options) GetComponentDescriptorByName(componentName string) (*cdv2.ComponentDescriptor, error) {
	cdList, err := o.GetComponentDescriptors()
	if err != nil {
		return nil, err
	}

	cds := []cdv2.ComponentDescriptor{}
	for _, cd := range cdList.Components {
		if cd.Name == componentName {
			cds = append(cds, cd)
		}
	}

	if len(cds) == 0 {
		return nil, fmt.Errorf("no component descriptor with name %s found", componentName)
	}

	if len(cds) > 1 {
		return nil, fmt.Errorf("more than one component descriptor with name %s found", componentName)
	}

	return &cds[0], nil
}

// GetResourceByName returns the entry with the given name from the "resources" section of the component descriptor.
// Returns an error if there is no such entry, or more than one.
func (o *Options) GetResourceByName(cd *cdv2.ComponentDescriptor, resourceName string) (*cdv2.Resource, error) {
	nameSelector := cdv2.NewNameSelector(resourceName)
	resources, err := cd.GetResourcesBySelector(nameSelector)
	if err != nil {
		return nil, fmt.Errorf("error selecting resources with name %s: %w", resourceName, err)
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no resource with name %s found", resourceName)
	}

	if len(resources) > 1 {
		return nil, fmt.Errorf("more than one resource with name %s found", resourceName)
	}

	return &resources[0], nil
}

// GetResourceOCIImageReference returns the OCI image reference of the resource with the given name
// from the given component descriptor.
func (o *Options) GetResourceOCIImageReference(cd *cdv2.ComponentDescriptor, resourceName string) (string, error) {
	resource, err := o.GetResourceByName(cd, resourceName)
	if err != nil {
		return "", err
	}

	access := &cdv2.OCIRegistryAccess{}
	if err := resource.Access.DecodeInto(access); err != nil {
		return "", fmt.Errorf("error decoding acceess data of resource %s: %w ", resourceName, err)
	}

	return access.ImageReference, nil
}

// GetResourceOCIRepositoryAndTag returns the repository and tag of the resource with the given name
// from the given component descriptor.
func (o *Options) GetResourceOCIRepositoryAndTag(cd *cdv2.ComponentDescriptor, resourceName string) (repository, tag string, err error) {
	imageReference, err := o.GetResourceOCIImageReference(cd, resourceName)
	if err != nil {
		return "", "", err
	}

	refSpec, err := oci.ParseRef(imageReference)
	if err != nil {
		return "", "", fmt.Errorf("error parsing oci image reference of resource %s: %w ", resourceName, err)
	}

	repository = refSpec.Name()

	if refSpec.Tag != nil {
		tag = *refSpec.Tag
	} else if refSpec.Digest != nil {
		tag = refSpec.Digest.String()
	} else {
		return "", "", fmt.Errorf("image reference of resource %s has neither tag, nor digest. ", resourceName)
	}

	return repository, tag, nil
}
