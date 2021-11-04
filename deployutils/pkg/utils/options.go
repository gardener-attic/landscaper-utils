// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"github.com/gardener/landscaper-utils/deployutils/pkg/logger"
	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
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
