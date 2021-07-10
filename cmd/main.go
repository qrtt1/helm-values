package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/qrtt1/friendly-yaml/internal/flatyaml"
	"os"
)

func main() {

	filename := flag.String("f", "values.yaml", "the path of a values file")
	filter := flag.String("e", "", "regular expression filters")
	showValues := flag.Bool("v", false, "show values")
	dumpToYaml := flag.Bool("y", false, "dump yaml (it works with the filter)")
	useShell := flag.Bool("i", false, "interactive shell")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s\n", flatyaml.GetVersion())
	}

	flag.Parse()

	_, err := os.Stat(*filename)
	if err != nil {
		flag.Usage()
		return
	}

	y := flatyaml.Values{}
	y.Load(*filename)

	shell := flatyaml.NewShell(y)
	if *useShell {
		runShell(y, shell)
		return
	}

	if *filter != "" {
		shell.ApplyFilter(*filter)
	}

	if *dumpToYaml {
		shell.CopyToBuffer(shell.CurrentFilter)
		shell.DumpToYaml()
	} else {
		shell.ShowConfigurations(*showValues, shell.CurrentFilter)
	}

	os.Exit(0)

}

func runShell(y flatyaml.Values, shell *flatyaml.Executor) {
	shell.ShowHelp()

	p := prompt.New(
		shell.Executor,
		shell.Complete,
		prompt.OptionTitle("helm-values: interactive helm values shell"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
