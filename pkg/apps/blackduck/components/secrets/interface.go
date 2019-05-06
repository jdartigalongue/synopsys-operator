package secrets

import "github.com/blackducksoftware/horizon/pkg/components"

type SecretInterface interface {
	GetSecrets() []*components.Secret
}