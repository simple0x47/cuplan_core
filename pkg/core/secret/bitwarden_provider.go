package secret

import (
	"encoding/json"
	"fmt"
	"github.com/simpleg-eu/cuplan_core/pkg/core"
	"os"
	"os/exec"
)

type BitwardenSecret struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	ProjectId      string `json:"project_id"`
	Key            string `json:"key"`
	Value          string `json:"value"`
	CreationDate   string `json:"creation_date"`
	RevisionDate   string `json:"revision_date"`
}

type BitwardenProvider struct {
	accessToken string
}

func NewBitwardenProvider(accessToken string) *BitwardenProvider {
	b := new(BitwardenProvider)
	b.accessToken = accessToken

	return b
}

func GetDefaultSecretsManagerAccessToken() string {
	return os.Getenv("SECRETS_MANAGER_ACCESS_TOKEN")
}

func (b BitwardenProvider) Get(secretId string) core.Result[string, core.Error] {
	cmd := exec.Command("bws", "get", "secret", secretId, "--access-token", b.accessToken)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return core.Err[string, core.Error](*core.NewError(core.CommandFailure, fmt.Sprintf("failed to get secret: %s", err.Error())))
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return core.Err[string, core.Error](*core.NewError(core.CommandFailure, fmt.Sprintf("secret provider exited with a non-zero code: %s", string(output))))
	}

	var secret BitwardenSecret

	err = json.Unmarshal(output, &secret)

	if err != nil {
		return core.Err[string, core.Error](*core.NewError(core.SerializationFailure, fmt.Sprintf("failed to extract secret from response: %s", err.Error())))
	}

	return core.Ok[string, core.Error](secret.Value)
}
