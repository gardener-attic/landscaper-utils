// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package machineimages

import (
	"fmt"
)

type OsImagesFilterKind string

const (
	OsImagesFilterKindAll            = OsImagesFilterKind("all")
	OsImagesFilterKindOutdated       = OsImagesFilterKind("outdated")
	OsImagesFilterKindPreview        = OsImagesFilterKind("preview")
	OsImagesFilterKindSupported      = OsImagesFilterKind("supported")
	OsImagesFilterKindDeprecated     = OsImagesFilterKind("deprecated")
	OsImagesFilterKindGardenlinux    = OsImagesFilterKind("gardenlinux")
	OsImagesFilterKindSuseChost      = OsImagesFilterKind("suse-chost")
	OsImagesFilterKindUbuntu         = OsImagesFilterKind("ubuntu")
	OsImagesFilterKindCoreos         = OsImagesFilterKind("coreos")
	OsImagesFilterKindFlatcar        = OsImagesFilterKind("flatcar")
	OsImagesFilterKindMemoryoneChost = OsImagesFilterKind("memoryone-chost")
)

const (
	ClassificationDeprecated = "deprecated"
	ClassificationPreview    = "preview"
	ClassificationSupported  = "supported"
)

const (
	OsNameGardenLinux    = "gardenlinux"
	OsNameSuseChost      = "suse-chost"
	OsNameUbuntu         = "ubuntu"
	OsNameCoreos         = "coreos"
	OsNameFlatcar        = "flatcar"
	OsNameMemoryoneChost = "memoryone-chost"
)

type OsImageFilter interface {
	match(image OsImage) (bool, error)
}

func filter(images []OsImage, f OsImageFilter) ([]OsImage, error) {
	result := []OsImage{}

	for _, image := range images {
		matched, err := f.match(image)
		if err != nil {
			return nil, err
		}

		if matched {
			result = append(result, image)
		}
	}

	return result, nil
}

type negatedFilter struct {
	filter OsImageFilter
}

// match returns true iff the original filter does not match.
func (n negatedFilter) match(image OsImage) (bool, error) {
	matched, err := n.filter.match(image)
	if err != nil {
		return false, err
	}

	return !matched, nil
}

type anyFilter struct {
	filters []OsImageFilter
}

// match returns true iff at least one of the original filters matches.
func (a anyFilter) match(image OsImage) (bool, error) {
	for _, f := range a.filters {
		matched, err := f.match(image)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}
	return false, nil
}

func filterOsImages(
	images []OsImage,
	includeFilterKinds []OsImagesFilterKind,
	excludeFilterKinds []OsImagesFilterKind,
) ([]OsImage, error) {
	includeFilters, err := createFilters(includeFilterKinds)
	if err != nil {
		return nil, err
	}

	excludeFilters, err := createFilters(excludeFilterKinds)
	if err != nil {
		return nil, err
	}

	result, err := filter(images, anyFilter{includeFilters})
	if err != nil {
		return nil, err
	}

	result, err = filter(result, negatedFilter{anyFilter{excludeFilters}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createFilters(filterKinds []OsImagesFilterKind) ([]OsImageFilter, error) {
	result := make([]OsImageFilter, len(filterKinds))

	for i, kind := range filterKinds {
		filter, err := createFilter(kind)
		if err != nil {
			return nil, err
		}

		result[i] = filter
	}

	return result, nil
}

func createFilter(filterKind OsImagesFilterKind) (OsImageFilter, error) {
	switch filterKind {
	case OsImagesFilterKindAll:
		return &allowAllFilter{}, nil
	case OsImagesFilterKindOutdated:
		return &outdatedFilter{}, nil
	case OsImagesFilterKindDeprecated:
		return &classificationImagesFilter{classification: ClassificationDeprecated}, nil
	case OsImagesFilterKindPreview:
		return &classificationImagesFilter{classification: ClassificationPreview}, nil
	case OsImagesFilterKindSupported:
		return &classificationImagesFilter{classification: ClassificationSupported}, nil
	case OsImagesFilterKindGardenlinux:
		return &osNameImagesFilter{osName: OsNameGardenLinux}, nil
	case OsImagesFilterKindSuseChost:
		return &osNameImagesFilter{osName: OsNameSuseChost}, nil
	case OsImagesFilterKindUbuntu:
		return &osNameImagesFilter{osName: OsNameUbuntu}, nil
	case OsImagesFilterKindCoreos:
		return &osNameImagesFilter{osName: OsNameCoreos}, nil
	case OsImagesFilterKindFlatcar:
		return &osNameImagesFilter{osName: OsNameFlatcar}, nil
	case OsImagesFilterKindMemoryoneChost:
		return &osNameImagesFilter{osName: OsNameMemoryoneChost}, nil
	default:
		return nil, fmt.Errorf("filter does not exist %s", filterKind)
	}
}

type allowAllFilter struct{}

func (a *allowAllFilter) match(_ OsImage) (bool, error) {
	return true, nil
}

type outdatedFilter struct{}

func (a *outdatedFilter) match(image OsImage) (bool, error) {
	expired, err := image.Version.isExpired()
	if err != nil {
		return false, err
	}
	return image.Version.hasClassification(ClassificationDeprecated) && expired, nil
}

type classificationImagesFilter struct {
	classification string
}

func (a *classificationImagesFilter) match(image OsImage) (bool, error) {
	return image.Version.hasClassification(a.classification), nil
}

type osNameImagesFilter struct {
	osName string
}

func (a *osNameImagesFilter) match(image OsImage) (bool, error) {
	return image.Name == a.osName, nil
}
