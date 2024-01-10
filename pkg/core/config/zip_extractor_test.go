package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"os"
	"runtime"
	"testing"
)

type ZipExtractorTestSuite struct {
	suite.Suite
	TestDataPath string
	Extractor    *ZipExtractor
}

func (z *ZipExtractorTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	z.TestDataPath = core.GetTestDataPath(testFile).Unwrap()
	logger, _ := zap.NewDevelopment()
	z.Extractor = NewZipExtractor(logger)
}

func (z *ZipExtractorTestSuite) TestZipExtractor_Extract_ValidZip_ExtractsExpectedFiles() {
	targetPath := uuid.New().String()
	packageData, err := os.ReadFile(fmt.Sprintf("%sdummy.zip", z.TestDataPath))
	if err != nil {
		assert.Fail(z.T(), fmt.Sprintf("failed to read 'dummy.zip': %s", err))
	}

	result := z.Extractor.Extract(packageData, targetPath)

	testDataPathResult := doesDirectoryExist(z.TestDataPath)
	_, executableErr := os.Stat(fmt.Sprintf("%s/cp-config", targetPath))
	_, configFileErr := os.Stat(fmt.Sprintf("%s/config/config.yaml", targetPath))
	_, logConfigFileErr := os.Stat(fmt.Sprintf("%s/config/log4rs.yaml", targetPath))
	_, anotherFileErr := os.Stat(fmt.Sprintf("%s/config/subfolder/another.yaml", targetPath))
	_ = os.RemoveAll(targetPath)
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
