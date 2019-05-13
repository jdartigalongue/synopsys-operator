package v1

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/secrets"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/blackduck/util"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	util2 "github.com/blackducksoftware/synopsys-operator/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type secret struct {
	blackduck *v1.Blackduck
	Config          *protoform.Config
	KubeClient      *kubernetes.Clientset
}

func NewSecret(blackduck *v1.Blackduck, config *protoform.Config, kubeClient *kubernetes.Clientset) secrets.SecretInterface {
	return &secret{blackduck: blackduck, Config: config, KubeClient: kubeClient}
}

func (s *secret) GetSecrets() []*components.Secret {
	var secretArr []*components.Secret


	if len(s.blackduck.Spec.ProxyCertificate) > 0 {
		cert, err := s.stringToCertificate(s.blackduck.Spec.ProxyCertificate)
		if err != nil {
			logrus.Warnf("The proxy certificate provided is invalid")
		} else {
			logrus.Debugf("Adding Proxy certificate with SN: %x", cert.SerialNumber)
			proxyCertificateSecret := components.NewSecret(horizonapi.SecretConfig{Namespace: s.blackduck.Spec.Namespace, Name: "blackduck-proxy-certificate", Type: horizonapi.SecretTypeOpaque})
			proxyCertificateSecret.AddData(map[string][]byte{"HUB_PROXY_CERT_FILE": []byte(s.blackduck.Spec.ProxyCertificate)})
			proxyCertificateSecret.AddLabels(utils.GetVersionLabel("secret", s.blackduck.Spec.Version))
			secretArr = append(secretArr, proxyCertificateSecret)
		}
	}

	if len(s.blackduck.Spec.AuthCustomCA) > 0 {
		cert, err := s.stringToCertificate(s.blackduck.Spec.AuthCustomCA)
		if err != nil {
			logrus.Warnf("The Auth Custom CA provided is invalid")
		} else {
			logrus.Debugf("Adding The Auth Custom CA with SN: %x", cert.SerialNumber)
			authCustomCASecret := components.NewSecret(horizonapi.SecretConfig{Namespace: s.blackduck.Spec.Namespace, Name: "auth-custom-ca", Type: horizonapi.SecretTypeOpaque})
			authCustomCASecret.AddData(map[string][]byte{"AUTH_CUSTOM_CA": []byte(s.blackduck.Spec.AuthCustomCA)})
			authCustomCASecret.AddLabels(utils.GetVersionLabel("secret", s.blackduck.Spec.Version))
			secretArr = append(secretArr, authCustomCASecret)
		}
	}

	secretArr = append(secretArr, s.getPostgresSecret())

	return secretArr
}

func (s *secret) stringToCertificate(certificate string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}


// GetPostgresSecret will return the postgres secret
func (s *secret) getPostgresSecret() *components.Secret {

	// Postgres secret / CM
	var adminPassword, userPassword string
	if s.blackduck.Spec.ExternalPostgres != nil {
		adminPassword = s.blackduck.Spec.ExternalPostgres.PostgresAdminPassword
		userPassword = s.blackduck.Spec.ExternalPostgres.PostgresAdminPassword
	} else {
		var err error
		adminPassword, userPassword, _, err = util.GetDefaultPasswords(s.KubeClient, s.Config.Namespace)
		if err != nil {
			return nil
		}
	}

	hubSecret := components.NewSecret(horizonapi.SecretConfig{Namespace: s.blackduck.Spec.Namespace, Name: "db-creds", Type: horizonapi.SecretTypeOpaque})

	if s.blackduck.Spec.ExternalPostgres != nil {
		hubSecret.AddData(map[string][]byte{"HUB_POSTGRES_ADMIN_PASSWORD_FILE": []byte(s.blackduck.Spec.ExternalPostgres.PostgresAdminPassword), "HUB_POSTGRES_USER_PASSWORD_FILE": []byte(s.blackduck.Spec.ExternalPostgres.PostgresUserPassword)})
	} else {
		hubSecret.AddData(map[string][]byte{"HUB_POSTGRES_ADMIN_PASSWORD_FILE": []byte(adminPassword), "HUB_POSTGRES_USER_PASSWORD_FILE": []byte(userPassword)})
	}
	hubSecret.AddLabels(utils.GetVersionLabel("postgres", s.blackduck.Spec.Version))

	// TODO add certificate + SEAL_KEY

	return hubSecret
}


func (s *secret) getTLSCertKeyOrCreate() (string, string, error) {
	if len(s.blackduck.Spec.Certificate) > 0 && len(s.blackduck.Spec.CertificateKey) > 0 {
		return s.blackduck.Spec.Certificate, s.blackduck.Spec.CertificateKey, nil
	}

	// Cert copy
	if len(s.blackduck.Spec.CertificateName) > 0 && !strings.EqualFold(s.blackduck.Spec.CertificateName, "default") {
		secret, err := util2.GetSecret(s.KubeClient, s.blackduck.Spec.CertificateName, "blackduck-certificate")
		if err == nil {
			cert, certok := secret.Data["WEBSERVER_CUSTOM_CERT_FILE"]
			key, keyok := secret.Data["WEBSERVER_CUSTOM_KEY_FILE"]
			if certok && keyok {
				return string(cert), string(key), nil
			}
		}
	}

	// default cert
	secret, err := util2.GetSecret(s.KubeClient, s.Config.Namespace, "blackduck-certificate")
	if err == nil {
		data := secret.Data
		if len(data) >= 2 {
			cert, certok := secret.Data["WEBSERVER_CUSTOM_CERT_FILE"]
			key, keyok := secret.Data["WEBSERVER_CUSTOM_KEY_FILE"]
			if !certok || !keyok {
				util2.DeleteSecret(s.KubeClient, s.blackduck.Spec.Namespace, "blackduck-certificate")
			} else {
				return string(cert), string(key), nil
			}
		}
	}


	return  "", "", errors.New("default certificate could not be found")
}