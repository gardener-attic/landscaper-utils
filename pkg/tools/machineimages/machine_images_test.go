package machineimages

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

var _ = Describe("machine images", func() {

	Context("validateFilters", func() {

		It("should accept disjoint filters", func() {
			err := validateFilters(
				[]OsImagesFilterKind{OsImagesFilterKindCoreos, OsImagesFilterKindFlatcar},
				[]OsImagesFilterKind{OsImagesFilterKindPreview})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject overlapping filters", func() {
			err := validateFilters(
				[]OsImagesFilterKind{OsImagesFilterKindCoreos, OsImagesFilterKindFlatcar},
				[]OsImagesFilterKind{OsImagesFilterKindPreview, OsImagesFilterKindCoreos})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("convertOsImagesToMachineImages", func() {

		newVersion := func(version string) MachineImageVersion {
			return map[string]interface{}{"version": version}
		}

		It("should convert os images", func() {
			osImages := []OsImage{
				{Name: OsNameUbuntu, Version: newVersion("0.10.0")},
				{Name: OsNameUbuntu, Version: newVersion("0.11.0")},
				{Name: OsNameCoreos, Version: newVersion("0.20.0")},
				{Name: OsNameCoreos, Version: newVersion("0.21.0")},
			}
			machineImages := convertOsImagesToMachineImages(osImages)
			Expect(machineImages).To(HaveLen(2))
			for _, machineImage := range machineImages {
				if machineImage.Name == OsNameUbuntu {
					Expect(machineImage.Versions).To(ConsistOf(
						newVersion("0.10.0"),
						newVersion("0.11.0"),
					))
				} else if machineImage.Name == OsNameCoreos {
					Expect(machineImage.Versions).To(ConsistOf(
						newVersion("0.20.0"),
						newVersion("0.21.0"),
					))
				} else {
					Fail("unexpected os image name")
				}
			}
		})
	})

	Context("computeMachineImages", func() {

		readMachineImages := func(path string) ([]MachineImage, error) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}

			machineImages := []MachineImage{}
			if err := yaml.Unmarshal(data, &machineImages); err != nil {
				return nil, err
			}

			return machineImages, nil
		}

		It("should compute the machine images", func() {
			lssOsImages, err := readMachineImages("./resources/images.yaml")
			Expect(err).NotTo(HaveOccurred())

			landscapeOsImages, err := readMachineImages("./resources/images-ls.yaml")
			Expect(err).NotTo(HaveOccurred())

			providerOsImages, err := readMachineImages("./resources/images-pr.yaml")
			Expect(err).NotTo(HaveOccurred())

			providerLandscapeOsImages, err := readMachineImages("./resources/images-ls-pr.yaml")
			Expect(err).NotTo(HaveOccurred())

			disabledMachineImages := []string{}
			includeFileters := []OsImagesFilterKind{}
			excludeFilters := []OsImagesFilterKind{}

			machineImages, err := ComputeMachineImages(
				lssOsImages,
				landscapeOsImages,
				providerOsImages,
				providerLandscapeOsImages,
				disabledMachineImages,
				includeFileters,
				excludeFilters,
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(machineImages).NotTo(BeNil())

			expectedImages, err := readMachineImages("./resources/images-out.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(machineImages).To(Equal(expectedImages))
		})
	})
})
