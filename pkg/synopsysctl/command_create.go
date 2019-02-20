// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synopsysctl

import (
	"encoding/json"
	"fmt"

	alertclientset "github.com/blackducksoftware/synopsys-operator/pkg/alert/client/clientset/versioned"
	alertv1 "github.com/blackducksoftware/synopsys-operator/pkg/api/alert/v1"
	blackduckv1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	opssightv1 "github.com/blackducksoftware/synopsys-operator/pkg/api/opssight/v1"
	blackduckclientset "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/client/clientset/versioned"
	opssightclientset "github.com/blackducksoftware/synopsys-operator/pkg/opssight/client/clientset/versioned"
	crddefaults "github.com/blackducksoftware/synopsys-operator/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Synopsys Resource in your cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		num_args := 1
		if len(args) != num_args {
			return fmt.Errorf("Must pass Namespace")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// createCmd represents the create command for Blackduck
var createBlackduckCmd = &cobra.Command{
	Use:   "blackduck",
	Short: "Create an instance of a Blackduck",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Must pass Namespace")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()

		// Create namespace for the Blackduck
		DeployCRDNamespace(restconfig)

		// Read Flags Into Default Blackduck Spec
		defaultBlackduckSpec := crddefaults.GetHubDefaultValue()
		flagset := cmd.Flags()
		flagset.VisitAll(checkBlackduckFlags)

		// Create and Deploy Blackduck CRD
		blackduck := &blackduckv1.Blackduck{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
			Spec: *defaultBlackduckSpec,
		}
		blackduckClient, err := blackduckclientset.NewForConfig(restconfig)
		_, err = blackduckClient.SynopsysV1().Blackducks(namespace).Create(blackduck)
		if err != nil {
			fmt.Printf("Error creating the Blackduck : %s\n", err)
			return
		}
	},
}

// createCmd represents the create command for OpsSight
var createOpsSightCmd = &cobra.Command{
	Use:   "opssight",
	Short: "Create an instance of OpsSight",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()

		// Create namespace for the OpsSight
		DeployCRDNamespace(restconfig)

		// Read Flags Into Default OpsSight Spec
		defaultOpsSightSpec := crddefaults.GetOpsSightDefaultValue()
		flagset := cmd.Flags()
		flagset.VisitAll(checkOpsSightFlags)

		// Create and Deploy OpsSight CRD
		opssight := &opssightv1.OpsSight{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
			Spec: *defaultOpsSightSpec,
		}
		opssightClient, err := opssightclientset.NewForConfig(restconfig)
		_, err = opssightClient.SynopsysV1().OpsSights(namespace).Create(opssight)
		if err != nil {
			fmt.Printf("Error creating the OpsSight : %s\n", err)
			return
		}
	},
}

