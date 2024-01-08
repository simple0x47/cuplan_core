package secret

import core "github.com/simpleg-eu/cuplan-core"

type Provider interface {
	Get(string) core.Result[string, core.Error]
}
