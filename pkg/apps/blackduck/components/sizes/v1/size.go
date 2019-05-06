package v1

import (
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/sizes"
)

var store = make(map[string]map[string]*sizes.Size)


func RegisterSize(name string, conf map[string]*sizes.Size){
	if _, ok := store[name]; ok{
		return
	}
	store[name] = conf
}

type v1 struct {}

func NewSize() sizes.SizeInterface {
	return &v1{}
}

func (*v1) GetSize(name string) map[string]*sizes.Size {
	val, ok := store[name]
	if ok {
		return val
	}
	return nil
}

