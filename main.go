package main

import (
	"log"
	"shared_infra/core"
	"shared_infra/network"
	"shared_infra/load_balancer"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Ottieni la configurazione comune
		environments, mappedEnv, defaultTags := core.GetConfig()

		// La logica di provisioning è ora nel file deployer.go
		network, err := network.Infrastructure(ctx, environments, mappedEnv, defaultTags)
		if err != nil {
			log.Printf("Network deployment failed: %v", err)
		}
		if err:= load_balancer.Infrastructure(ctx, environments, mappedEnv, defaultTags)
		return err
	})
}
