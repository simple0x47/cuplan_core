package config

import (
	"fmt"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/simpleg-eu/cuplan-core/pkg/core/secret"
	"github.com/stretchr/testify/assert"
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

var testDataPath string
var config Config
var downloader Downloader

func TestMain(m *testing.M) {
	setup()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestServerDownloader_Download_ReturnsBytes(t *testing.T) {
	result := downloader.Download(config.Host, config.Stage, config.Environment, config.Component)

	assert.True(t, result.IsOk())
	assert.True(t, len(result.Unwrap()) > 0, "Unexpected empty configuration package.")
}

func TestServerDownloader_Download_Fail_NonExistingHost(t *testing.T) {
	result := downloader.Download("https://hereconfigconfiglmao.com", config.Stage, config.Environment, config.Component)

	assert.False(t, result.IsOk())
	assert.Equal(t, core.ConfigurationRetrievalFailure, result.UnwrapErr().ErrorKind())
}

func setup() {
	_, testFile, _, _ := runtime.Caller(0)
	testDataPath = core.GetTestDataPath(testFile).Unwrap()
	configFilePath := testDataPath + "config.yaml"

	configFile, err := os.ReadFile(configFilePath)

	if err != nil {
		panic(fmt.Sprintf("Error reading configuration YAML file: %v\n", err))
	}

	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling the configuration YAML file: %v\n", err))
	}

	secretsProvider := secret.NewBitwardenProvider(secret.GetDefaultSecretsManagerAccessToken())

	downloader = NewServerDownloader(secretsProvider.Get(config.AccessTokenSecret).Unwrap(), time.Second*time.Duration(config.DownloadTimeoutInSeconds))
}
