package config

import "github.com/simpleg-eu/cuplan-core/pkg/core"

type Downloader interface {
	Download(host string, stage string, environment string, component string) core.Result[[]byte, core.Error]
}
