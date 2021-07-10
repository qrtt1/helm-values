package flatyaml

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Executor struct {
	values Values
	buffer map[string]interface{}
	helps  []prompt.Suggest

	CurrentFilter FilterFunc
}

type FilterFunc func(string) bool

func AllPassFilter(string) bool {
	return true
}

func NewShell(values Values) *Executor {
	buffer := make(map[string]interface{})
	helps := createHelps()

	return &Executor{values, buffer, helps, AllPassFilter}
}

func createHelps() []prompt.Suggest {
	helps := []prompt.Suggest{
		{Text: ":?", Description: "Show this help"},
		{Text: ":a", Description: "Show all configuration keys"},
		{Text: ":av", Description: "Show all configuration keys and values"},
		{Text: ":b", Description: "Show selected keys in the buffer"},
		{Text: ":d", Description: "Dump the buffer to yaml"},
		{Text: ":cb", Description: "Copy filtered keys to buffer"},
		{Text: ":rb", Description: "Reset buffer"},
		{Text: ":reset", Description: "Reset filter to all pass mode"},
		{Text: ":regex", Description: "Apply Regex Filter"},
		{Text: ":q", Description: "Bye !"},
	}
	return helps
}

func (e *Executor) ResetFilter() {
	e.CurrentFilter = AllPassFilter
}

func (e *Executor) ResetBuffer() {
	e.buffer = make(map[string]interface{})
}

func (e *Executor) Complete(d prompt.Document) []prompt.Suggest {
	if d.Text == "" {
		return []prompt.Suggest{}
	}

	if d.GetWordBeforeCursor() == "" {
		return []prompt.Suggest{}
	}

	if strings.HasPrefix(d.Text, ":") {
		return prompt.FilterContains(e.helps, d.GetWordBeforeCursor(), true)
	}

	s := e.values.suggests
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

func (e *Executor) ShowHelp() {
	fmt.Println("Any lines start with \":\" will consider a command")
	for _, help := range e.helps {
		fmt.Printf("%s\t%s\n", help.Text, help.Description)
	}
	fmt.Println()
	fmt.Println("Other inputs will put into the buffer when they match any keys")
	fmt.Println()
}

func (e *Executor) Executor(s string) {
	if s == "" {
		return
	}

	if strings.HasPrefix(s, ":regex ") {
		expr := strings.TrimSpace(s[len(":regex "):])
		compile, err := regexp.Compile(expr)
		if err != nil {
			fmt.Printf("Invalid regular expression %v", expr)
			return
		}
		e.CurrentFilter = func(s string) bool {
			return compile.MatchString(s)
		}
		fmt.Printf("Update filter to [%v] %p\n", expr, e.CurrentFilter)
		return
	}
	if strings.HasPrefix(s, ":") {
		switch s {
		case ":?":
			e.ShowHelp()
		case ":a":
			e.ShowConfigurations(false, e.CurrentFilter)
		case ":av":
			e.ShowConfigurations(true, e.CurrentFilter)
		case ":cb":
			e.CopyToBuffer(e.CurrentFilter)
		case ":rb":
		case ":b":
			e.ShowBuffer()
		case ":d":
			e.DumpToYaml()
		case ":reset":
			e.ResetFilter()
		case ":q":
			os.Exit(0)
		default:
			fmt.Printf("Unknown command %s\n\n", s)
			e.ShowHelp()
		}
		return
	}

	if val, ok := e.values.Settings[s]; ok {
		e.buffer[s] = val
		fmt.Printf("Put the key [%s] into the buffer\n", s)
	} else {
		fmt.Printf("Drop the unknown key [%s]\n", s)
	}

}

func (e *Executor) CopyToBuffer(filter FilterFunc) {
	m := e.values.Settings
	for _, k := range getSortedKeys(m, filter) {
		if filter(k) {
			e.buffer[k] = m[k]
		}
	}
}

func (e *Executor) ShowConfigurations(showValues bool, filter FilterFunc) {
	m := e.values.Settings
	for _, k := range getSortedKeys(m, filter) {
		if showValues {
			fmt.Printf("%s = %v\n", k, m[k])
		} else {
			fmt.Printf("%s\n", k)
		}
	}
}

func (e *Executor) ShowBuffer() {
	m := e.buffer
	fmt.Println("Selected configurations:")
	for _, k := range getSortedKeys(m, AllPassFilter) {
		fmt.Printf("%s = %v\n", k, m[k])
	}
	fmt.Println()
}

func (e *Executor) DumpToYaml() {
	root := make(map[interface{}]interface{})
	for k, v := range e.buffer {
		RebuildYaml(root, k, v)
	}
	d, err := yaml.Marshal(root)
	if err == nil {
		fmt.Printf("%v\n", string(d))
	}
}

func (e *Executor) ApplyFilter(expr string) {
	compile, err := regexp.Compile(expr)
	if err != nil {
		fmt.Printf("Invalid regular expression %v", expr)
		return
	}
	e.CurrentFilter = func(s string) bool {
		return compile.MatchString(s)
	}
}

func getSortedKeys(m map[string]interface{}, filter FilterFunc) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if filter != nil {
			if filter(k) {
				keys = append(keys, k)
			}
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
