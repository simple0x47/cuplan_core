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

func (s *ZipExtractorTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	s.TestDataPath = core.GetTestDataPath(testFile).Unwrap()
}

func (s *ZipExtractorTestSuite) TestZipExtractor_Extract_ValidZip_ExtractsExpectedFiles() {
	extractor := NewZipExtractor()
	uuid := uuid.New().String()
	packageData, err := os.ReadFile(fmt.Sprintf("%sdummy.zip", s.TestDataPath))
	if err != nil {
		assert.Fail(s.T(), fmt.Sprintf("failed to read 'dummy.zip': %s", err))
	}

	result := extractor.Extract(packageData, uuid)

	testDataPathResult := doesDirectoryExist(s.TestDataPath)
	_, executableErr := os.Stat(fmt.Sprintf("%s/cp-config", uuid))
	_, configFileErr := os.Stat(fmt.Sprintf("%s/config/config.yaml", uuid))
	_, logConfigFileErr := os.Stat(fmt.Sprintf("%s/config/log4rs.yaml", uuid))
	_, anotherFileErr := os.Stat(fmt.Sprintf("%s/config/subfolder/another.yaml", uuid))
	os.RemoveAll(uuid)
	assert.True(s.T(), result.IsOk())
	assert.True(s.T(), testDataPathResult)
	assert.Equal(s.T(), nil, executableErr, "Expected executable file does not exist.")
	assert.Equal(s.T(), nil, configFileErr, "Expected configuration file does not exist.")
	assert.Equal(s.T(), nil, logConfigFileErr, "Expected log configuration file does not exist.")
	assert.Equal(s.T(), nil, anotherFileErr, "Expected another file does not exist.")
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
