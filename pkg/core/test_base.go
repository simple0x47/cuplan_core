package core

import (
	"strings"
)

func GetTestDataPath(testFile string) Result[string, Error] {
	index := strings.Index(testFile, "/pkg/")

	if index == -1 {
		return Err[string, Error](*NewError(MissingFilePath, "Could not find '/pkg/' in test file's directory."))
	}

	result := testFile[:index] + "/test/data/" + testFile[index:len(testFile)-3] + "/"

	return Ok[string, Error](result)
}
