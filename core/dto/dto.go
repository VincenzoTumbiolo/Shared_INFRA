package dto

type Environments struct {
	CreatedBy     string       `env:"VAR_created_by" default:"services-Backend"`
	CreatedAt     string       `env:"VAR_created_at" default:"14/02/2022"`
	ProjectPrefix string       `env:"VAR_project_prefix" default:"shared"`
	Env           string       `env:"VAR_env" default:"dev"`
	Region        string       `env:"VAR_region" default:"eu-central-1"`
	AccountId     string       `env:"VAR_account_id"`
	PulumiAction  PulumiAction `env:"VAR_pulumi_action" default:"preview"`

	//LOAD BALANCER
	AlbListenerPort         int    `env:"VAR_albListenerPort" default:"443"`
	AlbListenerProtocol     string `env:"VAR_albListenerProtocol" default:"HTTPS"`
	CertificateId           string `env:"VAR_httpsCertificateId"`
	NlbListenerPort         int    `env:"VAR_nlbListenerPort" default:"443"`
	NlbListenerProtocol     string `env:"VAR_nlbListenerProtocol" default:"TCP"`
	NlbTgPort               int    `env:"VAR_nlbTgPort" default:"443"`
	NlbTgProtocol           string `env:"VAR_nlbTgProtocol" default:"TCP"`
	NlbListenerRulePriority int    `env:"VAR_nlbListenerRulePriority" default:"49999"`
}