// createCmd represents the create command for Alert
var createAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Create an instance of Alert",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()

		// Create namespace for the Alert
		DeployCRDNamespace(restconfig)

		// Read Flags Into Default Alert Spec
		defaultAlertSpec := crddefaults.GetAlertDefaultValue()
		flagset := cmd.Flags()
		flagset.VisitAll(checkAlertFlags)

		// Create and Deploy Alert CRD
		alert := &alertv1.Alert{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
			Spec: *defaultAlertSpec,
		}
		alertClient, err := alertclientset.NewForConfig(restconfig)
		_, err = alertClient.SynopsysV1().Alerts(namespace).Create(alert)
		if err != nil {
			fmt.Printf("Error creating the Alert : %s\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Add Blackduck Flags
	createBlackduckCmd.Flags().StringVar(&create_blackduck_size, "size", create_blackduck_size, "Blackduck size - small, medium, large")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_dbPrototype, "db-prototype", create_blackduck_dbPrototype, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_externalPostgres_postgresHost, "external-postgres-host", create_blackduck_externalPostgres_postgresHost, "TODO")
	createBlackduckCmd.Flags().IntVar(&create_blackduck_externalPostgres_postgresPort, "external-postgres-port", create_blackduck_externalPostgres_postgresPort, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_externalPostgres_postgresAdmin, "external-postgres-admin", create_blackduck_externalPostgres_postgresAdmin, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_externalPostgres_postgresUser, "external-postgres-user", create_blackduck_externalPostgres_postgresUser, "TODO")
	createBlackduckCmd.Flags().BoolVar(&create_blackduck_externalPostgres_postgresSsl, "external-postgres-ssl", create_blackduck_externalPostgres_postgresSsl, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_externalPostgres_postgresAdminPassword, "external-postgres-admin-password", create_blackduck_externalPostgres_postgresAdminPassword, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_externalPostgres_postgresUserPassword, "external-postgres-user-password", create_blackduck_externalPostgres_postgresUserPassword, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_pvcStorageClass, "pvc-storage-class", create_blackduck_pvcStorageClass, "TODO")
	createBlackduckCmd.Flags().BoolVar(&create_blackduck_livenessProbes, "liveness-probes", create_blackduck_livenessProbes, "Enable liveness probes")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_scanType, "scan-type", create_blackduck_scanType, "TODO")
	createBlackduckCmd.Flags().BoolVar(&create_blackduck_persistentStorage, "persistent-storage", create_blackduck_persistentStorage, "Enable persistent storage")
	createBlackduckCmd.Flags().StringSliceVar(&create_blackduck_PVC_json_slice, "pvc", create_blackduck_PVC_json_slice, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_certificateName, "db-certificate-name", create_blackduck_certificateName, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_certificate, "certificate", create_blackduck_certificate, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_certificateKey, "certificate-key", create_blackduck_certificateKey, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_proxyCertificate, "proxy-certificate", create_blackduck_proxyCertificate, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_type, "type", create_blackduck_type, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_desiredState, "desired-state", create_blackduck_desiredState, "TODO")
	createBlackduckCmd.Flags().StringSliceVar(&create_blackduck_environs, "environs", create_blackduck_environs, "TODO")
	createBlackduckCmd.Flags().StringSliceVar(&create_blackduck_imageRegistries, "image-registries", create_blackduck_imageRegistries, "List of image registries")
	createBlackduckCmd.Flags().StringSliceVar(&create_blackduck_imageUIDMap_json_slice, "image-uid-map", create_blackduck_imageUIDMap_json_slice, "TODO")
	createBlackduckCmd.Flags().StringVar(&create_blackduck_licenseKey, "license-key", create_blackduck_licenseKey, "TODO")
	createCmd.AddCommand(createBlackduckCmd)

	// Add OpsSight Flags
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceptor_name, "perceptor-name", create_opssight_perceptor_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceptor_image, "perceptor-image", create_opssight_perceptor_image, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_port, "perceptor-port", create_opssight_perceptor_port, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_checkForStalledScansPauseHours, "perceptor-check-scan-hours", create_opssight_perceptor_checkForStalledScansPauseHours, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_stalledScanClientTimeoutHours, "perceptor-scan-client-timeout-hours", create_opssight_perceptor_stalledScanClientTimeoutHours, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_modelMetricsPauseSeconds, "perceptor-metrics-pause-seconds", create_opssight_perceptor_modelMetricsPauseSeconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_unknownImagePauseMilliseconds, "perceptor-unknown-image-pause-milliseconds", create_opssight_perceptor_unknownImagePauseMilliseconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceptor_clientTimeoutMilliseconds, "perceptor-client-timeout-milliseconds", create_opssight_perceptor_clientTimeoutMilliseconds, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_name, "scannerpod-name", create_opssight_scannerPod_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_scanner_name, "scannerpod-scanner-name", create_opssight_scannerPod_scanner_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_scanner_image, "scannerpod-scanner-image", create_opssight_scannerPod_scanner_image, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_scannerPod_scanner_port, "scannerpod-scanner-port", create_opssight_scannerPod_scanner_port, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_scannerPod_scanner_clientTimeoutSeconds, "scannerpod-scanner-client-timeout-seconds", create_opssight_scannerPod_scanner_clientTimeoutSeconds, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_imageFacade_name, "scannerpod-imagefacade-name", create_opssight_scannerPod_imageFacade_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_imageFacade_image, "scannerpod-imagefacade-image", create_opssight_scannerPod_imageFacade_image, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_scannerPod_imageFacade_port, "scannerpod-imagefacade-port", create_opssight_scannerPod_imageFacade_port, "TODO")
	createOpsSightCmd.Flags().StringSliceVar(&create_opssight_scannerPod_imageFacade_internalRegistries_json_slice, "scannerpod-imagefacade-internal-registries", create_opssight_scannerPod_imageFacade_internalRegistries_json_slice, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_imageFacade_imagePullerType, "scannerpod-imagefacade-image-puller-type", create_opssight_scannerPod_imageFacade_imagePullerType, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_imageFacade_serviceAccount, "scannerpod-imagefacade-service-account", create_opssight_scannerPod_imageFacade_serviceAccount, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_scannerPod_replicaCount, "scannerpod-replica-count", create_opssight_scannerPod_replicaCount, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_scannerPod_imageDirectory, "scannerpod-image-directory", create_opssight_scannerPod_imageDirectory, "TODO")
	createOpsSightCmd.Flags().BoolVar(&create_opssight_perceiver_enableImagePerceiver, "enable-image-perceiver", create_opssight_perceiver_enableImagePerceiver, "TODO")
	createOpsSightCmd.Flags().BoolVar(&create_opssight_perceiver_enablePodPerceiver, "enable-pod-perceiver", create_opssight_perceiver_enablePodPerceiver, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_imagePerceiver_name, "imageperceiver-name", create_opssight_perceiver_imagePerceiver_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_imagePerceiver_image, "imageperceiver-image", create_opssight_perceiver_imagePerceiver_image, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_podPerceiver_name, "podperceiver-name", create_opssight_perceiver_podPerceiver_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_podPerceiver_image, "podperceiver-image", create_opssight_perceiver_podPerceiver_image, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_podPerceiver_namespaceFilter, "podperceiver-namespace-filter", create_opssight_perceiver_podPerceiver_namespaceFilter, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceiver_annotationIntervalSeconds, "perceiver-annotation-interval-seconds", create_opssight_perceiver_annotationIntervalSeconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceiver_dumpIntervalMinutes, "perceiver-dump-interval-minutes", create_opssight_perceiver_dumpIntervalMinutes, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_perceiver_serviceAccount, "perceiver-service-account", create_opssight_perceiver_serviceAccount, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_perceiver_port, "perceiver-port", create_opssight_perceiver_port, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_prometheus_name, "prometheus-name", create_opssight_prometheus_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_prometheus_name, "prometheus-image", create_opssight_prometheus_name, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_prometheus_port, "prometheus-port", create_opssight_prometheus_port, "TODO")
	createOpsSightCmd.Flags().BoolVar(&create_opssight_enableSkyfire, "enable-skyfire", create_opssight_enableSkyfire, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_skyfire_name, "skyfire-name", create_opssight_skyfire_name, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_skyfire_image, "skyfire-image", create_opssight_skyfire_image, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_port, "skyfire-port", create_opssight_skyfire_port, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_prometheusPort, "skyfire-prometheus-port", create_opssight_skyfire_prometheusPort, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_skyfire_serviceAccount, "skyfire-service-account", create_opssight_skyfire_serviceAccount, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_hubClientTimeoutSeconds, "skyfire-hub-client-timeout-seconds", create_opssight_skyfire_hubClientTimeoutSeconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_hubDumpPauseSeconds, "skyfire-hub-dump-pause-seconds", create_opssight_skyfire_hubDumpPauseSeconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_kubeDumpIntervalSeconds, "skyfire-kube-dump-interval-seconds", create_opssight_skyfire_kubeDumpIntervalSeconds, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_skyfire_perceptorDumpIntervalSeconds, "skyfire-perceptor-dump-interval-seconds", create_opssight_skyfire_perceptorDumpIntervalSeconds, "TODO")
	createOpsSightCmd.Flags().StringSliceVar(&create_opssight_blackduck_hosts, "blackduck-hosts", create_opssight_blackduck_hosts, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_blackduck_user, "blackduck-user", create_opssight_blackduck_user, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_port, "blackduck-port", create_opssight_blackduck_port, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_concurrentScanLimit, "blackduck-concurrent-scan-limit", create_opssight_blackduck_concurrentScanLimit, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_totalScanLimit, "blackduck-total-scan-limit", create_opssight_blackduck_totalScanLimit, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_blackduck_passwordEnvVar, "blackduck-password-environment-variable", create_opssight_blackduck_passwordEnvVar, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_initialCount, "blackduck-initial-count", create_opssight_blackduck_initialCount, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_maxCount, "blackduck-max-count", create_opssight_blackduck_maxCount, "TODO")
	createOpsSightCmd.Flags().IntVar(&create_opssight_blackduck_deleteHubThresholdPercentage, "blackduck-delete-blackduck-threshold-percentage", create_opssight_blackduck_deleteHubThresholdPercentage, "TODO")
	createOpsSightCmd.Flags().BoolVar(&create_opssight_enableMetrics, "enable-metrics", create_opssight_enableMetrics, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_defaultCPU, "default-cpu", create_opssight_defaultCPU, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_defaultMem, "default-mem", create_opssight_defaultMem, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_logLevel, "log-level", create_opssight_logLevel, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_configMapName, "config-map-name", create_opssight_configMapName, "TODO")
	createOpsSightCmd.Flags().StringVar(&create_opssight_secretName, "secret-name", create_opssight_secretName, "TODO")
	createCmd.AddCommand(createOpsSightCmd)

	// Add Alert Flags
	createAlertCmd.Flags().StringVar(&create_alert_registry, "alert-registry", create_alert_registry, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_imagePath, "image-path", create_alert_imagePath, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_alertImageName, "alert-image-name", create_alert_alertImageName, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_alertImageVersion, "alert-image-version", create_alert_alertImageVersion, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_cfsslImageName, "cfssl-image-name", create_alert_cfsslImageName, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_cfsslImageVersion, "cfssl-image-version", create_alert_cfsslImageVersion, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_blackduckHost, "blackduck-host", create_alert_blackduckHost, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_blackduckUser, "blackduck-user", create_alert_blackduckUser, "TODO")
	createAlertCmd.Flags().IntVar(&create_alert_blackduckPort, "blackduck-port", create_alert_blackduckPort, "TODO")
	createAlertCmd.Flags().IntVar(&create_alert_port, "port", create_alert_port, "TODO")
	createAlertCmd.Flags().BoolVar(&create_alert_standAlone, "stand-alone", create_alert_standAlone, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_alertMemory, "alert-memory", create_alert_alertMemory, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_cfsslMemory, "cfssl-memory", create_alert_cfsslMemory, "TODO")
	createAlertCmd.Flags().StringVar(&create_alert_state, "alert-state", create_alert_state, "TODO")
	createCmd.AddCommand(createAlertCmd)
}

