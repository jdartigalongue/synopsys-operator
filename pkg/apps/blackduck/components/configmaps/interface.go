package configmaps

import "github.com/blackducksoftware/horizon/pkg/components"

type ConfigMapInterface interface {
	GetCM() []*components.ConfigMap
}