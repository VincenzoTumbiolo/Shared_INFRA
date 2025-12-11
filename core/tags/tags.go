package tags

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"shared_infra/core/dto"

	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/config/tags"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetDefaultTags(env dto.Environments) map[string]string {
	e := tags.DefaultTags{
		CreationDate:         env.CreatedAt,
		Environment:          "Dev",
		InfraManagedBy:       "VincenzoTumbiolo",
		ApplicationManagedBy: "VincenzoTumbiolo",
		ProjectOwner:         "VTechStudio",
		Project:              "SharedInfra",
		CreatedBy:            env.CreatedBy,
		ProjectPrefix:        strings.ReplaceAll(env.ProjectPrefix, "-", "_"),
		Suffix:               "D",
	}

	var result map[string]string
	data, err := json.Marshal(e)
	if err != nil {
		slog.Error("Error marshalling Environments struct", "error", err)
		panic(err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		slog.Error("Error unmarshalling Environments struct", "error", err)
		panic(err)
	}
	return result
}

func ServiceNameTag(service string, name string) pulumi.String {
	first := string(unicode.ToUpper(rune(name[0])))

	rest := ""
	if len(name) > 1 {
		rest = strings.ToLower(name[1:])
	}

	return pulumi.String(fmt.Sprintf("%s-%s", first+rest, service))
}
