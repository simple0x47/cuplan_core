package config

import "github.com/simpleg-eu/cuplan_core/pkg/core"

// Downloader
// Interface which provides a facility to download configuration packages.
type Downloader interface {
	// Download
	// Downloads the latest configuration from a specific configuration provider.
	//
	// * host - The host from which the configuration package is going to be downloaded.
	//
	// * stage - Flavour of the configuration package to be downloaded.
	//
	// * environment - Environment, i.e. 'development', 'staging', 'production'.
	//
	//* component - Microservice for which the configuration package is downloaded.
	//
	// Returns: Bytes slice containing the configuration package or an error.
	Download(host string, stage string, environment string, component string) core.Result[[]byte, core.Error]
}
