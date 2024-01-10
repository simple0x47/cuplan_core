package config

import (
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime"
	"testing"
	"time"
)

type FileProviderTestSuite struct {
	suite.Suite
	provider          *FileProvider
	configurationFile string
}

func TestFileProviderTestSuite(t *testing.T) {
	suite.Run(t, new(FileProviderTestSuite))
}

func (f *FileProviderTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	testDataPath := core.GetTestDataPath(testFile).Unwrap()
	f.provider = NewFileProvider(testDataPath, core.NewCache(time.Hour), time.Hour)
	f.configurationFile = "application.yaml"
}

func (f *FileProviderTestSuite) TestFileProvider_Get_ReturnsExpectedValue() {
	const key = "Example:Inner:Value"
	const value = 5

	result := f.provider.Get(f.configurationFile, key)

	assert.True(f.T(), result.IsOk())
	assert.Equal(f.T(), value, result.Unwrap())
}

func (f *FileProviderTestSuite) TestFileProvider_Get_RootKey_ReturnsExpectedValue() {
	const key = "Root"
	const value = "yes"

	result := f.provider.Get(f.configurationFile, key)

	assert.True(f.T(), result.IsOk())
	assert.Equal(f.T(), value, result.Unwrap())
}

func (f *FileProviderTestSuite) TestFileProvider_Get_EvenLevelKey_ReturnsExpectedValue() {
	const key = "Example:Yeah"
	const value = true

	result := f.provider.Get(f.configurationFile, key)

	assert.True(f.T(), result.IsOk())
	assert.Equal(f.T(), value, result.Unwrap())
}

func (f *FileProviderTestSuite) TestFileProvider_Get_CachedKey_ReturnsExpectedValue() {
	f.TestFileProvider_Get_ReturnsExpectedValue()
	f.TestFileProvider_Get_ReturnsExpectedValue()
}
