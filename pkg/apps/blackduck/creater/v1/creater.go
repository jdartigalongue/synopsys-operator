package v1

import (
	"fmt"
	"github.com/blackducksoftware/synopsys-operator/pkg/api"
	clusterRoleBidingv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/clusterrolebinding/v1"
	configmapv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/configmaps/v1"
	pvcv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/pvc/v1"
	postgresv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/postgres/v1"
	secretv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/secrets/v1"
	serviceAccountv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/serviceaccount/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/crdupdater"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"reflect"
	"strings"

	// Rc
	authenticationv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/authentication/v1"
	cfsslv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/cfssl/v1"
	documentationv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/documentation/v1"
	jobrunnerv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/jobrunner/v1"
	registrationv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/registration/v1"
	scanv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/scan/v1"
	solrv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/solr/v1"
	uploadcachev1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/uploadcache/v1"
	webapplogstashv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/webapplogstash/v1"
	webserver1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/webserver/v1"
	zookeeperv1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/zookeeper/v1"

	// Services
	authenticationSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/authentication/v1"
	cfsslSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/cfssl/v1"
	documentationSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/documentation/v1"
	logstashSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/logstash/v1"
	postgresSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/postgres/v1"
	rabbitmqSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/rabbitmq/v1"
	registrationSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/regisrtration/v1"
	scanSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/scan/v1"
	solrSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/solr/v1"
	uploadcacheSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/uploadcache/v1"
	webappSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/webapp/v1"
	webserverSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/webserver/v1"
	zookeeperSvcV1 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services/zookeeper/v1"

	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/sizes/v1"
	c "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/creater"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	blackduckclientset "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/client/clientset/versioned"
)

type Creater struct {
	c.Creater
}

var imageTags = map[string]c.Blackduck{
	"2019.4.1": {
		Rc: map[string]c.ComponentReplicationController{
			"postgres": {
				Func: postgresv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-postgres": {
						Image:"registry.access.redhat.com/rhscl/postgresql-96-rhel7:1",
						Tag:"1",
					},
				},
			},
			"authentication": {
				Func: authenticationv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-authentication": {Tag:"2019.4.1"},
				},
			},
			"documentation": {
				Func: documentationv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-documentation":  {Tag:"2019.4.1"},
				},
			},
			"registration": {
				Func: registrationv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-registration":  {Tag:"2019.4.1"},
				},
			},
			"scan": {
				Func: scanv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-scan":  {Tag:"2019.4.1"},
				},
			},
			"webapp-logstash": {
				Func: webapplogstashv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-webapp":    {Tag:"2019.4.1"},
					"blackduck-logstash":  {Tag:"1.0.4"},
				},
			},
			"cfssl": {
				Func: cfsslv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-cfssl":  {Tag:"1.0.0"},
				},
			},
			"webserver": {
				Func: webserver1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-nginx":  {Tag:"1.0.7"},
				},
			},
			"solr": {
				Func: solrv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-solr":  {Tag:"1.0.0"},
				},
			},
			"zookeeper": {
				Func: zookeeperv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-zookeeper":  {Tag:"1.0.0"},
				},
			},
			"upload-cache": {
				Func: uploadcachev1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-uploadcache":  {Tag:"1.0.8"},
				},
			},
			"jobrunner": {
				Func: jobrunnerv1.NewDeploymentVersion,
				Tags: map[string]c.TagOrImage{
					"blackduck-jobrunner":  {Tag:"2019.4.1"},
				},
			},
		},
		Service: map[string]func(blackduck *blackduckapi.Blackduck) services.ServiceInterface{
			"authentication": authenticationSvcV1.NewService,
			"cfssl":          cfsslSvcV1.NewService,
			"documentation":  documentationSvcV1.NewService,
			"logstash":       logstashSvcV1.NewService,
			"postgres":       postgresSvcV1.NewService,
			"rabbitmq":       rabbitmqSvcV1.NewService,
			"registration":   registrationSvcV1.NewService,
			"scan":           scanSvcV1.NewService,
			"solr":           solrSvcV1.NewService,
			"uploadcache":    uploadcacheSvcV1.NewService,
			"webapp":         webappSvcV1.NewService,
			"webserver":      webserverSvcV1.NewService,
			"zookeeper":      zookeeperSvcV1.NewService,
		},
		Configmap:          configmapv1.NewConfigmap,
		PVC:                pvcv1.NewPvc,
		Secret:             secretv1.NewSecret,
		Size:               v1.NewSize(),
		ServiceAccount:     serviceAccountv1.NewServiceaccount,
		ClusterRoleBinding: clusterRoleBidingv1.NewClusterrolebinding,
	},
}

