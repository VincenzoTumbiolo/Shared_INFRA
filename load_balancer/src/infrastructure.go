package src

import (
	"errors"
	"fmt"
	"strings"

	vtech_aws_dto "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/dto/aws"
	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/envs"
	"github.com/VincenzoTumbiolo/Shared_INFRA_load_balancer/src/load_balancer"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Infrastructure sets up the infrastructure for the project.
func Infrastructure(ctx *pulumi.Context, mod *vtech_aws.AWSModule, env envs.Environments) error {
	baseName := fmt.Sprintf("%s-%s", env.Env, env.ProjectPrefix)
	bucketName := fmt.Sprintf("%s-lb-logs", baseName)
	certificateArn := fmt.Sprintf("arn:aws:acm:%s:%s:certificate/%s", env.Region, env.AccountId, env.CertificateId)
	sslPolicy := "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"

	if env.VpcId == nil {
		return errors.New("ERROR | Missing Vpc Id")
	}
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
		Id: pulumi.StringRef(*env.VpcId), // o un'altra logica
	})
	if err != nil {
		return err
	}

	_, err = load_balancer.CreateElasticLoadBalancers(ctx, baseName, env, mod.DefaultTags, vpc, vtech_aws_dto.ElbModuleInput{
		Env:                     env.Env,
		ProjectPrefix:           env.ProjectPrefix,
		PrivateSubnetIds:        strings.Split(*env.PrivateSubnetIds, ","),
		VpcId:                   *env.VpcId,
		LogBucket:               bucketName,
		AlbListenerPort:         env.AlbListenerPort,
		AlbListenerProtocol:     env.AlbListenerProtocol,
		CertificateArn:          &certificateArn,
		SslPolicy:               &sslPolicy,
		NlbListenerPort:         env.NlbListenerPort,
		NlbListenerProtocol:     env.NlbListenerProtocol,
		NlbTgProtocol:           env.NlbTgProtocol,
		NlbTgPort:               env.NlbTgPort,
		NlbListenerRulePriority: env.NlbListenerRulePriority,
		Tags:                    mod.DefaultTags,
	})
	if err != nil {
		return err
	}

	return nil
}
