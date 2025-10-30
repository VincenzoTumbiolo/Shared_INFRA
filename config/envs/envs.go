package envs

import (
	"log/slog"
	"reflect"

	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/config/configx"
)

type Environments struct {
	CreatedBy     string `env:"VAR_created_by" default:"services-Backend"`
	CreatedAt     string `env:"VAR_created_at" default:"14/02/2022"`
	ProjectPrefix string `env:"VAR_project_prefix" default:"shared"`
	Env           string `env:"VAR_env" default:"dev"`
	Region        string `env:"VAR_region" default:"eu-central-1"`
	AccountId     string `env:"VAR_account_id"`

	//LOAD BALANCER
	AlbListenerPort         int    `env:"VAR_albListenerPort" default:"443"`
	AlbListenerProtocol     string `env:"VAR_albListenerProtocol" default:"HTTPS"`
	CertificateId           string `env:"VAR_certificateArn"`
	NlbListenerPort         int    `env:"VAR_nlbListenerPort" default:"443"`
	NlbListenerProtocol     string `env:"VAR_nlbListenerProtocol" default:"TCP"`
	NlbTgPort               int    `env:"VAR_nlbTgPort" default:"443"`
	NlbTgProtocol           string `env:"VAR_nlbTgProtocol" default:"TCP"`
	NlbListenerRulePriority int    `env:"VAR_nlbListenerRulePriority" default:"49999"`

	// COMPUTED ENVS
	VpcId            *string `env:"SHARED_VPC_ID"`
	PrivateSubnetIds *string `env:"PRIVATE_SUBNETS"`
}

func GetEnvironments() (Environments, map[string]string) {
	var environments Environments

	if err := configx.Load(&environments); err != nil {
		slog.Error("Error loading environment variables", "error", err)
		panic(err)
	}

	resultStr := make(map[string]string)
	val := reflect.ValueOf(environments)

	for i := range val.NumField() {
		field := val.Field(i)
		resultStr[val.Type().Field(i).Name] = field.String()
	}

	return environments, resultStr
}
