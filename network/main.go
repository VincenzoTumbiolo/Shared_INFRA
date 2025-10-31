package main

import (
	"fmt"

	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/envs"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/tags"
	"github.com/VincenzoTumbiolo/Shared_INFRA_network/src"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(deploy)
}

func deploy(ctx *pulumi.Context) error {
	fmt.Println("Config")
	environments, environmentsMap, defaultTags := setConfig()

	vtechpulumiMod := vtech_aws.New(
		ctx,
		defaultTags,
		environmentsMap,
	)

	fmt.Println("Start deploy")
	res, err := src.Infrastructure(ctx, vtechpulumiMod, environments)
	if err != nil {
		return err
	}
	ctx.Export("vpcId", res.VpcId)
	for i, sub := range res.PrivateSubnets {
		ctx.Export(fmt.Sprintf("privateSubnet%d", i), sub)
	}
	fmt.Println("End deploy")

	return nil
}

func setConfig() (envs.Environments, pulumi.StringMap, pulumi.StringMap) {

	environments, mapEnv := envs.GetEnvironments()

	defaultTags := tags.GetDefaultTags(environments)

	return environments, pulumi.ToStringMap(mapEnv), pulumi.ToStringMap(defaultTags)
}
