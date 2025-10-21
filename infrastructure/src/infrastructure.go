package src

import (
	"fmt"

	vtechpulumi "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/rest"
	"github.com/VincenzoTumbiolo/Shared_INFRA/config"
	"github.com/VincenzoTumbiolo/Shared_INFRA/src/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Infrastructure sets up the infrastructure for the project.
func Infrastructure(ctx *pulumi.Context, mod *vtechpulumi.RESTModule, env config.Environments) error {
	baseName := fmt.Sprintf("%s-%s", env.Env, env.ProjectPrefix)

	err := network.NewNetwork(ctx, mod, baseName, "10.10", "16", "24")
	if err != nil {
		return err
	}

	return nil
}