//"2019.4.1": {
//	"blackduck-authentication": "2019.4.1",
//	"blackduck-documentation":  "2019.4.1",
//	"blackduck-jobrunner":      "2019.4.1",
//	"blackduck-registration":   "2019.4.1",
//	"blackduck-scan":           "2019.4.1",
//	"blackduck-webapp":         "2019.4.1",
//	"blackduck-cfssl":          "1.0.0",
//	"blackduck-logstash":       "1.0.4",
//	"blackduck-nginx":          "1.0.7",
//	"blackduck-solr":           "1.0.0",
//	"blackduck-zookeeper":      "1.0.0",
//	"blackduck-upload-cache":   "1.0.8",
//	"appcheck-worker":          "2019.01",
//	"rabbitmq":                 "1.0.0",
//},

func NewCreater(config *protoform.Config, kubeConfig *rest.Config, kubeClient *kubernetes.Clientset, bdClient *blackduckclientset.Clientset) *Creater {
	return &Creater{c.Creater{Config: config, KubeConfig: kubeConfig, KubeClient: kubeClient, BlackduckClient:bdClient}}
}

func (hc *Creater) Ensure(blackduck *blackduckapi.Blackduck) error {
	newBlackuck := blackduck.DeepCopy()
	components := hc.GetComponents(blackduck, imageTags[blackduck.Spec.Version])


	if strings.EqualFold(blackduck.Spec.DesiredState, "STOP") {
		commonConfig := crdupdater.NewCRUDComponents(hc.KubeConfig, hc.KubeClient, hc.Config.DryRun, false, blackduck.Spec.Namespace,
			&api.ComponentList{PersistentVolumeClaims: components.PersistentVolumeClaims}, "app=blackduck")
		_, errors := commonConfig.CRUDComponents()
		if len(errors) > 0 {
			return fmt.Errorf("stop blackduck: %+v", errors)
		}
	} else {
		commonConfig := crdupdater.NewCRUDComponents(hc.KubeConfig, hc.KubeClient, hc.Config.DryRun, false, blackduck.Spec.Namespace,
			&api.ComponentList{PersistentVolumeClaims: components.PersistentVolumeClaims}, "app=blackduck,component=pvc")
		isPatched, errors := commonConfig.CRUDComponents()
		if len(errors) > 0 {
			return fmt.Errorf("update pvc: %+v", errors)
		}


		// install postgres
		commonConfig = crdupdater.NewCRUDComponents(hc.KubeConfig, hc.KubeClient, hc.Config.DryRun, isPatched, blackduck.Spec.Namespace,
			components, "app=blackduck,component=postgres")
		isPatched, errors = commonConfig.CRUDComponents()
		if len(errors) > 0 {
			return fmt.Errorf("update postgres component: %+v", errors)
		}
		// log.Debugf("created/updated postgres component for %s", blackduck.Spec.Namespace)

		// Check postgres and initialize if needed.
		if blackduck.Spec.ExternalPostgres == nil {
			//// TODO return whether we re-initialized or not
			//err = hc.initPostgres(&blackduck.Spec)
			//if err != nil {
			//	return err
			//}
		}

		// install cfssl
		commonConfig = crdupdater.NewCRUDComponents(hc.KubeConfig, hc.KubeClient, hc.Config.DryRun, isPatched, blackduck.Spec.Namespace,
			components, "app=blackduck,component in (configmap,serviceAccount,cfssl)")
		isPatched, errors = commonConfig.CRUDComponents()
		if len(errors) > 0 {
			return fmt.Errorf("update cfssl component: %+v", errors)
		}

		if err := util.ValidatePodsAreRunningInNamespace(hc.KubeClient, blackduck.Spec.Namespace, hc.Config.PodWaitTimeoutSeconds); err != nil {
			return err
		}

		// deploy non postgres and uploadcache component
		commonConfig = crdupdater.NewCRUDComponents(hc.KubeConfig, hc.KubeClient, hc.Config.DryRun, isPatched, blackduck.Spec.Namespace,
			components, "app=blackduck,component notin (postgres,configmap,serviceAccount,cfssl)")
		isPatched, errors = commonConfig.CRUDComponents()
		if len(errors) > 0 {
			return fmt.Errorf("update non postgres and cfssl component: %+v", errors)
		}
		// log.Debugf("created/updated non postgres and upload cache component for %s", blackduck.Spec.Namespace)

		//// add security context constraint if bdba enabled
		//if hc.isBinaryAnalysisEnabled(&blackduck.Spec) {
		//	// log.Debugf("created/updated upload cache component for %s", blackduck.Spec.Namespace)
		//	err = hc.addAnyUIDToServiceAccount(&blackduck.Spec)
		//	if err != nil {
		//		log.Error(err)
		//	}
		//}
		//
		//// add security context constraint if bdba enabled
		//if hc.isBinaryAnalysisEnabled(&blackduck.Spec) {
		//	err = hc.addAnyUIDToServiceAccount(&blackduck.Spec)
		//	if err != nil {
		//		log.Error(err)
		//	}
		//}

		//if strings.ToUpper(blackduck.Spec.ExposeService) == "NODEPORT" {
		//	newBlackuck.Status.IP, err = bdutils.GetNodePortIPAddress(hc.KubeClient, blackduck.Spec.Namespace, "webserver-exposed")
		//} else if strings.ToUpper(blackduck.Spec.ExposeService) == "LOADBALANCER" {
		//	newBlackuck.Status.IP, err = bdutils.GetLoadBalancerIPAddress(hc.KubeClient, blackduck.Spec.Namespace, "webserver-exposed")
		//}
		//
		//// Create Route on Openshift
		//if strings.ToUpper(blackduck.Spec.ExposeService) == "OPENSHIFT" && hc.routeClient != nil {
		//	route, err := util.GetOpenShiftRoutes(hc.routeClient, blackduck.Spec.Namespace, blackduck.Spec.Namespace)
		//	if err != nil {
		//		route, err = util.CreateOpenShiftRoutes(hc.routeClient, blackduck.Spec.Namespace, blackduck.Spec.Namespace, "Service", "webserver", "port-webserver", routev1.TLSTerminationPassthrough)
		//		if err != nil {
		//			log.Errorf("unable to create the openshift route due to %+v", err)
		//		}
		//	}
		//	if route != nil {
		//		newBlackuck.Status.IP = route.Spec.Host
		//	}
		//}
		//
		//if err := util.ValidatePodsAreRunningInNamespace(hc.KubeClient, blackduck.Spec.Namespace, 600); err != nil {
		//	return err
		//}

		// TODO wait for webserver to be up before we register
		//if len(blackduck.Spec.LicenseKey) > 0 {
		//	if err := hc.registerIfNeeded(blackduck); err != nil {
		//		log.Infof("couldn't register blackduck %s: %v", blackduck.Name, err)
		//	}
		//}
	}

	if blackduck.Spec.PersistentStorage {
		pvcVolumeNames := map[string]string{}
		for _, v := range blackduck.Spec.PVC {
			pvName, err := hc.GetPVCVolumeName(blackduck.Spec.Namespace, v.Name)
			if err != nil {
				continue
			}
			pvcVolumeNames[v.Name] = pvName
		}
		newBlackuck.Status.PVCVolumeName = pvcVolumeNames
	}

	if !reflect.DeepEqual(blackduck, newBlackuck) {
		if _, err := hc.BlackduckClient.SynopsysV1().Blackducks(hc.Config.Namespace).Update(newBlackuck); err != nil {
			return err
		}
	}

	return nil
}

func (hc *Creater) Versions() []string {
	var versions []string
	for k := range imageTags {
		versions = append(versions, k)
	}
	return versions
}
