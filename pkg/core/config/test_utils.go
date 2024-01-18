package config

import (
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"time"
)

func GetConfigForTest(testFile string) core.Result[*FileProvider, core.Error] {
	f := NewFileProvider(core.GetTestDataPath(testFile).Unwrap(), core.NewCache(time.Hour), time.Hour)

	return core.Ok[*FileProvider, core.Error](f)
}
