package flatyaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNodeBuilder(t *testing.T) {
	values := Values{}
	values.Load("../testdata/primehub-values.yaml")

	key := "datasetUpload.interface.webFrontEndImage.resources"
	root := make(map[interface{}]interface{})
	actual := RebuildYaml(root, key, values.Settings[key])
	fmt.Printf("%v\n", actual)

	var expected map[interface{}]interface{}
	loadYaml("../testdata/datasetUpload-resources.yaml", &expected)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, \nActual %v", expected, actual)
	}

}

func loadYaml(file string, output *map[interface{}]interface{}) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, &output)
	}
}
