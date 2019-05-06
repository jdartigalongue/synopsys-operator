package rc

import "github.com/blackducksoftware/horizon/pkg/components"

type ReplicationControllerInterface interface {
	GetRc() *components.ReplicationController
	// TODO add deployment, rc
}
