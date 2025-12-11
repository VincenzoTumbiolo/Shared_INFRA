package envs

import (
	"log/slog"
	"reflect"
	"shared_infra/core/dto"

	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/config/configx"
)

func GetEnvironments() (dto.Environments, map[string]string) {
	var environments dto.Environments

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
