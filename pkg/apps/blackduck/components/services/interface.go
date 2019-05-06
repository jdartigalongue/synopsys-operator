package services

import "github.com/blackducksoftware/horizon/pkg/components"

type ServiceInterface interface {
	GetService() *components.Service
	// TODO add deployment, rc
}
