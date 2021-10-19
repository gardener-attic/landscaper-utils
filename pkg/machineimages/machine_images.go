package machineimages

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

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

func computeMachineImages(
	lssOsImages []MachineImage,
	landscapeOsImages []MachineImage,
	providerOsImages []MachineImage,
	providerLandscapeOsImages []MachineImage,
	disableMachineImages []string,
	includeFilters []OsImagesFilterKind,
	excludeFilters []OsImagesFilterKind,
) (
	[]MachineImage,
	error,
) {
	if len(includeFilters) == 0 {
		includeFilters = append(includeFilters, OsImagesFilterKindAll)
	}

	err := validateFilters(includeFilters, excludeFilters)
	if err != nil {
		return nil, err
	}

	flatLandscapeOsImages := flatImages(landscapeOsImages)
	flatLssOsImages := flatImages(lssOsImages)
	flatOsImages := append(flatLandscapeOsImages, flatLssOsImages...)
	flatOsImages = removeDuplicates(flatOsImages)

	flatOsImages, err = filterOsImages(flatOsImages, includeFilters, excludeFilters)
	if err != nil {
		return nil, err
	}

	if len(flatOsImages) == 0 {
		return []MachineImage{}, nil
	}

	machineImages := convertOsImagesToMachineImages(flatOsImages)
	sort.SliceStable(machineImages, func(i, j int) bool {
		return machineImages[i].Name < machineImages[j].Name
	})
	sort.SliceStable(machineImages, func(i, j int) bool {
		return machineImages[i].Name == OsNameGardenLinux && machineImages[j].Name != OsNameGardenLinux
	})

	machineImages = getFilteredMachineImages(machineImages, disableMachineImages,
		providerLandscapeOsImages, providerOsImages)
	return machineImages, nil
}

func getFilteredMachineImages(
	machineImages []MachineImage,
	disableMachineImages []string,
	providerLandscapeOsImages []MachineImage,
	providerOsImages []MachineImage,
) []MachineImage {
	filteredImages := []MachineImage{}
	for _, nextImage := range machineImages {
		if contains(disableMachineImages, nextImage.Name) {
			continue
		}

		versionsWithConfig := []MachineImageVersion{}
		for _, nextVersion := range nextImage.Versions {
			versionNumber := nextVersion.getVersion()
			config := getVersionConfig(nextImage.Name, *versionNumber, providerLandscapeOsImages, providerOsImages)
			if config != nil {
				for nextKey, nextValue := range *config {
					nextVersion[nextKey] = nextValue
				}
				versionsWithConfig = append(versionsWithConfig, nextVersion)
			}
		}

		if len(versionsWithConfig) > 0 {
			filteredImages = append(filteredImages, MachineImage{
				Name:     nextImage.Name,
				Versions: versionsWithConfig,
			})
		}
	}

	return filteredImages
}

func getVersionConfig(imageName, versionNumber string, providerLandscapeOsImages, providerOsImages []MachineImage) *MachineImageVersion {
	config := getVersionConfigInternal(imageName, versionNumber, providerLandscapeOsImages)

	if config != nil {
		return config
	}

	config = getVersionConfigInternal(imageName, versionNumber, providerOsImages)
	return config
}

func getVersionConfigInternal(imageName, versionNumber string, images []MachineImage) *MachineImageVersion {
	for _, nextImage := range images {
		if nextImage.Name == imageName {
			for _, nextVersion := range nextImage.Versions {
				if nextVersion.getVersion() != nil && *nextVersion.getVersion() == versionNumber {
					return &nextVersion
				}
			}
		}
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func removeDuplicates(images []OsImage) []OsImage {
	result := []OsImage{}
	for _, nextImage := range images {
		found := false
		for _, nextResult := range result {
			if reflect.DeepEqual(nextImage, nextResult) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, nextImage)
		}
	}
	return result
}

func flatImages(images []MachineImage) []OsImage {
	result := []OsImage{}
	for _, nextImage := range images {
		for _, nextVersion := range nextImage.Versions {
			result = append(result, OsImage{
				Name:    nextImage.Name,
				Version: nextVersion,
			})
		}
	}
	return result
}

func convertOsImagesToMachineImages(images []OsImage) []MachineImage {
	m := map[string][]MachineImageVersion{}

	for _, image := range images {
		versions, ok := m[image.Name]
		if !ok {
			versions = []MachineImageVersion{image.Version}
		} else {
			versions = append(versions, image.Version)
		}
		m[image.Name] = versions
	}

	result := []MachineImage{}
	for name, versions := range m {
		result = append(result, MachineImage{
			Name:     name,
			Versions: versions,
		})
	}

	return result
}

func validateFilters(includeFilters []OsImagesFilterKind, excludeFilters []OsImagesFilterKind) error {
	for _, include := range includeFilters {
		for _, exclude := range excludeFilters {
			if include == exclude {
				return fmt.Errorf("exclude filter list contains element of include list")
			}
		}
	}

	return nil
}
