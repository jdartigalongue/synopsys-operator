package pvc

import "github.com/blackducksoftware/horizon/pkg/components"

type PVCInterface interface {
	GetPVCs() []*components.PersistentVolumeClaim
	// TODO add deployment, rc
}
