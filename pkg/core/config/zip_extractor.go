package config

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/simpleg-eu/cuplan-core/pkg/core"
	"io"
	"os"
	"path/filepath"
)

type ZipExtractor struct {
}

func NewZipExtractor() *ZipExtractor {
	return new(ZipExtractor)
}

func (z ZipExtractor) Extract(packageData []byte, targetPath string) core.Result[core.Empty, core.Error] {
	err := os.MkdirAll(targetPath, os.ModePerm)

	if err != nil {
		return core.Err[core.Empty, core.Error](*core.NewError(core.ExtractionFailure, fmt.Sprintf("failed to create target directory: %s", err)))
	}

	reader := bytes.NewReader(packageData)
	zipReader, err := zip.NewReader(reader, int64(len(packageData)))

	if err != nil {
		return core.Err[core.Empty, core.Error](*core.NewError(core.ExtractionFailure, fmt.Sprintf("failed to unzip package data: %s", err)))
	}

	for _, file := range zipReader.File {
		rc, err := file.Open()

		if err != nil {
			return core.Err[core.Empty, core.Error](*core.NewError(core.ExtractionFailure, fmt.Sprintf("failed to open file within package data: %s", err)))
		}

		defer rc.Close()

		if file.FileInfo().IsDir() {
			continue
		}

		extractedFilePath := filepath.Join(targetPath, file.Name)

		err = os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm)

		if err != nil {
			return core.Err[core.Empty, core.Error](*core.NewError(core.ExtractionFailure, fmt.Sprintf("failed to create sub-directory for package data's extraction: %s", err)))
		}

		extractedFile, err := os.Create(extractedFilePath)

		if err != nil {
			return core.Err[core.Empty, core.Error](*core.NewError(core.ExtractionFailure, fmt.Sprintf("failed to create extracted file: %s", err)))
		}
		defer extractedFile.Close()

		_, err = io.Copy(extractedFile, rc)

		if err != nil {
			return core.Err[core.Empty, core.Error](*core.NewError(core.ConfigurationRetrievalFailure, fmt.Sprintf("failed to copy package data's file content: %s", err)))
		}
	}

	return core.Ok[core.Empty, core.Error](core.Empty{})
}
