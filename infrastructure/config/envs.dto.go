package config

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
