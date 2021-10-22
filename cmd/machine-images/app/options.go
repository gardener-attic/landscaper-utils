package app

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gardener/landscaper-utils/pkg/tools/machineimages"
	"sigs.k8s.io/yaml"

	"github.com/spf13/pflag"
)

const (
	EnvVarImportsPath = "IMPORTS_PATH"
	EnvVarExportsPath = "EXPORTS_PATH"
)

type options struct {
	// ImportsPath is the path to the imports file.
	ImportsPath string
	// ExportsPath is the path to the exports file.
	ExportsPath string
}

func newOptions() *options {
	return &options{}
}

func (o *options) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ImportsPath, "imports-path", "i", "", "The path to the imports file")
	fs.StringVarP(&o.ExportsPath, "exports-path", "e", "", "The path to the exports file")
}

// complete parses all options and flags and initializes the basic functions
func (o *options) complete() error {
	if len(o.ImportsPath) == 0 {
		o.ImportsPath = os.Getenv(EnvVarImportsPath)
	}

	if len(o.ExportsPath) == 0 {
		o.ExportsPath = os.Getenv(EnvVarExportsPath)
	}

	return o.validate()
}

func (o *options) validate() error {
	if len(o.ImportsPath) == 0 {
		return errors.New("an imports path must be provided. ")
	}

	if len(o.ExportsPath) == 0 {
		return errors.New("an exports path must be provided. ")
	}

	return nil
}

func (o *options) run(ctx context.Context) error {
	imports, err := o.readImports()
	if err != nil {
		return err
	}

	result, err := machineimages.ComputeMachineImages(
		imports.MachineImages,
		imports.MachineImagesLs,
		imports.MachineImagesProvider,
		imports.MachineImagesProviderLs,
		imports.DisableMachineImages,
		imports.IncludeFilters,
		imports.ExcludeFilters,
	)
	if err != nil {
		return err
	}

	err = o.writeExports(&machineimages.Exports{ResultMachineImages: result})
	return err
}

func (o *options) readImports() (*machineimages.Imports, error) {
	data, err := ioutil.ReadFile(o.ImportsPath)
	if err != nil {
		return nil, err
	}

	imports := &machineimages.Imports{}
	if err := yaml.Unmarshal(data, imports); err != nil {
		return nil, err
	}

	return imports, nil
}

func (o *options) writeExports(exports *machineimages.Exports) error {
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
