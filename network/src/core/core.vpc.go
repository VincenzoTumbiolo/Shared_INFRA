package core

import (
	"fmt"

	vtech_aws "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/modules/aws"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/dto"
	"github.com/VincenzoTumbiolo/Shared_INFRA_config/tags"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewNetwork(ctx *pulumi.Context, mod *vtech_aws.AWSModule, baseName string, input dto.VpcInput) (*dto.VpcOut, error) {
	cidr := fmt.Sprintf("%s.0.0/%s", input.BaseNetwork, input.RangeNetwork)
	// --- VPC ---
	vpc, err := ec2.NewVpc(ctx, baseName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(cidr),
		EnableDnsSupport:   pulumi.Bool(true),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": tags.ServiceNameTag("Vpc", baseName),
		},
	})
	if err != nil {
		return nil, err
	}

	// --- Internet Gateway ---
	igw, err := ec2.NewInternetGateway(ctx, baseName+"-igw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": tags.ServiceNameTag("Igw", baseName),
		},
	})
	if err != nil {
		return nil, err
	}

	// --- Public Subnets (2 AZs) ---
	azs := []string{"eu-central-1a", "eu-central-1b"}
	publicCidrs := []string{getSubnetCdir(input.BaseNetwork, "0.0", input.SubnetRangeNetwork), getSubnetCdir(input.BaseNetwork, "1.0", input.SubnetRangeNetwork)}
	privateCidrs := []string{getSubnetCdir(input.BaseNetwork, "10.0", input.SubnetRangeNetwork), getSubnetCdir(input.BaseNetwork, "11.0", input.SubnetRangeNetwork)}
	isolatedCidrs := []string{getSubnetCdir(input.BaseNetwork, "20.0", input.SubnetRangeNetwork), getSubnetCdir(input.BaseNetwork, "21.0", input.SubnetRangeNetwork)}

	publicSubnets := []pulumi.StringOutput{}
	privateSubnets := []pulumi.StringOutput{}
	isolatedSubnets := []pulumi.StringOutput{}

	// Route table per public
	rtPublic, err := ec2.NewRouteTable(ctx, baseName+"-rt-public", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
		},
		Tags: pulumi.StringMap{
			"Name": tags.ServiceNameTag("PublicRt", baseName),
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range azs {
		// PUBLIC
		pub, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-public-%d", baseName, i+1), &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String(publicCidrs[i]),
			AvailabilityZone:    pulumi.String(azs[i]),
			MapPublicIpOnLaunch: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"Name": tags.ServiceNameTag(fmt.Sprintf("PublicSub%d", i+1), baseName),
				"Tier": pulumi.String("public"),
			},
		})
		if err != nil {
			return nil, err
		}
		publicSubnets = append(publicSubnets, pub.ID().ToStringOutput())

		// Associa public subnet al route table public
		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("%s-rtassoc-public-%d", baseName, i+1), &ec2.RouteTableAssociationArgs{
			RouteTableId: rtPublic.ID(),
			SubnetId:     pub.ID(),
		})
		if err != nil {
			return nil, err
		}
	}
	var nat *ec2.NatGateway

	if input.EnableNAT {
		// --- EIP + NAT per private subnets ---
		eip, err := ec2.NewEip(ctx, baseName+"-eip", &ec2.EipArgs{
			Domain: pulumi.String("vpc"),
			Tags: pulumi.StringMap{
				"Name": tags.ServiceNameTag("Eip", baseName),
			},
		})
		if err != nil {
			return nil, err
		}
		nat, err := ec2.NewNatGateway(ctx, baseName+"-nat", &ec2.NatGatewayArgs{
			AllocationId: eip.ID(),
			SubnetId:     publicSubnets[0],
			Tags: pulumi.StringMap{
				"Name": tags.ServiceNameTag("Nat", baseName),
			},
		}, pulumi.DependsOn([]pulumi.Resource{igw}))
		if err != nil {
			return nil, err
		}

		for i := range azs {
			// PRIVATE
			priv, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-private-%d", baseName, i+1), &ec2.SubnetArgs{
				VpcId:            vpc.ID(),
				CidrBlock:        pulumi.String(privateCidrs[i]),
				AvailabilityZone: pulumi.String(azs[i]),
				Tags: pulumi.StringMap{
					"Name": tags.ServiceNameTag(fmt.Sprintf("PrivateSub%d", i+1), baseName),
					"Tier": pulumi.String("private"),
				},
			})
			if err != nil {
				return nil, err
			}
			privateSubnets = append(privateSubnets, priv.ID().ToStringOutput())

			rtPriv, err := ec2.NewRouteTable(ctx, fmt.Sprintf("%s-rt-private-%d", baseName, i+1), &ec2.RouteTableArgs{
				VpcId: vpc.ID(),
				Routes: ec2.RouteTableRouteArray{
					ec2.RouteTableRouteArgs{
						CidrBlock:    pulumi.String("0.0.0.0/0"),
						NatGatewayId: nat.ID(),
					},
				},
				Tags: pulumi.StringMap{
					"Name": tags.ServiceNameTag(fmt.Sprintf("PrivateRt%d", i+1), baseName),
				},
			})
			if err != nil {
				return nil, err
			}

			_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("%s-rtassoc-private-%d", baseName, i+1), &ec2.RouteTableAssociationArgs{
				RouteTableId: rtPriv.ID(),
				SubnetId:     priv.ID(),
			})
			if err != nil {
				return nil, err
			}

			// ISOLATED
			iso, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-isolated-%d", baseName, i+1), &ec2.SubnetArgs{
				VpcId:            vpc.ID(),
				CidrBlock:        pulumi.String(isolatedCidrs[i]),
				AvailabilityZone: pulumi.String(azs[i]),
				Tags: pulumi.StringMap{
					"Name": tags.ServiceNameTag(fmt.Sprintf("IsolatedSub%d", i+1), baseName),
					"Tier": pulumi.String("isolated"),
				},
			})
			if err != nil {
				return nil, err
			}
			isolatedSubnets = append(isolatedSubnets, iso.ID().ToStringOutput())
		}
	}

	return &dto.VpcOut{
		VpcId:           igw.VpcId,               // pulumi.StringOutput
		PublicSubnets:   publicSubnets,           // pulumi.StringArrayOutput
		PrivateSubnets:  privateSubnets,          // pulumi.StringArrayOutput
		IsolatedSubnets: isolatedSubnets,         // pulumi.StringArrayOutput
		InternetGateway: igw.Arn,                 // pulumi.StringOutput
		NatGateway:      &nat.NetworkInterfaceId, // pulumi.StringOutput
	}, nil
}

func getSubnetCdir(baseNetwork string, subnetNetwork string, rangeNetwork string) string {
	return fmt.Sprintf("%s.%s/%s", baseNetwork, subnetNetwork, rangeNetwork)
}
