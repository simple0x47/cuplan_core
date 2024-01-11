package config

import (
	"github.com/google/uuid"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

const host = "https://simpleg.eu"
const stage = "dummy"
const environment = "development"
const component = "dummy"
const downloadAgainAfter = time.Second
const filePath = "application.yaml"
const configKey = "Parent:Child"
const value string = "1234abcd"

type MockDownloader struct {
	mock.Mock
}

type MockExtractor struct {
	mock.Mock
}

type MockProvider struct {
	mock.Mock
}

type ClientTestSuite struct {
	suite.Suite
	WorkingPath string
	PackageData []byte
	Downloader  *MockDownloader
	Extractor   *MockExtractor
	Provider    *MockProvider
	Client      *Client
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (m *MockDownloader) Download(host string, stage string, environment string, component string) core.Result[[]byte, core.Error] {
	args := m.Called(host, stage, environment, component)
	result, _ := args.Get(0).(core.Result[[]byte, core.Error])
	return core.Ok[[]byte, core.Error](result.Unwrap())
}

func (m *MockExtractor) Extract(packageData []byte, targetPath string) core.Result[core.Empty, core.Error] {
	m.Called(packageData, targetPath)
	return core.Ok[core.Empty, core.Error](core.Empty{})
}

func (m *MockProvider) Get(filePath string, key string) core.Result[any, core.Error] {
	args := m.Called(filePath, key)
	return core.Ok[any, core.Error](args.Get(0))
}

func (m *MockProvider) CleanCache() {
	m.Called()
}

func (c *ClientTestSuite) SetupTest() {
	c.WorkingPath = uuid.New().String()
	c.PackageData = make([]byte, 0)
	c.Downloader = new(MockDownloader)
	c.Extractor = new(MockExtractor)
	c.Provider = new(MockProvider)
	logger, _ := zap.NewDevelopment()
	c.Client = NewClient(logger, host, stage, environment, component, c.WorkingPath, downloadAgainAfter, c.Downloader, c.Extractor, c.Provider)

	c.Downloader.On("Download", host, stage, environment, component).Return(core.Ok[[]byte, core.Error](c.PackageData))
	c.Extractor.On("Extract", c.PackageData, c.WorkingPath).Return()
	c.Provider.On("Get", filePath, configKey).Return(value)
	c.Provider.On("CleanCache").Return()
}

func (c *ClientTestSuite) TestClient_Get_ReturnsExpectedValue() {
	defer c.Client.Close()
	result := c.Client.Get(filePath, configKey)

	c.AssertExpectedValue(result)
	c.AssertCompleteFlowExecutedTimes(1)
}

func (c *ClientTestSuite) TestClient_Get_AfterDownloadAgain_Downloads() {
	defer c.Client.Close()
	c.Client.Get(filePath, configKey)
	time.Sleep(time.Second * 2)

	secondResult := c.Client.Get(filePath, configKey)

	c.AssertExpectedValue(secondResult)
	c.AssertCompleteFlowExecutedTimes(2)
}

func (c *ClientTestSuite) TestClient_Close_RemovesWorkingPath() {
	c.Client.Close()

	_, err := os.Stat(c.WorkingPath)

	assert.True(c.T(), os.IsNotExist(err))
}

func (c *ClientTestSuite) AssertExpectedValue(result core.Result[any, core.Error]) {
	assert.True(c.T(), result.IsOk())
	assert.Equal(c.T(), value, result.Unwrap())
}

func (c *ClientTestSuite) AssertCompleteFlowExecutedTimes(times int) {
	c.Downloader.AssertNumberOfCalls(c.T(), "Download", times)
	c.Extractor.AssertNumberOfCalls(c.T(), "Extract", times)
	c.Provider.AssertNumberOfCalls(c.T(), "CleanCache", times)
	c.Provider.AssertNumberOfCalls(c.T(), "Get", times)
}
