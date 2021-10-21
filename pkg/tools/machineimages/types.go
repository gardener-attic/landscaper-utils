package machineimages

import "time"

type Imports struct {
	MachineImages           []MachineImage       `json:"machineImages" yaml:"machineImages"`
	MachineImagesLs         []MachineImage       `json:"machineImagesLs" yaml:"machineImagesLs"`
	MachineImagesProvider   []MachineImage       `json:"machineImagesProvider" yaml:"machineImagesProvider"`
	MachineImagesProviderLs []MachineImage       `json:"machineImagesProviderLs" yaml:"machineImagesProviderLs"`
	IncludeFilters          []OsImagesFilterKind `json:"includeFilters" yaml:"includeFilters"`
	ExcludeFilters          []OsImagesFilterKind `json:"excludeFilters" yaml:"excludeFilters"`
	DisableMachineImages    []string             `json:"disableMachineImages" yaml:"disableMachineImages"`
}

type Exports struct {
	ResultMachineImages []MachineImage `json:"resultMachineImages" yaml:"resultMachineImages"`
}

type MachineImage struct {
	Name     string                `json:"name,omitempty"`
	Versions []MachineImageVersion `json:"versions,omitempty"`
}

type MachineImageVersion map[string]interface{}

func (v MachineImageVersion) getClassification() *string {
	m := map[string]interface{}(v)
	value, ok := m["classification"].(string)
	if !ok {
		return nil
	}

	return &value
}

func (v MachineImageVersion) getVersion() *string {
	m := map[string]interface{}(v)
	value, ok := m["version"].(string)
	if !ok {
		return nil
	}

	return &value
}

func (v MachineImageVersion) hasClassification(classification string) bool {
	c := v.getClassification()
	return c != nil && *c == classification
}

func (v MachineImageVersion) getExpirationDate() (*time.Time, error) {
	m := map[string]interface{}(v)
	value, ok := m["expirationDate"].(string)
	if !ok {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (v MachineImageVersion) isExpired() (bool, error) {
	t, err := v.getExpirationDate()
	if err != nil {
		return false, err
	}

	return t != nil && time.Now().After(*t), nil
}

type OsImage struct {
	Name    string              `json:"name,omitempty"`
	Version MachineImageVersion `json:"version,omitempty"`
}
