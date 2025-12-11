package load_balancer

import (
	"errors"
	"fmt"
	"strings"

	"shared_infra/core/dto"
	"shared_infra/load_balancer/connection"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"

	vtech_aws_dto "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/dto/aws"
	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Infrastructure(ctx *pulumi.Context, environments dto.Environments, environmentsMap pulumi.StringMap, defaultTags pulumi.StringMap, network dto.VpcOut) error {
	vtechpulumiMod := vtech_aws.New(
		ctx,
		defaultTags,
		environmentsMap,
	)

	fmt.Println("Start Deploy Load Balancer")

	baseName := fmt.Sprintf("%s-%s", environments.Env, environments.ProjectPrefix)
	bucketName := fmt.Sprintf("%s-lb-logs", baseName)
	certificateArn := fmt.Sprintf("arn:aws:acm:%s:%s:certificate/%s", environments.Region, environments.AccountId, environments.CertificateId)
	sslPolicy := "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"

	if network.VpcId == "empty" {
		return errors.New("ERROR | Missing Vpc Id")
	}
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
		Id: pulumi.StringRef(network.VpcId),
	})
	if err != nil {
		return err
	}

	_, err = connection.CreateElasticLoadBalancers(ctx, baseName, environments, vtechpulumiMod.DefaultTags, vpc, vtech_aws_dto.ElbModuleInput{
		Env:                     environments.Env,
		ProjectPrefix:           environments.ProjectPrefix,
		PrivateSubnetIds:        strings.Split(network.PrivateSubnetIds, ","),
		VpcId:                   network.VpcId,
		LogBucket:               bucketName,
		AlbListenerPort:         environments.AlbListenerPort,
		AlbListenerProtocol:     environments.AlbListenerProtocol,
		CertificateArn:          &certificateArn,
		SslPolicy:               &sslPolicy,
		NlbListenerPort:         environments.NlbListenerPort,
		NlbListenerProtocol:     environments.NlbListenerProtocol,
		NlbTgProtocol:           environments.NlbTgProtocol,
		NlbTgPort:               environments.NlbTgPort,
		NlbListenerRulePriority: environments.NlbListenerRulePriority,
		Tags:                    vtechpulumiMod.DefaultTags,
	})
	if err != nil {
		return err
	}
	fmt.Println("End deploy Load Balancer")

	return nil
}
