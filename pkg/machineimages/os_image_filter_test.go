package machineimages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("os image filter", func() {

	Context("filterOsImages", func() {

		It("should include all images", func() {
			images := []OsImage{{Name: OsNameUbuntu}, {Name: OsNameCoreos}, {Name: OsNameGardenLinux}}
			filteredImages, err := filterOsImages(
				images,
				[]OsImagesFilterKind{OsImagesFilterKindAll},
				nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filteredImages).To(Equal(images))
		})

		It("should exclude all images", func() {
			images := []OsImage{{Name: OsNameUbuntu}, {Name: OsNameCoreos}, {Name: OsNameGardenLinux}}
			filteredImages, err := filterOsImages(
				images,
				nil,
				[]OsImagesFilterKind{OsImagesFilterKindAll})
			Expect(err).NotTo(HaveOccurred())
			Expect(filteredImages).To(Equal([]OsImage{}))
		})

		It("should include some images using the name filter", func() {
			images := []OsImage{{Name: OsNameUbuntu}, {Name: OsNameCoreos}, {Name: OsNameGardenLinux}}
			filteredImages, err := filterOsImages(
				images,
				[]OsImagesFilterKind{OsNameUbuntu, OsNameGardenLinux},
				nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filteredImages).To(ConsistOf(OsImage{Name: OsNameUbuntu}, OsImage{Name: OsNameGardenLinux}))
		})

		It("should exclude some images using the name filter", func() {
			images := []OsImage{{Name: OsNameUbuntu}, {Name: OsNameCoreos}, {Name: OsNameGardenLinux}}
			filteredImages, err := filterOsImages(
				images,
				[]OsImagesFilterKind{OsImagesFilterKindAll},
				[]OsImagesFilterKind{OsNameUbuntu, OsNameGardenLinux})
			Expect(err).NotTo(HaveOccurred())
			Expect(filteredImages).To(ConsistOf(OsImage{Name: OsNameCoreos}))
		})
	})
})
