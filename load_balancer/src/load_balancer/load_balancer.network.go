package load_balancer

import (
	"fmt"

	vtech_aws_dto "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/dto/aws"
	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/aws/load_balancer"
	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/aws/s3"

	vtech_network "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/aws/network"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/dto"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/envs"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elb"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ====== Modulo principale ======

func CreateElasticLoadBalancers(ctx *pulumi.Context, baseName string, env envs.Environments, defaultTags pulumi.StringMapInput, vpc *ec2.LookupVpcResult, in vtech_aws_dto.ElbModuleInput) (*dto.ElbModuleResources, error) {

	sgName := fmt.Sprintf("%s-lb-sg", baseName)

	sg, err := vtech_network.CreateSecurityGroup(ctx, sgName, vtech_aws_dto.SecurityGroupArgs{
		VpcID: &env.VpcId,
		Ingress: []vtech_aws_dto.SecurityGroupRule{
			{
				Protocol:    "TCP",
				FromPort:    env.AlbListenerPort,
				ToPort:      env.AlbListenerPort,
				CidrBlock:   vpc.CidrBlock,
				Description: "",
			},
		},
		Egress: []vtech_aws_dto.SecurityGroupRule{
			{
				Protocol:    "-1",
				FromPort:    0,
				ToPort:      0,
				CidrBlock:   "0.0.0.0/0",
				Description: "",
			},
		},
		Tags: defaultTags,
	})
	account, err := elb.GetServiceAccount(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	if _, err = s3.CreateS3Bucket(ctx, vtech_aws_dto.S3BucketInput{
		Name:      in.LogBucket,
		Versioned: true, // var.versioned
		PolicyStatements: []vtech_aws_dto.S3PolicyStatement{
			{
				Sid:                  "0",
				Actions:              []string{"s3:PutObject"},
				Effect:               "Allow",
				Resources:            []string{fmt.Sprintf("arn:aws:s3:::%s/*", in.LogBucket)},
				PrincipalType:        "AWS",
				PrincipalIdentifiers: []string{account.Arn},
			},
			{
				Sid:                  "1",
				Actions:              []string{"s3:PutObject"},
				Effect:               "Allow",
				Resources:            []string{fmt.Sprintf("arn:aws:s3:::%s/*", in.LogBucket)},
				PrincipalType:        "Service",
				PrincipalIdentifiers: []string{"delivery.logs.amazonaws.com"},
			},
			{
				Sid:                  "2",
				Actions:              []string{"s3:GetBucketAcl"},
				Effect:               "Allow",
				Resources:            []string{fmt.Sprintf("arn:aws:s3:::%s", in.LogBucket)},
				PrincipalType:        "Service",
				PrincipalIdentifiers: []string{"delivery.logs.amazonaws.com"},
			},
		},
		Tags: in.Tags,
	}); err != nil {
		return nil, err
	}

	sgID := sg.ID()

	// ---- main_alb (module lb_service) ----
	alb, err := load_balancer.CreateService(ctx, vtech_aws_dto.LoadBalancerInput{
		LbName:            fmt.Sprintf("%s-%s-alb", in.Env, in.ProjectPrefix),
		LbType:            "application",
		LbSecurityGroupId: &sgID,
		LbSubnetIds:       in.PrivateSubnetIds,
		LogBucket:         in.LogBucket,
		Tags:              in.Tags,
	})
	if err != nil {
		return nil, err
	}

	var cert pulumi.StringPtrInput
	if in.CertificateArn != nil {
		cert = pulumi.StringPtr(*in.CertificateArn)
	}

	var sslPolicy pulumi.StringPtrInput
	if in.SslPolicy != nil {
		sslPolicy = pulumi.StringPtr(*in.SslPolicy)
	}

	albListener, err := lb.NewListener(ctx, "main_alb_listener", &lb.ListenerArgs{
		LoadBalancerArn: alb.Arn,
		Port:            pulumi.Int(in.AlbListenerPort),
		Protocol:        pulumi.String(in.AlbListenerProtocol),
		CertificateArn:  cert,
		SslPolicy:       sslPolicy,
		DefaultActions: lb.ListenerDefaultActionArray{
			&lb.ListenerDefaultActionArgs{
				Type: pulumi.String("fixed-response"),
				FixedResponse: &lb.ListenerDefaultActionFixedResponseArgs{
					ContentType: pulumi.String("text/plain"),
					MessageBody: pulumi.StringPtr("FORWARD ERROR"),
					StatusCode:  pulumi.String("400"),
				},
			},
		},
		Tags: in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// ---- aws_lb_listener_rule.healtcheck_listener_rule ----
	healthRule, err := lb.NewListenerRule(ctx, "healtcheck_listener_rule", &lb.ListenerRuleArgs{
		ListenerArn: albListener.Arn,
		Priority:    pulumi.Int(in.NlbListenerRulePriority),
		Actions: lb.ListenerRuleActionArray{
			&lb.ListenerRuleActionArgs{
				Type: pulumi.String("fixed-response"),
				FixedResponse: &lb.ListenerRuleActionFixedResponseArgs{
					ContentType: pulumi.String("text/plain"),
					MessageBody: pulumi.String("healt check passed"),
					StatusCode:  pulumi.String("200"),
				},
			},
		},
		Conditions: lb.ListenerRuleConditionArray{
			&lb.ListenerRuleConditionArgs{
				PathPattern: &lb.ListenerRuleConditionPathPatternArgs{
					Values: pulumi.StringArray{pulumi.String("/healtcheck")},
				},
			},
		},
		Tags: in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// ---- main_nlb (module lb_service) ----
	nlbRes, err := load_balancer.CreateService(ctx, vtech_aws_dto.LoadBalancerInput{
		LbName:      fmt.Sprintf("%s-%s-nlb", in.Env, in.ProjectPrefix),
		LbType:      "network",
		LbSubnetIds: in.PrivateSubnetIds,
		LogBucket:   in.LogBucket,
		Tags:        in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// ---- main_nlb_target_group (module lb_target_group) ----
	tg, err := load_balancer.CreateTargetGroup(ctx, vtech_aws_dto.TargetGroupInput{
		Name:       fmt.Sprintf("%s-%s-target-group", in.Env, in.ProjectPrefix),
		Port:       in.NlbTgPort,
		TargetType: "alb",
		Protocol:   in.NlbTgProtocol,
		VpcId:      in.VpcId,

		HealthCheckPath:               "/healtcheck",
		HealthCheckPort:               fmt.Sprintf("%d", in.AlbListenerPort),
		HealthCheckProtocol:           in.AlbListenerProtocol,
		HealthCheckHealthyThreshold:   0, // → default AWS
		HealthCheckUnhealthyThreshold: 0, // → default AWS
		HealthCheckMatcher:            "200-399",

		Tags: in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// ---- aws_lb_listener.main_nlb_listener (forward to TG) ----
	nlbListener, err := lb.NewListener(ctx, "main_nlb_listener", &lb.ListenerArgs{
		LoadBalancerArn: nlbRes.Arn,
		Port:            pulumi.Int(in.NlbListenerPort),
		Protocol:        pulumi.String(in.NlbListenerProtocol),
		DefaultActions: lb.ListenerDefaultActionArray{
			&lb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				TargetGroupArn: tg.Arn,
			},
		},
		Tags: in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// ---- aws_lb_target_group_attachment.nlb ----
	tgAttachment, err := lb.NewTargetGroupAttachment(ctx, "nlb_attachment", &lb.TargetGroupAttachmentArgs{
		TargetGroupArn: tg.Arn,
		TargetId:       alb.ID(),                       // ALB come target dell’NLB TG
		Port:           pulumi.Int(in.AlbListenerPort), // porta ALB
	})
	if err != nil {
		return nil, err
	}

	vpcLink, err := apigateway.NewVpcLink(ctx, "vpcLink", &apigateway.VpcLinkArgs{
		Name:      pulumi.String(fmt.Sprintf("%s-%s-vpc-link", in.Env, in.ProjectPrefix)),
		TargetArn: nlbRes.Arn,
		Tags:      in.Tags,
	})
	if err != nil {
		return nil, err
	}

	return &dto.ElbModuleResources{
		Alb:            alb,
		AlbListener:    albListener,
		AlbHealthRule:  healthRule,
		Nlb:            nlbRes,
		NlbTargetGroup: tg,
		NlbListener:    nlbListener,
		TgAttachment:   tgAttachment,
		VpcLink:        vpcLink,
	}, nil
}
