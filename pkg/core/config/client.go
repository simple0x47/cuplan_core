package config

import (
	"fmt"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"go.uber.org/zap"
	"os"
	"time"
)

// Client provides a way to obtain configurations from remote locations.
type Client struct {
	logger             *zap.Logger
	host               string
	stage              string
	environment        string
	component          string
	workingPath        string
	downloadAgainAfter time.Duration
	downloader         Downloader
	extractor          Extractor
	provider           Provider
	lastDownload       time.Time
}

// NewClient creates a new instance of Client.
func NewClient(logger *zap.Logger,
	host string,
	stage string,
	environment string,
	component string,
	workingPath string,
	downloadAgainAfter time.Duration,
	downloader Downloader,
	extractor Extractor,
	provider Provider) *Client {
	client := new(Client)

	client.logger = logger
	client.host = host
	client.stage = stage
	client.environment = environment
	client.component = component
	client.workingPath = workingPath
	client.downloadAgainAfter = downloadAgainAfter
	client.downloader = downloader
	client.extractor = extractor
	client.provider = provider
	client.lastDownload = time.Unix(0, 0)

	return client
}

// Close deletes the working path.
func (c *Client) Close() {
	err := os.RemoveAll(c.workingPath)
	if err != nil {
		c.logger.Warn(fmt.Sprintf("Failed to remove working path '%s'.", c.workingPath))
	}
}

// Get retrieves the configuration located within the specified file and at the specified key.
// The different levels are separated by ':', i.e. "Root:Parent:Example".
func (c *Client) Get(filePath string, key string) core.Result[any, core.Error] {
	timeSinceLastDownload := time.Now().Sub(c.lastDownload)

	if timeSinceLastDownload > c.downloadAgainAfter || !_doesDirectoryExist(c.workingPath) {
		initResult := c.initializeConfig()

		if !initResult.IsOk() {
			return core.Err[any, core.Error](initResult.UnwrapErr())
		}
	}

	return c.provider.Get(filePath, key)
}

func (c *Client) initializeConfig() core.Result[core.Empty, core.Error] {
	err := os.MkdirAll(c.workingPath, os.ModePerm)
	if err != nil {
		c.logger.Info(fmt.Sprintf("Failed to create working path '%s'.", c.workingPath))
	}

	downloadResult := c.downloader.Download(c.host, c.stage, c.environment, c.component)

	if !downloadResult.IsOk() {
		return core.Err[core.Empty, core.Error](downloadResult.UnwrapErr())
	}

	c.lastDownload = time.Now()
	c.provider.CleanCache()
	return c.extractor.Extract(downloadResult.Unwrap(), c.workingPath)
}

func _doesDirectoryExist(directory string) bool {
	_, err := os.Stat(directory)

	return err == nil
}
