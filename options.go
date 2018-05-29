package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type options struct {
	ifName          string
	fileName        string
	structName      string
	varName         string
	originalComment string
	newComment      string
}

func (o *options) setup() {
	flag.StringVar(&o.ifName, "interface", "", "The name of the interface to be created; Must be specified.")
	flag.StringVar(&o.fileName, "filename", "", "The name of the file to read; Must be specified.")
}

func (o *options) setStructNameAndVarName() {
	o.structName = fmt.Sprintf("%sImpl", o.ifName)
	o.varName = fmt.Sprintf("default%s", o.ifName)
}

func (o *options) validate() {
	if o.ifName == "" {
		fmt.Errorf("interface name must be specified")
		os.Exit(2)
	}

	o.ifName = strings.Title(o.ifName)

	if o.fileName == "" {
		fmt.Errorf("file name must be specified")
		os.Exit(2)
	}

	if !strings.HasSuffix(o.fileName, ".go") {
		fmt.Errorf("specified file must be a go file")
		os.Exit(2)
	}
}

// This is a long shot and will be used to identify the go gen comment and replace it with a new one that will
// generate phil's genmock one
func (o *options) generateOriginalComment() {
	o.originalComment = fmt.Sprintf("//go:generate mockable -interface=%s -filename=%s", o.ifName, o.fileName)
}

func (o *options) generateNewComment() {
	o.newComment = fmt.Sprintf("//go:generate genmock -interface=%s -mock-package=. -package=.", o.ifName)

}
