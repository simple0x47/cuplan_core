package config

import (
	"fmt"
	"github.com/simpleg-eu/cuplan_core/pkg/core"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
	"time"
)

const keySeparator = ":"

type FileProvider struct {
	targetPath           string
	cache                *core.Cache
	expireCacheItemAfter time.Duration
}

func NewFileProvider(targetPath string, cache *core.Cache, expireCacheItemAfter time.Duration) *FileProvider {
	provider := new(FileProvider)

	provider.targetPath = targetPath
	provider.cache = cache
	provider.expireCacheItemAfter = expireCacheItemAfter

	return provider
}

func (f *FileProvider) Get(filePath string, key string) core.Result[any, core.Error] {
	filePath = fmt.Sprintf("%s/%s", f.targetPath, filePath)

	cache := f.cache.Get(filePath)

	var config map[string]any

	if cache.IsSome() {
		var ok bool
		config, ok = cache.Unwrap().(map[string]any)

		if !ok {
			return core.Err[any, core.Error](*core.NewError(core.InvalidCache, "failed to read cache as 'map[string]any'"))
		}
	} else {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return core.Err[any, core.Error](*core.NewError(core.NotFound, fmt.Sprintf("couldn't find file: %s", filePath)))
		}

		yamlConfig, err := os.ReadFile(filePath)

		if err != nil {
			return core.Err[any, core.Error](*core.NewError(core.IOFailure, fmt.Sprintf("failed to open file '%s': %s", filePath, err)))
		}

		if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
			return core.Err[any, core.Error](*core.NewError(core.SerializationFailure, fmt.Sprintf("failed to read file's content '%s' as YAML: %s", filePath, err)))
		}

		f.cache.Set(filePath, config, f.expireCacheItemAfter)
	}

	return getValueFromKeys[any](key, config)
}

func getValueFromKeys[T any](key string, object map[string]any) core.Result[T, core.Error] {
	subKeys := strings.Split(key, keySeparator)

	value := object[subKeys[0]]
	var ok = true

	subKeys = subKeys[1:]
	for i, subKey := range subKeys {
		if len(subKeys)-1 != i {
			if value, ok = value.(map[string]any)[subKey].(map[string]any); !ok {
				return core.Err[T, core.Error](*core.NewError(core.InvalidInput, fmt.Sprintf("failed to read key '%s'", subKey)))
			}
		} else {
			value = value.(map[string]any)[subKey]
		}
	}

	finalValue, ok := value.(T)

	if !ok {
		return core.Err[T, core.Error](*core.NewError(core.InvalidInput, fmt.Sprintf("failed to get key '%s' as T '%v'", subKeys[len(subKeys)-1], reflect.TypeOf(value).String())))
	}

	return core.Ok[T, core.Error](finalValue)
}
