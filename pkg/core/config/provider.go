package config

import "github.com/simpleg-eu/cuplan_core/pkg/core"

// Provider provides a configuration value. I know, crazy, huh?
type Provider interface {
	// Get provides the configuration value for the specified key.
	Get(filePath string, key string) core.Result[any, core.Error]

	// CleanCache cleans any possible caching mechanism.
	CleanCache()
}
