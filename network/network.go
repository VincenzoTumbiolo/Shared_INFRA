package network

import (
	"fmt"

	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"

	"shared_infra/core/dto"
	"shared_infra/network/vpc"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Infrastructure(ctx *pulumi.Context, environments dto.Environments, environmentsMap pulumi.StringMap, defaultTags pulumi.StringMap) (*dto.VpcOut, error) {

	vtechpulumiMod := vtech_aws.New(
		ctx,
		defaultTags,
		environmentsMap,
	)

	fmt.Println("Start deploy")
	baseName := fmt.Sprintf("%s-%s", environments, environments.ProjectPrefix)

	vpc, err := vpc.NewNetwork(ctx, vtechpulumiMod, baseName, "10.10", "16", "24")
	if err != nil {
		return nil, err
	}

	ctx.Export("vpcId", vpc.VpcId)
	for i, sub := range vpc.PrivateSubnets {
		ctx.Export(fmt.Sprintf("privateSubnet%d", i), sub)
	}
	fmt.Println("End deploy")

	return vpc, nil
}
