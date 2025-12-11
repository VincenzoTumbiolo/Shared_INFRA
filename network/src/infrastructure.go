package src

import (
	"fmt"

	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/dto"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/envs"
	"github.com/VincenzoTumbiolo/Shared_INFRA_network/src/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Infrastructure sets up the infrastructure for the project.
func Infrastructure(ctx *pulumi.Context, mod *vtech_aws.AWSModule, env envs.Environments) (*dto.VpcOut, error) {
	baseName := fmt.Sprintf("%s-%s", env.Env, env.ProjectPrefix)

	vpc, err := core.NewNetwork(ctx, mod, baseName, "10.10", "16", "24")
	if err != nil {
		return nil, err
	}

	return vpc, nil
}
