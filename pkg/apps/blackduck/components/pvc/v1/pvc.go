package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"
)

type pvc struct {
	blackduck *v1.Blackduck
}

func (p *pvc) GetPVCs() []*components.PersistentVolumeClaim {
	var pvcs []*components.PersistentVolumeClaim

	defaultPVC := map[string]string{
		"blackduck-postgres":          "150Gi",
		"blackduck-authentication":    "2Gi",
		"blackduck-cfssl":             "2Gi",
		"blackduck-registration":      "2Gi",
		"blackduck-solr":              "2Gi",
		"blackduck-webapp":            "2Gi",
		"blackduck-logstash":          "20Gi",
		"blackduck-zookeeper-data":    "2Gi",
		"blackduck-zookeeper-datalog": "2Gi",
		"blackduck-rabbitmq":          "5Gi",
		"blackduck-uploadcache-data":  "100Gi",
		"blackduck-uploadcache-key":   "2Gi",
	}

	//if p.blackduck.Spec.ExternalPostgres != nil {
	//	delete(defaultPVC, "blackduck-postgres")
	//}
	//if !p.blackduck.Spec.isBinaryAnalysisEnabled {
	//	delete(defaultPVC, "blackduck-rabbitmq")
	//}

	if p.blackduck.Spec.PersistentStorage {
		for k, v := range defaultPVC {
			size := v
			storageClass := ""
			for _, claim := range p.blackduck.Spec.PVC {
				if strings.EqualFold(claim.Name, k) {
					if len(claim.StorageClass) > 0 {
						storageClass = claim.StorageClass
					}
					if len(claim.Size) > 0 {
						size = claim.StorageClass
					}
				}
				break
			}
			pvcs = append(pvcs, p.createPVC(k, size, v, storageClass, horizonapi.ReadWriteOnce, utils.GetLabel("pvc")))
		}
	}

	return pvcs
}

func (p *pvc) createPVC(name string, requestedSize string, defaultSize string, storageclass string, accessMode horizonapi.PVCAccessModeType, label map[string]string) *components.PersistentVolumeClaim {
	// Workaround so that storageClass does not get set to "", which prevent Kube from using the default storageClass
	var class *string
	if len(storageclass) > 0 {
		class = &storageclass
	} else if len(p.blackduck.Spec.PVCStorageClass) > 0 {
		class = &p.blackduck.Spec.PVCStorageClass
	} else {
		class = nil
	}

	var size string
	_, err := resource.ParseQuantity(requestedSize)
	if err != nil {
		size = defaultSize
	} else {
		size = requestedSize
	}

	pvc, _ := components.NewPersistentVolumeClaim(horizonapi.PVCConfig{
		Name:      name,
		Namespace: p.blackduck.Spec.Namespace,
		Size:      size,
		Class:     class,
	})

	pvc.AddAccessMode(accessMode)
	pvc.AddLabels(label)

	return pvc
}


func NewPvc(blackduck *v1.Blackduck) types.PVCInterface {
	return &pvc{blackduck: blackduck}
}