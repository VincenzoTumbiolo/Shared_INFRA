package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/VincenzoTumbiolo/Infra-PlumiCommons-Package/infrastructure/config/tags"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetDefaultTags(env Environments) map[string]string {
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

func SubnetGroupNameTag(defaultTags pulumi.StringMap) pulumi.String {
	return pulumi.String(fmt.Sprintf("%sSubgroup%s", defaultTags["ProjectPrefix"], defaultTags["Environment"]))
}

func DBNameTag(defaultTags pulumi.StringMap, tpe string, index *int) pulumi.StringOutput {
	return pulumi.String(fmt.Sprintf("%sDatabase%s", defaultTags["ProjectPrefix"], defaultTags["Environment"])).ToStringOutput()
}

func VPCNameTag(defaultTags pulumi.StringMap) pulumi.String {
	return pulumi.String(fmt.Sprintf("%sVpc%s", defaultTags["ProjectPrefix"], defaultTags["Environment"]))
}

func IGWNameTag(defaultTags pulumi.StringMap) pulumi.String {
	return pulumi.String(fmt.Sprintf("%sIgw%s", defaultTags["ProjectPrefix"], defaultTags["Environment"]))
}

func SubnetNameTag(defaultTags pulumi.StringMap, visibility SubnetVisibility) pulumi.String {
	return pulumi.String(fmt.Sprintf("%s%s%sSubnet", defaultTags["ProjectPrefix"], defaultTags["Environment"], visibility))
}

func SGNameTag(defaultTags pulumi.StringMap, name *string) pulumi.String {
	return pulumi.String(fmt.Sprintf("%s%sSG%s", defaultTags["ProjectPrefix"], defaultTags["Environment"], capitalizeFirst(*name)))
}

func LambdaNameTag(defaultTags pulumi.StringMap, name *string) pulumi.String {
	if name != nil {
		return pulumi.String(fmt.Sprintf("%s%sLambda%s", defaultTags["ProjectPrefix"], defaultTags["Environment"], capitalizeFirst(*name)))
	}
	return pulumi.String(fmt.Sprintf("%s%sLambda", defaultTags["ProjectPrefix"], defaultTags["Environment"]))
}

func ApiGWNameTag(defaultTags pulumi.StringMap) pulumi.String {
	return pulumi.String(fmt.Sprintf("%s%sApiGW", defaultTags["Project"], defaultTags["Environment"]))
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}

	// Converti la prima lettera in maiuscolo
	first := string(unicode.ToUpper(rune(s[0])))

	// Il resto in minuscolo
	rest := ""
	if len(s) > 1 {
		rest = strings.ToLower(s[1:])
	}

	return first + rest
}
