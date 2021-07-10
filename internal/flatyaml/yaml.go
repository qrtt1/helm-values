package flatyaml

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Values struct {
	Settings map[string]interface{}
	yaml     map[interface{}]interface{}
	suggests []prompt.Suggest
}

func (y *Values) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &y.yaml)
	if err != nil {
		return err
	}

	y.Settings = make(map[string]interface{})
	y.visit("", y.yaml)

	y.suggests = make([]prompt.Suggest, 1000)
	for k, v := range y.Settings {
		y.suggests = append(y.suggests, prompt.Suggest{Text: k, Description: fmt.Sprintf("%v", v)})
	}

	return nil
}

func (y *Values) visit(parent string, root map[interface{}]interface{}) {
	for k, v := range root {
		switch key := k.(type) {
		case string:
			switch value := v.(type) {
			case map[interface{}]interface{}:
				y.updateSetting(parent, key, value)
				if parent == "" {
					y.visit(key, value)
				} else {
					y.visit(fmt.Sprintf("%s.%s", parent, key), value)
				}
			case []interface{}:
				y.updateSetting(parent, key, value)
				y.visitList(fmt.Sprintf("%s.%s", parent, key), value)
			default:
				y.updateSetting(parent, key, value)
			}
		default:
			fmt.Printf("Unhandled Type %s[%T]\n", key, key)
		}

	}
}

func (y *Values) visitList(parent string, list []interface{}) {
	for index, elem := range list {
		switch value := elem.(type) {
		case map[interface{}]interface{}:
			y.Settings[fmt.Sprintf("%s[%d]", parent, index)] = value
			y.visit(fmt.Sprintf("%s[%d]", parent, index), value)
		default:
			y.Settings[fmt.Sprintf("%s[%d]", parent, index)] = value
		}
	}
}

func (y *Values) updateSetting(parent string, key string, value interface{}) {
	if parent == "" {
		y.Settings[key] = value
	} else {
		y.Settings[fmt.Sprintf("%s.%v", parent, key)] = value
	}
}
