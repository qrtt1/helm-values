package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FlatternYaml struct {
	settings map[string]interface{}
	yaml     map[interface{}]interface{}
}

func (y *FlatternYaml) load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &y.yaml)
	if err != nil {
		return err
	}

	y.settings = make(map[string]interface{})
	y.visit(y.yaml, "")

	return nil
}

func (y *FlatternYaml) visit(root map[interface{}]interface{}, parent string) {
	for k, v := range root {
		switch key := k.(type) {
		case string:
			switch value := v.(type) {
			case map[interface{}]interface{}:
				if parent == "" {
					y.visit(value, key)
				} else {
					y.visit(value, fmt.Sprintf("%s.%s", parent, key))
				}
			case bool, int, string, nil:
				if parent == "" {
					y.settings[key] = value
					//fmt.Printf("%s => %v\n", key, value)
				} else {
					y.settings[fmt.Sprintf("%s.%s", parent, key)] = value
					//fmt.Printf("%s => %v\n", fmt.Sprintf("%s.%s", parent, key), value)
				}
			case []interface{}:
				y.visitList(value, fmt.Sprintf("%s.%s", parent, key))
			default:
				fmt.Printf("Unhandled Type %s[%T]\n", key, key)
			}
		default:
			fmt.Printf("Unhandled Type %s[%T]\n", key, key)
		}

	}
}

func (y *FlatternYaml) visitList(list []interface{}, parent string) {
	for index, elem := range list {
		switch value := elem.(type) {
		case map[interface{}]interface{}:
			y.visit(value, fmt.Sprintf("%s[%d]", parent, index))
		case bool, int, string, nil:
			y.settings[fmt.Sprintf("%s[%d]", parent, index)] = value
		default:
			fmt.Printf(".??.%v TTT[%T]!\n", value, value)
		}
	}
}

func main() {

	y := FlatternYaml{}
	y.load("examples/values.yaml")

	for k, v := range y.settings {
		fmt.Printf("%s: %v\n", k, v)
	}

}
