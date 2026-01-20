package dto

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubnetVisibility string

const (
	Public   SubnetVisibility = "Public"
	Private  SubnetVisibility = "Private"
	Isolated SubnetVisibility = "Isolated"
)

type VpcInput struct {
	BaseNetwork        string
	RangeNetwork       string
	SubnetRangeNetwork string
	EnableNAT          bool
}

type VpcOut struct {
	VpcId           pulumi.StringOutput
	PublicSubnets   []pulumi.StringOutput
	PrivateSubnets  []pulumi.StringOutput
	IsolatedSubnets []pulumi.StringOutput
	InternetGateway pulumi.StringOutput
	NatGateway      *pulumi.StringOutput
}

type ElbModuleInput struct {
	// Naming
	Env           string // var.env
	ProjectPrefix string // var.tags.project_prefix (usato nei nomi risorse)

	// Network
	PrivateSubnetIds pulumi.StringArray // local.private_subnets_id (già calcolati: 2 subnets)
	VpcId            pulumi.StringInput // var.vpcId

	// Logs
	LogBucket string // module.lb_logs.s3.bucket

	// ALB
	AlbListenerPort     int     // var.alb_forward_listener_port
	AlbListenerProtocol string  // var.alb_forward_listener_protocol (es. "HTTPS")
	CertificateArn      *string // var.certificate_arn
	SslPolicy           *string // var.ssl_policy

	// NLB
	NlbListenerPort     int    // var.nlb_forward_listener_port
	NlbListenerProtocol string // var.nlb_forward_listener_protocol

	// TG (per NLB)
	NlbTgProtocol string // var.nlb_forward_target_group_protocol
	NlbTgPort     int    // var.nlb_forward_target_group_port

	// Listener Rule
	NlbListenerRulePriority int // var.nlb_listener_rule_priority

	// Tags base
	Tags pulumi.StringMap // local.tags
}

type ElbModuleResources struct {
	Alb            *lb.LoadBalancer
	AlbListener    *lb.Listener
	AlbHealthRule  *lb.ListenerRule
	Nlb            *lb.LoadBalancer
	NlbTargetGroup *lb.TargetGroup
	NlbListener    *lb.Listener
	TgAttachment   *lb.TargetGroupAttachment
	VpcLink        *apigateway.VpcLink
}
