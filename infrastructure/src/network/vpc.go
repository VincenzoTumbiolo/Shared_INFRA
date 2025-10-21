package network

import (
	"fmt"

	vtechpulumi "github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/services/rest"
	"github.com/VincenzoTumbiolo/Shared_INFRA/config"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpcOut struct {
	VpcId           pulumi.StringOutput
	PublicSubnets   pulumi.StringArrayOutput
	PrivateSubnets  pulumi.StringArrayOutput
	IsolatedSubnets pulumi.StringArrayOutput
	InternetGateway pulumi.StringOutput
	NatGateway      pulumi.StringOutput
}

func NewNetwork(ctx *pulumi.Context, mod *vtechpulumi.RESTModule, baseName string, baseNetwork string, rangeNetwork string, subnetRangeNetwork string) error {
	cidr := fmt.Sprintf("%s.0.0/%s", baseNetwork, rangeNetwork)
	// --- VPC ---
	vpc, err := ec2.NewVpc(ctx, baseName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(cidr),
		EnableDnsSupport:   pulumi.Bool(true),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": config.VPCNameTag(mod.DefaultTags),
		},
	})
	if err != nil {
		return err
	}

	// --- Internet Gateway ---
	igw, err := ec2.NewInternetGateway(ctx, baseName+"-igw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": config.VPCNameTag(mod.DefaultTags),
		},
	})
	if err != nil {
		return err
	}

	// --- Public Subnets (2 AZs) ---
	azs := []string{"eu-central-1a", "eu-central-1b"}
	publicCidrs := []string{getSubnetCdir(baseNetwork, "0.0", subnetRangeNetwork), getSubnetCdir(baseNetwork, "1.0", subnetRangeNetwork)}
	privateCidrs := []string{getSubnetCdir(baseNetwork, "10.0", subnetRangeNetwork), getSubnetCdir(baseNetwork, "11.0", subnetRangeNetwork)}
	isolatedCidrs := []string{getSubnetCdir(baseNetwork, "20.0", subnetRangeNetwork), getSubnetCdir(baseNetwork, "21.0", subnetRangeNetwork)}

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
		Tags: pulumi.StringMap{"Name": pulumi.String(baseName + "-rt-public")},
	})
	if err != nil {
		return err
	}

	// --- EIP + NAT per private subnets ---
	eip, err := ec2.NewEip(ctx, baseName+"-eip-nat", &ec2.EipArgs{
		Domain: pulumi.String("vpc"),
		Tags:   pulumi.StringMap{"Name": pulumi.String(baseName + "-eip-nat")},
	})
	if err != nil {
		return err
	}

	nat, err := ec2.NewNatGateway(ctx, baseName+"-nat", &ec2.NatGatewayArgs{
		AllocationId: eip.ID(),
		SubnetId:     pulumi.StringOutput{}, // da riempire dopo con la prima public subnet
		Tags:         pulumi.StringMap{"Name": pulumi.String(baseName + "-nat")},
	}, pulumi.DependsOn([]pulumi.Resource{igw}))
	if err != nil {
		return err
	}

	// --- Creazione Subnet & Route Tables ---
	for i := range azs {
		// PUBLIC
		pub, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-public-%d", baseName, i+1), &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String(publicCidrs[i]),
			AvailabilityZone:    pulumi.String(azs[i]),
			MapPublicIpOnLaunch: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("%s-public-%d", baseName, i+1)),
				"Tier": pulumi.String("public"),
			},
		})
		if err != nil {
			return err
		}
		publicSubnets = append(publicSubnets, pub.ID().ToStringOutput())

		// Associa public subnet al route table public
		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("%s-rtassoc-public-%d", baseName, i+1), &ec2.RouteTableAssociationArgs{
			RouteTableId: rtPublic.ID(),
			SubnetId:     pub.ID(),
		})
		if err != nil {
			return err
		}

		// PRIVATE
		priv, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-private-%d", baseName, i+1), &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(privateCidrs[i]),
			AvailabilityZone: pulumi.String(azs[i]),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("%s-private-%d", baseName, i+1)),
				"Tier": pulumi.String("private"),
			},
		})
		if err != nil {
			return err
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
			Tags: pulumi.StringMap{"Name": pulumi.String(fmt.Sprintf("%s-rt-private-%d", baseName, i+1))},
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("%s-rtassoc-private-%d", baseName, i+1), &ec2.RouteTableAssociationArgs{
			RouteTableId: rtPriv.ID(),
			SubnetId:     priv.ID(),
		})
		if err != nil {
			return err
		}

		// ISOLATED
		iso, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-isolated-%d", baseName, i+1), &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(isolatedCidrs[i]),
			AvailabilityZone: pulumi.String(azs[i]),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("%s-isolated-%d", baseName, i+1)),
				"Tier": pulumi.String("isolated"),
			},
		})
		if err != nil {
			return err
		}
		isolatedSubnets = append(isolatedSubnets, iso.ID().ToStringOutput())
	}

	return nil
}

func getSubnetCdir(baseNetwork string, subnetNetwork string, rangeNetwork string) string {
	return fmt.Sprintf("%s.%s/%s", baseNetwork, subnetNetwork, rangeNetwork)
}
