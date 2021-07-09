package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func visitList(list []interface{}, parent string) {
	for index, elem := range list {
		switch value := elem.(type) {
		case map[interface{}]interface{}:
			visit(value, fmt.Sprintf("%s[%d]", parent, index), nil)
		case bool, int, string, nil:
			fmt.Printf("%s => %v\n", fmt.Sprintf("%s[%d]", parent, index), value)
		default:
			fmt.Printf(".??.%v TTT[%T]!\n", value, value)
		}
	}
}

func visit(root map[interface{}]interface{}, parent string, content *map[string]interface{}) {
	for k, v := range root {
		switch key := k.(type) {
		case string:
			switch value := v.(type) {
			case map[interface{}]interface{}:
				if parent == "" {
					visit(value, key, nil)
				} else {
					visit(value, fmt.Sprintf("%s.%s", parent, key), nil)
				}
			case bool, int, string, nil:
				if parent == "" {
					//*content[key] = value
					fmt.Printf("%s => %v\n", key, value)
				} else {
					fmt.Printf("%s => %v\n", fmt.Sprintf("%s.%s", parent, key), value)
				}
			case []interface{}:
				visitList(value, fmt.Sprintf("%s.%s", parent, key))
			default:
				fmt.Printf("Unhandled Type %s[%T]\n", key, key)
			}
		default:
			// each type should be the type of string
			fmt.Printf("Unhandled Type %s[%T]\n", key, key)
		}

	}
}

func main() {
	var content map[interface{}]interface{}
	data, err := ioutil.ReadFile("examples/values.yaml")
	var flatternContent map[string]interface{}
	if err != nil {
		fmt.Printf("%v", err)
	} else {
		err := yaml.Unmarshal(data, &content)
		if err == nil {
			visit(content, "", &flatternContent)
		} else {
			fmt.Printf("%v", err)
		}

	}
}
