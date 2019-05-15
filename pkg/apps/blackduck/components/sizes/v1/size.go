package v1

import (
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
)

var store = make(map[string]map[string]*types.Size)


func RegisterSize(name string, conf map[string]*types.Size){
	if _, ok := store[name]; ok{
		return
	}
	store[name] = conf
}

type v1 struct {}

func NewSize()types.SizeInterface {
	return &v1{}
}

func (*v1) GetSize(name string) map[string]*types.Size {
	val, ok := store[name]
	if ok {
		return val
	}
	return nil
}

