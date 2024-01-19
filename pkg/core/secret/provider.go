package secret

import "github.com/simpleg-eu/cuplan_core/pkg/core"

type Provider interface {
	Get(string) core.Result[string, core.Error]
}
