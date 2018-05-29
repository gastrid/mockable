package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

func main() {
	o := &options{}
	o.setup()

	flag.Parse()

	if o.ifName == "" || o.fileName == "" {
		fmt.Printf("You must specify an interface name and the name of the file to be made mockable.")
		flag.Usage()
		os.Exit(2)
	}

	// We generate the original comment string we'll be looking for
	// before validation
	o.generateOriginalComment()

	// Some standard name validation
	o.validate()

	// We generate the names for the struct type and variable
	o.setStructNameAndVarName()

	o.generateNewComment()

	generateMock(o)
}

func generateMock(o *options) {

	// Parse and look for the requested file
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(fileinfo os.FileInfo) bool {
		return fileinfo.Name() == o.fileName
	}, parser.ParseComments)

	if err != nil {
		fmt.Errorf("Failed to parse package path. %v", err)
		os.Exit(2)
	}

	if len(pkgs) != 1 {
		fmt.Errorf("That's weird, there should only be one package with this name")
		os.Exit(2)
	}

	var targetFile *ast.File

	for _, p := range pkgs {
		for i, f := range p.Files {
			if i == o.fileName {
				targetFile = f
			}
		}
	}

	if targetFile == nil {
		fmt.Errorf("Oops, file name not found in this folder")
		os.Exit(2)
	}

	// First, we want to work out if there already is an interface with this name in the document
	iv := &interfaceVisitor{
		ifName: o.ifName,
	}
	ast.Walk(iv, targetFile)

	if iv.abort == true {
		// If the interface name aready exists, we don't need to regenerate the whole thing and exit nicely
		return
	}

	v := &functionVisitor{
		ifName:     o.ifName,
		structName: o.structName,
		varName:    o.varName,
		fset:       fset,
	}

	// looping over comment groups to replace the
	// old generate comment with the new one
	for _, cg := range targetFile.Comments {
		for _, c := range cg.List {
			if c.Text == o.originalComment {
				c.Text = o.newComment
			}
		}
	}

	v.cmap = ast.NewCommentMap(fset, targetFile, targetFile.Comments)

	// We look for all independent functions, store their signature and turn them into methods
	ast.Walk(v, targetFile)

	// We generate the interface, struct, setter function and function layer
	MockDecls := v.generateMockable()
	targetFile.Decls = append(targetFile.Decls, MockDecls...)
	// I think this is useful?
	targetFile.Comments = v.cmap.Filter(targetFile).Comments()
	var buf bytes.Buffer
	err = format.Node(&buf, fset, targetFile)
	ioutil.WriteFile(o.fileName, []byte(buf.String()), 0666)
}
