package config

import (
	"fmt"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/simpleg-eu/cuplan-core/pkg/core/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
	"testing"
	"time"
)

type Config struct {
	AccessTokenSecret        string `yaml:"AccessTokenSecret"`
	Host                     string `yaml:"Host"`
	Stage                    string `yaml:"Stage"`
	Environment              string `yaml:"Environment"`
	Component                string `yaml:"Component"`
	DownloadTimeoutInSeconds int    `yaml:"DownloadTimeoutInSeconds"`
}

type ServerDownloaderTestSuite struct {
	suite.Suite
	TestDataPath string
	Config       Config
	Downloader   Downloader
}

func TestServerDownloaderTestSuite(t *testing.T) {
	suite.Run(t, new(ServerDownloaderTestSuite))
}

func (s *ServerDownloaderTestSuite) SetupTest() {
	_, testFile, _, _ := runtime.Caller(0)
	s.TestDataPath = core.GetTestDataPath(testFile).Unwrap()
	configFilePath := fmt.Sprintf("%sconfig.yaml", s.TestDataPath)

	configFile, err := os.ReadFile(configFilePath)

	if err != nil {
		panic(fmt.Sprintf("Error reading configuration YAML file: %v\n", err))
	}

	err = yaml.Unmarshal(configFile, &s.Config)

	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling the configuration YAML file: %v\n", err))
	}

	secretsProvider := secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken())
	logger, _ := zap.NewDevelopment()
	s.Downloader = NewServerDownloader(logger, secretsProvider.Get(s.Config.AccessTokenSecret).Unwrap(), time.Second*time.Duration(s.Config.DownloadTimeoutInSeconds))
}

func (s *ServerDownloaderTestSuite) TestServerDownloader_Download_ReturnsBytes() {
	result := s.Downloader.Download(s.Config.Host, s.Config.Stage, s.Config.Environment, s.Config.Component)

	assert.True(s.T(), result.IsOk())
	assert.True(s.T(), len(result.Unwrap()) > 0, "Unexpected empty configuration package.")
}

func (s *ServerDownloaderTestSuite) TestServerDownloader_Download_Fail_NonExistingHost() {
	result := s.Downloader.Download("https://hereconfigconfiglmao.com", s.Config.Stage, s.Config.Environment, s.Config.Component)

	assert.False(s.T(), result.IsOk())
	assert.Equal(s.T(), core.ConfigurationRetrievalFailure, result.UnwrapErr().ErrorKind)
}