func checkBlackduckFlags(f *pflag.Flag) {
	if f.Changed {
		fmt.Printf("Flag %s: CHANGED\n", f.Name)
		switch f.Name {
		case "namespace":
			defaultBlackduckSpec.Namespace = namespace
		case "size":
			defaultBlackduckSpec.Size = create_blackduck_size
		case "db-prototype":
			defaultBlackduckSpec.DbPrototype = create_blackduck_dbPrototype
		case "external-postgres-host":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresHost = create_blackduck_externalPostgres_postgresHost
		case "external-postgres-port":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresPort = create_blackduck_externalPostgres_postgresPort
		case "external-postgres-admin":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresAdmin = create_blackduck_externalPostgres_postgresAdmin
		case "external-postgres-user":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresUser = create_blackduck_externalPostgres_postgresUser
		case "external-postgres-ssl":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresSsl = create_blackduck_externalPostgres_postgresSsl
		case "external-postgres-admin-password":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresAdminPassword = create_blackduck_externalPostgres_postgresAdminPassword
		case "external-postgres-user-password":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.ExternalPostgres.PostgresUserPassword = create_blackduck_externalPostgres_postgresUserPassword
		case "pvc-storage-class":
			if defaultBlackduckSpec.ExternalPostgres == nil {
				defaultBlackduckSpec.ExternalPostgres = &blackduckv1.PostgresExternalDBConfig{}
			}
			defaultBlackduckSpec.PVCStorageClass = create_blackduck_pvcStorageClass
		case "liveness-probes":
			defaultBlackduckSpec.LivenessProbes = create_blackduck_livenessProbes
		case "scan-type":
			defaultBlackduckSpec.ScanType = create_blackduck_scanType
		case "persistent-storage":
			defaultBlackduckSpec.PersistentStorage = create_blackduck_persistentStorage
		case "pvc":
			for _, pvc_json := range create_blackduck_PVC_json_slice {
				pvc := &blackduckv1.PVC{}
				json.Unmarshal([]byte(pvc_json), pvc)
				defaultBlackduckSpec.PVC = append(defaultBlackduckSpec.PVC, *pvc)
			}
		case "db-certificate-name":
			defaultBlackduckSpec.CertificateName = create_blackduck_certificateName
		case "certificate":
			defaultBlackduckSpec.Certificate = create_blackduck_certificate
		case "certificate-key":
			defaultBlackduckSpec.CertificateKey = create_blackduck_certificateKey
		case "proxy-certificate":
			defaultBlackduckSpec.ProxyCertificate = create_blackduck_proxyCertificate
		case "type":
			defaultBlackduckSpec.Type = create_blackduck_type
		case "desired-state":
			defaultBlackduckSpec.DesiredState = create_blackduck_desiredState
		case "environs":
			defaultBlackduckSpec.Environs = create_blackduck_environs
		case "image-registries":
			defaultBlackduckSpec.ImageRegistries = create_blackduck_imageRegistries
		case "image-uid-map":
			type uid struct {
				Key   string `json:"key"`
				Value int64  `json:"value"`
			}
			defaultBlackduckSpec.ImageUIDMap = make(map[string]int64)
			for _, uid_json := range create_blackduck_imageUIDMap_json_slice {
				uid_struct := &uid{}
				json.Unmarshal([]byte(uid_json), uid_struct)
				defaultBlackduckSpec.ImageUIDMap[uid_struct.Key] = uid_struct.Value
			}
		case "license-key":
			defaultBlackduckSpec.LicenseKey = create_blackduck_licenseKey
		default:
			fmt.Printf("Flag %s: Not Found\n", f.Name)
		}
	}
	fmt.Printf("Flag %s: UNCHANGED\n", f.Name)
}

