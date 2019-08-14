package controllers_utils

import (
	"github.com/blackducksoftware/synopsys-operator/meta-builder/flying-dutchman"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"strings"
)

func CreateInstructionManual(instructionManualLocation string) (*flying_dutchman.RuntimeObjectDependencyYaml, error) {
	// Read Dependency YAML File into Struct
	var content []byte

	if uri, err := url.Parse(instructionManualLocation); err == nil && (strings.EqualFold(uri.Scheme, "https") || strings.EqualFold(uri.Scheme, "http")) {
		content, err = HttpGet(instructionManualLocation)
		if err != nil {
			return nil, err
		}
	} else {
		content, err = ioutil.ReadFile(instructionManualLocation)
		if err != nil {
			return nil, err
		}
	}

	dependencyYamlStruct := &flying_dutchman.RuntimeObjectDependencyYaml{}
	err := yaml.Unmarshal(content, dependencyYamlStruct)
	if err != nil {
		return nil, err
	}
	return dependencyYamlStruct, nil
}
