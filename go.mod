module github.com/blackducksoftware/synopsys-operator

go 1.12

require (
	github.com/Azure/go-autorest/autorest v0.0.0-00010101000000-000000000000 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.2.0 // indirect
	github.com/blackducksoftware/horizon v0.0.0-20190603173136-e141457f7a80
	github.com/docker/spdystream v0.0.0-00010101000000-000000000000 // indirect
	github.com/evanphx/json-patch v0.0.0-00010101000000-000000000000 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/google/go-cmp v0.3.0
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/gophercloud/gophercloud v0.0.0-00010101000000-000000000000 // indirect
	github.com/imdario/mergo v0.3.7
	github.com/jackc/pgx v3.3.0+incompatible // indirect
	github.com/juju/errors v0.0.0-20190207033735-e65537c515d7
	github.com/lib/pq v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/openshift/api v3.9.0+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/prometheus/client_golang v0.9.4
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	gopkg.in/inf.v0 v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/api v0.0.0-20190612125737-db0771252981
	k8s.io/apiextensions-apiserver v0.0.0-20190626210203-fdc73e13f9a6
	k8s.io/apimachinery v0.0.0-20190624085041-961b39a1baa0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/utils v0.0.0-00010101000000-000000000000 // indirect
	sigs.k8s.io/yaml v1.1.0
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.40.0
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.3.0
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.2.0
	github.com/blackducksoftware/horizon => github.com/blackducksoftware/horizon v0.0.0-20190603173136-e141457f7a80
	github.com/docker/spdystream => github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c
	github.com/elazarl/goproxy => github.com/elazarl/goproxy v0.0.0-20190711103511-473e67f1d7d2
	github.com/elazarl/goproxy/ext => github.com/elazarl/goproxy/ext v0.0.0-20190711103511-473e67f1d7d2
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.4.0
	github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.4.0
	github.com/google/go-cmp => github.com/google/go-cmp v0.3.0
	github.com/google/gofuzz => github.com/google/gofuzz v1.0.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.0
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.2.0
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.7
	github.com/jackc/pgx => github.com/jackc/pgx v3.3.0+incompatible
	github.com/juju/errors => github.com/juju/errors v0.0.0-20190207033735-e65537c515d7
	github.com/juju/loggo => github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8
	github.com/juju/testing => github.com/juju/testing v0.0.0-20190613124551-e81189438503
	github.com/lib/pq => github.com/lib/pq v1.1.1
	github.com/mitchellh/go-homedir => github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.5.0
	github.com/openshift/api => github.com/openshift/api v3.9.0+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v3.9.0+incompatible
	github.com/pkg/errors => github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra => github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.3
	github.com/spf13/viper => github.com/spf13/viper v1.4.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.3.0
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net => golang.org/x/net v0.0.0-20190611141213-3f473d35a33a
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190613124609-5ed2794edfdc
	google.golang.org/appengine => google.golang.org/appengine v1.6.1
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20190611190212-a7e196e89fd3
	google.golang.org/grpc => google.golang.org/grpc v1.21.1
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1
	gopkg.in/mgo.v2 => gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog => k8s.io/klog v0.3.3
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190709113604-33be087ad058
	k8s.io/utils => k8s.io/utils v0.0.0-20190607212802-c55fbcfc754a
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.1.0
)
