package machineimages

type Imports struct {
	MachineImages           []MachineImage `json:"machineImages" yaml:"machineImages"`
	MachineImagesLs         []MachineImage `json:"machineImagesLs" yaml:"machineImagesLs"`
	MachineImagesProvider   []MachineImage `json:"machineImagesProvider" yaml:"machineImagesProvider"`
	MachineImagesProviderLs []MachineImage `json:"machineImagesProviderLs" yaml:"machineImagesProviderLs"`
}

type Exports struct {
	ResultMachineImages []MachineImage `json:"resultMachineImages" yaml:"resultMachineImages"`
}
