package core

import (
	"shared_infra/core/dto"
	"shared_infra/core/envs"
	"shared_infra/core/tags"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetConfig() (dto.Environments, pulumi.StringMap, pulumi.StringMap) {

	environments, mapEnv := envs.GetEnvironments()

	defaultTags := tags.GetDefaultTags(environments)

	return environments, pulumi.ToStringMap(mapEnv), pulumi.ToStringMap(defaultTags)
}