func checkOpsSightFlags(f *pflag.Flag) {
	if f.Changed {
		fmt.Printf("Flag %s: CHANGED\n", f.Name)
		switch f.Name {
		case "perceptor-name":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.Name = create_opssight_perceptor_name
		case "perceptor-image":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.Image = create_opssight_perceptor_image
		case "perceptor-port":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.Port = create_opssight_perceptor_port
		case "perceptor-check-scan-hours":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.CheckForStalledScansPauseHours = create_opssight_perceptor_checkForStalledScansPauseHours
		case "perceptor-scan-client-timeout-hours":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.StalledScanClientTimeoutHours = create_opssight_perceptor_stalledScanClientTimeoutHours
		case "perceptor-metrics-pause-seconds":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.ModelMetricsPauseSeconds = create_opssight_perceptor_modelMetricsPauseSeconds
		case "perceptor-unknown-image-pause-milliseconds":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.UnknownImagePauseMilliseconds = create_opssight_perceptor_unknownImagePauseMilliseconds
		case "perceptor-client-timeout-milliseconds":
			if defaultOpsSightSpec.Perceptor == nil {
				defaultOpsSightSpec.Perceptor = &opssightv1.Perceptor{}
			}
			defaultOpsSightSpec.Perceptor.ClientTimeoutMilliseconds = create_opssight_perceptor_clientTimeoutMilliseconds
		case "scannerpod-name":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			defaultOpsSightSpec.ScannerPod.Name = create_opssight_scannerPod_name
		case "scannerpod-scanner-name":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.Scanner == nil {
				defaultOpsSightSpec.ScannerPod.Scanner = &opssightv1.Scanner{}
			}
			defaultOpsSightSpec.ScannerPod.Scanner.Name = create_opssight_scannerPod_scanner_name
		case "scannerpod-scanner-image":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.Scanner == nil {
				defaultOpsSightSpec.ScannerPod.Scanner = &opssightv1.Scanner{}
			}
			defaultOpsSightSpec.ScannerPod.Scanner.Image = create_opssight_scannerPod_scanner_image
		case "scannerpod-scanner-port":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.Scanner == nil {
				defaultOpsSightSpec.ScannerPod.Scanner = &opssightv1.Scanner{}
			}
			defaultOpsSightSpec.ScannerPod.Scanner.Port = create_opssight_scannerPod_scanner_port
		case "scannerpod-scanner-client-timeout-seconds":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.Scanner == nil {
				defaultOpsSightSpec.ScannerPod.Scanner = &opssightv1.Scanner{}
			}
			defaultOpsSightSpec.ScannerPod.Scanner.ClientTimeoutSeconds = create_opssight_scannerPod_scanner_clientTimeoutSeconds
		case "scannerpod-imagefacade-name":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			defaultOpsSightSpec.ScannerPod.ImageFacade.Name = create_opssight_scannerPod_imageFacade_name
		case "scannerpod-imagefacade-image":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			defaultOpsSightSpec.ScannerPod.ImageFacade.Image = create_opssight_scannerPod_imageFacade_image
		case "scannerpod-imagefacade-port":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			defaultOpsSightSpec.ScannerPod.ImageFacade.Port = create_opssight_scannerPod_imageFacade_port
		case "scannerpod-imagefacade-internal-registries":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			for _, registry_json := range create_opssight_scannerPod_imageFacade_internalRegistries_json_slice {
				registry := &opssightv1.RegistryAuth{}
				json.Unmarshal([]byte(registry_json), registry)
				defaultOpsSightSpec.ScannerPod.ImageFacade.InternalRegistries = append(defaultOpsSightSpec.ScannerPod.ImageFacade.InternalRegistries, *registry)
			}
		case "scannerpod-imagefacade-image-puller-type":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			defaultOpsSightSpec.ScannerPod.ImageFacade.ImagePullerType = create_opssight_scannerPod_imageFacade_imagePullerType
		case "scannerpod-imagefacade-service-account":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			if defaultOpsSightSpec.ScannerPod.ImageFacade == nil {
				defaultOpsSightSpec.ScannerPod.ImageFacade = &opssightv1.ImageFacade{}
			}
			defaultOpsSightSpec.ScannerPod.ImageFacade.ServiceAccount = create_opssight_scannerPod_imageFacade_serviceAccount
		case "scannerpod-replica-count":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			defaultOpsSightSpec.ScannerPod.ReplicaCount = create_opssight_scannerPod_replicaCount
		case "scannerpod-image-directory":
			if defaultOpsSightSpec.ScannerPod == nil {
				defaultOpsSightSpec.ScannerPod = &opssightv1.ScannerPod{}
			}
			defaultOpsSightSpec.ScannerPod.ImageDirectory = create_opssight_scannerPod_imageDirectory
		case "enable-image-perceiver":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.EnableImagePerceiver = create_opssight_perceiver_enableImagePerceiver
		case "enable-pod-perceiver":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.EnablePodPerceiver = create_opssight_perceiver_enablePodPerceiver
		case "imageperceiver-name":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			if defaultOpsSightSpec.Perceiver.ImagePerceiver == nil {
				defaultOpsSightSpec.Perceiver.ImagePerceiver = &opssightv1.ImagePerceiver{}
			}
			defaultOpsSightSpec.Perceiver.ImagePerceiver.Name = create_opssight_perceiver_imagePerceiver_name
		case "imageperceiver-image":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			if defaultOpsSightSpec.Perceiver.ImagePerceiver == nil {
				defaultOpsSightSpec.Perceiver.ImagePerceiver = &opssightv1.ImagePerceiver{}
			}
			defaultOpsSightSpec.Perceiver.ImagePerceiver.Image = create_opssight_perceiver_imagePerceiver_image
		case "podperceiver-name":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			if defaultOpsSightSpec.Perceiver.PodPerceiver == nil {
				defaultOpsSightSpec.Perceiver.PodPerceiver = &opssightv1.PodPerceiver{}
			}
			defaultOpsSightSpec.Perceiver.PodPerceiver.Name = create_opssight_perceiver_podPerceiver_name
		case "podperceiver-image":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			if defaultOpsSightSpec.Perceiver.PodPerceiver == nil {
				defaultOpsSightSpec.Perceiver.PodPerceiver = &opssightv1.PodPerceiver{}
			}
			defaultOpsSightSpec.Perceiver.PodPerceiver.Image = create_opssight_perceiver_podPerceiver_image
		case "podperceiver-namespace-filter":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			if defaultOpsSightSpec.Perceiver.PodPerceiver == nil {
				defaultOpsSightSpec.Perceiver.PodPerceiver = &opssightv1.PodPerceiver{}
			}
			defaultOpsSightSpec.Perceiver.PodPerceiver.NamespaceFilter = create_opssight_perceiver_podPerceiver_namespaceFilter
		case "perceiver-annotation-interval-seconds":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.AnnotationIntervalSeconds = create_opssight_perceiver_annotationIntervalSeconds
		case "perceiver-dump-interval-minutes":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.DumpIntervalMinutes = create_opssight_perceiver_dumpIntervalMinutes
		case "perceiver-service-account":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.ServiceAccount = create_opssight_perceiver_serviceAccount
		case "perceiver-port":
			if defaultOpsSightSpec.Perceiver == nil {
				defaultOpsSightSpec.Perceiver = &opssightv1.Perceiver{}
			}
			defaultOpsSightSpec.Perceiver.Port = create_opssight_perceiver_port
		case "prometheus-name":
			if defaultOpsSightSpec.Prometheus == nil {
				defaultOpsSightSpec.Prometheus = &opssightv1.Prometheus{}
			}
			defaultOpsSightSpec.Prometheus.Name = create_opssight_prometheus_name
		case "prometheus-image":
			if defaultOpsSightSpec.Prometheus == nil {
				defaultOpsSightSpec.Prometheus = &opssightv1.Prometheus{}
			}
			defaultOpsSightSpec.Prometheus.Image = create_opssight_prometheus_image
		case "prometheus-port":
			if defaultOpsSightSpec.Prometheus == nil {
				defaultOpsSightSpec.Prometheus = &opssightv1.Prometheus{}
			}
			defaultOpsSightSpec.Prometheus.Port = create_opssight_prometheus_port
		case "enable-skyfire":
			defaultOpsSightSpec.EnableSkyfire = create_opssight_enableSkyfire
		case "skyfire-name":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.Name = create_opssight_skyfire_name
		case "skyfire-image":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.Image = create_opssight_skyfire_image
		case "skyfire-port":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.Port = create_opssight_skyfire_port
		case "skyfire-prometheus-port":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.PrometheusPort = create_opssight_skyfire_prometheusPort
		case "skyfire-service-account":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.ServiceAccount = create_opssight_skyfire_serviceAccount
		case "skyfire-hub-client-timeout-seconds":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.HubClientTimeoutSeconds = create_opssight_skyfire_hubClientTimeoutSeconds
		case "skyfire-hub-dump-pause-seconds":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.HubDumpPauseSeconds = create_opssight_skyfire_hubDumpPauseSeconds
		case "skyfire-kube-dump-interval-seconds":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.KubeDumpIntervalSeconds = create_opssight_skyfire_kubeDumpIntervalSeconds
		case "skyfire-perceptor-dump-interval-seconds":
			if defaultOpsSightSpec.Skyfire == nil {
				defaultOpsSightSpec.Skyfire = &opssightv1.Skyfire{}
			}
			defaultOpsSightSpec.Skyfire.PerceptorDumpIntervalSeconds = create_opssight_skyfire_perceptorDumpIntervalSeconds
		case "blackduck-hosts":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.Hosts = create_opssight_blackduck_hosts
		case "blackduck-user":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.User = create_opssight_blackduck_user
		case "blackduck-port":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.Port = create_opssight_blackduck_port
		case "blackduck-concurrent-scan-limit":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.ConcurrentScanLimit = create_opssight_blackduck_concurrentScanLimit
		case "blackduck-total-scan-limit":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.TotalScanLimit = create_opssight_blackduck_totalScanLimit
		case "blackduck-password-environment-variable":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.PasswordEnvVar = create_opssight_blackduck_passwordEnvVar
		case "blackduck-initial-count":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.InitialCount = create_opssight_blackduck_initialCount
		case "blackduck-max-count":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.MaxCount = create_opssight_blackduck_maxCount
		case "blackduck-delete-blackduck-threshold-percentage":
			if defaultOpsSightSpec.Blackduck == nil {
				defaultOpsSightSpec.Blackduck = &opssightv1.Blackduck{}
			}
			defaultOpsSightSpec.Blackduck.DeleteHubThresholdPercentage = create_opssight_blackduck_deleteHubThresholdPercentage
		case "enable-metrics":
			defaultOpsSightSpec.EnableMetrics = create_opssight_enableMetrics
		case "default-cpu":
			defaultOpsSightSpec.DefaultCPU = create_opssight_defaultCPU
		case "default-mem":
			defaultOpsSightSpec.DefaultMem = create_opssight_defaultMem
		case "log-level":
			defaultOpsSightSpec.LogLevel = create_opssight_logLevel
		case "config-map-name":
			defaultOpsSightSpec.ConfigMapName = create_opssight_configMapName
		case "secret-name":
			defaultOpsSightSpec.SecretName = create_opssight_secretName
		default:
			fmt.Printf("Flag %s: Not Found\n", f.Name)
		}
	}
	fmt.Printf("Flag %s: UNCHANGED\n", f.Name)

}

