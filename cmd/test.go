package main

import (
	v12 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/creater/v1"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	v1.NewCreater(0, nil, nil).Ensure(&v12.Blackduck{
		ObjectMeta: v13.ObjectMeta{
			Name:            "test",
		},
		Spec: v12.BlackduckSpec{
			Namespace:     "myns",
			Size:          "small",
			Version:       "2019.4.1",
			ExposeService: "",
			DbPrototype:   "",
			ExternalPostgres: &v12.PostgresExternalDBConfig{
				PostgresHost:          "",
				PostgresPort:          0,
				PostgresAdmin:         "",
				PostgresUser:          "",
				PostgresSsl:           false,
				PostgresAdminPassword: "",
				PostgresUserPassword:  "",
			},
			PVCStorageClass:   "",
			LivenessProbes:    false,
			ScanType:          "",
			PersistentStorage: true,
			PVC:               nil,
			CertificateName:   "",
			Certificate:       "",
			CertificateKey:    "",
			ProxyCertificate:  "",
			AuthCustomCA:      "",
			Type:              "",
			DesiredState:      "",
			Environs:          nil,
			ImageRegistries:   nil,
			ImageUIDMap:       nil,
			LicenseKey:        "",
			RegistryConfiguration: v12.RegistryConfiguration{
				Registry:    "",
				Namespace:   "",
				PullSecrets: nil,
			},
		},
	})
}