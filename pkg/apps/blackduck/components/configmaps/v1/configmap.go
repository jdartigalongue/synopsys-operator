package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/configmaps"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"strconv"
	"strings"
)

type configmap struct {
	blackduck *v1.Blackduck
}

func (c *configmap) GetCM() []*components.ConfigMap {

	var configMaps []*components.ConfigMap

	hubConfig := components.NewConfigMap(horizonapi.ConfigMapConfig{Namespace: c.blackduck.Spec.Namespace, Name: "blackduck-config"})
	hubData := map[string]string{
		"RUN_SECRETS_DIR": "/tmp/secrets",
		"HUB_VERSION":     c.blackduck.Spec.Version,
	}

	for _, value := range c.blackduck.Spec.Environs {
		values := strings.SplitN(value, ":", 2)
		if len(values) == 2 {
			mapKey := strings.TrimSpace(values[0])
			mapValue := strings.TrimSpace(values[1])
			if len(mapKey) > 0 && len(mapValue) > 0 {
				hubData[mapKey] = mapValue
			}
		}
	}

	// merge default and input environs
	environs := GetHubKnobs()
	hubData = util.MergeEnvMaps(hubData, environs)

	hubConfig.AddData(hubData)
	hubConfig.AddLabels(utils.GetVersionLabel("configmap",c.blackduck.Spec.Version))
	configMaps = append(configMaps, hubConfig)

	configMaps = append(configMaps, c.getPostgresCM())

	return configMaps
}

func (c *configmap) getPostgresCM() *components.ConfigMap {
	// DB
	hubDbConfig := components.NewConfigMap(horizonapi.ConfigMapConfig{Namespace: c.blackduck.Spec.Namespace, Name: "blackduck-db-config"})
	if c.blackduck.Spec.ExternalPostgres != nil {
		hubDbConfig.AddData(map[string]string{
			"HUB_POSTGRES_ADMIN": c.blackduck.Spec.ExternalPostgres.PostgresAdmin,
			"HUB_POSTGRES_USER":  c.blackduck.Spec.ExternalPostgres.PostgresUser,
			"HUB_POSTGRES_PORT":  strconv.Itoa(c.blackduck.Spec.ExternalPostgres.PostgresPort),
			"HUB_POSTGRES_HOST":  c.blackduck.Spec.ExternalPostgres.PostgresHost,
		})
	} else {
		hubDbConfig.AddData(map[string]string{
			"HUB_POSTGRES_ADMIN": "blackduck",
			"HUB_POSTGRES_USER":  "blackduck_user",
			"HUB_POSTGRES_PORT":  "5432",
			"HUB_POSTGRES_HOST":  "postgres",
		})
	}

	if c.blackduck.Spec.ExternalPostgres != nil {
		hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL": strconv.FormatBool(c.blackduck.Spec.ExternalPostgres.PostgresSsl)})
		if c.blackduck.Spec.ExternalPostgres.PostgresSsl {
			hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL_CERT_AUTH": "false"})
		}
	} else {
		hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL": "false"})
	}
	hubDbConfig.AddLabels(utils.GetVersionLabel("postgres",c.blackduck.Spec.Version))

	return hubDbConfig
}


func NewConfigmap(blackduck *v1.Blackduck) configmaps.ConfigMapInterface {
	return &configmap{blackduck: blackduck}
}

// GetHubKnobs returns the default environs
func GetHubKnobs() map[string]string {
	return map[string]string{
		"IPV4_ONLY":                         "0",
		"USE_ALERT":                         "0",
		"USE_BINARY_UPLOADS":                "0",
		"RABBIT_MQ_HOST":                    "rabbitmq",
		"RABBIT_MQ_PORT":                    "5671",
		"BROKER_URL":                        "amqps://rabbitmq/protecodesc",
		"BROKER_USE_SSL":                    "yes",
		"CFSSL":                             "cfssl:8888",
		"HUB_LOGSTASH_HOST":                 "logstash",
		"SCANNER_CONCURRENCY":               "1",
		"HTTPS_VERIFY_CERTS":                "yes",
		"RABBITMQ_DEFAULT_VHOST":            "protecodesc",
		"RABBITMQ_SSL_FAIL_IF_NO_PEER_CERT": "false",
		"CLIENT_CERT_CN":                    "binaryscanner",
		"ENABLE_SOURCE_UPLOADS":             "false",
		"DATA_RETENTION_IN_DAYS":            "180",
		"MAX_TOTAL_SOURCE_SIZE_MB":          "4000",
	}
}