package config

import (
	"context"
	"fmt"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"io"
	"net/http"
	"time"
)

type ServerDownloader struct {
	// TODO: Can expire, a renewal mechanism must be developed.
	accessToken     string
	downloadTimeout time.Duration
}

func NewServerDownloader(accessToken string, downloadTimeout time.Duration) *ServerDownloader {
	s := new(ServerDownloader)
	s.accessToken = accessToken
	s.downloadTimeout = downloadTimeout

	return s
}

func (s ServerDownloader) Download(host string, stage string, environment string, component string) core.Result[[]byte, core.Error] {
	url := fmt.Sprintf("%s/config?stage=%s&environment=%s&component=%s", host, stage, environment, component)

	ctx, cancel := context.WithTimeout(context.Background(), s.downloadTimeout)
	defer cancel()

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return core.Err[[]byte, core.Error](*core.NewError(core.ConfigurationRetrievalFailure, fmt.Sprintf("failed to get config from server: %s", err)))
	}

	request.WithContext(ctx)
	request.Header.Set("Authorization", "Bearer "+s.accessToken)

	client := http.Client{
		Timeout: s.downloadTimeout,
	}

	response, err := client.Do(request)

	if err != nil {
		return core.Err[[]byte, core.Error](*core.NewError(core.ConfigurationRetrievalFailure, fmt.Sprintf("failed to make GET request: %s", err)))
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return core.Err[[]byte, core.Error](*core.NewError(core.ConfigurationRetrievalFailure, fmt.Sprintf("received an unexpected status code %d", response.StatusCode)))
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return core.Err[[]byte, core.Error](*core.NewError(core.ConfigurationRetrievalFailure, fmt.Sprintf("failed to read response's body: %s", err)))
	}

	return core.Ok[[]byte, core.Error](body)
}