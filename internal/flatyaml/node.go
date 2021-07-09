package flatyaml

import (
	"strings"
)

func RebuildYaml(root map[interface{}]interface{}, key string, content interface{}) map[interface{}]interface{} {
	n := root
	s := strings.Split(key, ".")
	buildNode(n, s, content)
	return n
}

func buildNode(n map[interface{}]interface{}, path []string, content interface{}) {
	if len(path) == 0 {
		return
	}

	if len(path) == 1 {
		n[path[0]] = content
		return
	}

	if _, ok := n[path[0]]; !ok {
		n[path[0]] = make(map[interface{}]interface{})
	}

	buildNode(n[path[0]].(map[interface{}]interface{}), path[1:], content)
}
