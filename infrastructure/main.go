package main

import (
	"fmt"

	vtechpulumi "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/rest"
	"github.com/VincenzoTumbiolo/Shared_INFRA/config"
	"github.com/VincenzoTumbiolo/Shared_INFRA/src"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(deploy)
}

func deploy(ctx *pulumi.Context) error {
	fmt.Println("Config")
	environments, environmentsMap, defaultTags := setConfig()

	vtechpulumiMod := vtechpulumi.New(
		ctx,
		defaultTags,
		environmentsMap,
	)

	fmt.Println("Start deploy")
	err := src.Infrastructure(ctx, vtechpulumiMod, environments)
	if err != nil {
		return err
	}
	fmt.Println("End deploy")

	return nil
}

func setConfig() (config.Environments, pulumi.StringMap, pulumi.StringMap) {

	environments, mapEnv := config.GetEnvironments()

	defaultTags := config.GetDefaultTags(environments)

	return environments, pulumi.ToStringMap(mapEnv), pulumi.ToStringMap(defaultTags)
}
