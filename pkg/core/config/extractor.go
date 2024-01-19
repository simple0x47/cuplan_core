package config

import "github.com/simpleg-eu/cuplan_core/pkg/core"

// Extractor
// Interface which provides a facility to extract a configuration package.
type Extractor interface {
	// Extract
	// Extracts the configuration package's content into the targetPath.
	//
	// * packageData - Package's raw data.
	//
	// * targetPath - Path where the configuration will be extracted into.
	Extract(packageData []byte, targetPath string) core.Result[core.Empty, core.Error]
}
