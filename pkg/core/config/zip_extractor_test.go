package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"runtime"
	"testing"
)

type ZipExtractorTestSuite struct {
	suite.Suite
	TestDataPath string
}

func (z *ZipExtractorTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	z.TestDataPath = core.GetTestDataPath(testFile).Unwrap()
}

func (z *ZipExtractorTestSuite) TestZipExtractor_Extract_ValidZip_ExtractsExpectedFiles() {
	extractor := NewZipExtractor()
	uuid := uuid.New().String()
	packageData, err := os.ReadFile(fmt.Sprintf("%sdummy.zip", z.TestDataPath))
	if err != nil {
		assert.Fail(z.T(), fmt.Sprintf("failed to read 'dummy.zip': %z", err))
	}

	result := extractor.Extract(packageData, uuid)

	testDataPathResult := doesDirectoryExist(z.TestDataPath)
	_, executableErr := os.Stat(fmt.Sprintf("%z/cp-config", uuid))
	_, configFileErr := os.Stat(fmt.Sprintf("%z/config/config.yaml", uuid))
	_, logConfigFileErr := os.Stat(fmt.Sprintf("%z/config/log4rs.yaml", uuid))
	_, anotherFileErr := os.Stat(fmt.Sprintf("%z/config/subfolder/another.yaml", uuid))
	os.RemoveAll(uuid)
	assert.True(z.T(), result.IsOk())
	assert.True(z.T(), testDataPathResult)
	assert.Equal(z.T(), nil, executableErr, "Expected executable file does not exist.")
	assert.Equal(z.T(), nil, configFileErr, "Expected configuration file does not exist.")
	assert.Equal(z.T(), nil, logConfigFileErr, "Expected log configuration file does not exist.")
	assert.Equal(z.T(), nil, anotherFileErr, "Expected another file does not exist.")
}

func TestZipExtractorTestSuite(t *testing.T) {
	suite.Run(t, new(ZipExtractorTestSuite))
}

func doesDirectoryExist(directory string) bool {
	_, err := os.Stat(directory)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
