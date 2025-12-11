package dto

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpcOut struct {
	VpcId           pulumi.StringOutput
	PublicSubnets   []pulumi.StringOutput
	PrivateSubnets  []pulumi.StringOutput
	IsolatedSubnets []pulumi.StringOutput
	InternetGateway pulumi.StringOutput
	NatGateway      pulumi.StringOutput
}
