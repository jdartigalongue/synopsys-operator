package store

import (
	"fmt"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"reflect"

	//_ "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"log"
)

type Components struct {
	Rc                 map[types.ComponentName]types.ReplicationControllerCreater
	//Rc                 map[types.ComponentName]func(replicationController *types.ReplicationController, blackduck *v1.Blackduck) types.ReplicationControllerInterface
	Service            map[types.ComponentName]types.ServiceCreater
	Configmap          map[types.ComponentName]types.ConfigmapCreater
	PVC                map[types.ComponentName]types.PvcCreater
	Secret             map[types.ComponentName]types.SecretCreater
	Size               map[types.ComponentName]types.SizeInterface
	ServiceAccount     map[types.ComponentName]types.ServiceAccountCreater
	ClusterRoleBinding map[types.ComponentName]types.ClusterRoleBindingCreater
}

var ComponentStore Components

func Register(name types.ComponentName, function interface{}) {
	//val := reflect.ValueOf(function)
	//typ := val.String()
	fmt.Println(reflect.TypeOf(function))
	switch function.(type) {
	case func(replicationController *types.ReplicationController, blackduck *v1.Blackduck) types.ReplicationControllerInterface:
		if ComponentStore.Rc == nil{
			ComponentStore.Rc = make(map[types.ComponentName]types.ReplicationControllerCreater)
		}
		ComponentStore.Rc[name] = function.(func(replicationController *types.ReplicationController, blackduck *v1.Blackduck) types.ReplicationControllerInterface)
	//case types.ReplicationControllerCreater:
	//	if ComponentStore.Rc == nil{
	//		ComponentStore.Rc = make(map[types.ComponentName]types.ReplicationControllerCreater)
	//	}
	//	ComponentStore.Rc[name] = function.(types.ReplicationControllerCreater)
	case types.ServiceCreater:
		ComponentStore.Service[name] = function.(types.ServiceCreater)
	case types.ConfigmapCreater:
		ComponentStore.Configmap[name] = function.(types.ConfigmapCreater)
	case types.PvcCreater:
		ComponentStore.PVC[name] = function.(types.PvcCreater)
	case types.SecretCreater:
		ComponentStore.Secret[name] = function.(types.SecretCreater)
	case types.SizeInterface:
		ComponentStore.Size[name] = function.(types.SizeInterface)
	case types.ServiceAccountCreater:
		ComponentStore.ServiceAccount[name] = function.(types.ServiceAccountCreater)
	case types.ClusterRoleBindingInterface:
		ComponentStore.ClusterRoleBinding[name] = function.(types.ClusterRoleBindingCreater)
	default:
		log.Fatal("Couldn't import " + name)
	}
}