func checkAlertFlags(f *pflag.Flag) {
	if f.Changed {
		fmt.Printf("Flag %s: CHANGED\n", f.Name)
		switch f.Name {
		case "namespace":
			defaultAlertSpec.Namespace = namespace
		case "alert-registry":
			defaultAlertSpec.Registry = create_alert_registry
		case "image-path":
			defaultAlertSpec.ImagePath = create_alert_imagePath
		case "alert-image-name":
			defaultAlertSpec.AlertImageName = create_alert_alertImageName
		case "alert-image-version":
			defaultAlertSpec.AlertImageVersion = create_alert_alertImageVersion
		case "cfssl-image-name":
			defaultAlertSpec.CfsslImageName = create_alert_cfsslImageName
		case "cfssl-image-version":
			defaultAlertSpec.CfsslImageVersion = create_alert_cfsslImageVersion
		case "blackduck-host":
			defaultAlertSpec.BlackduckHost = create_alert_blackduckHost
		case "blackduck-user":
			defaultAlertSpec.BlackduckUser = create_alert_blackduckUser
		case "blackduck-port":
			defaultAlertSpec.BlackduckPort = &create_alert_blackduckPort
		case "port":
			defaultAlertSpec.Port = &create_alert_port
		case "stand-alone":
			defaultAlertSpec.StandAlone = &create_alert_standAlone
		case "alert-memory":
			defaultAlertSpec.AlertMemory = create_alert_alertMemory
		case "cfssl-memory":
			defaultAlertSpec.CfsslMemory = create_alert_cfsslMemory
		default:
			fmt.Printf("Flag %s: Not Found\n", f.Name)
		}
	}
	fmt.Printf("Flag %s: UNCHANGED\n", f.Name)
}