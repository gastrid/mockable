package main

import (
	"go/ast"
)

type interfaceVisitor struct {
	ifName string
	abort  bool
}

func (v *interfaceVisitor) Visit(n ast.Node) ast.Visitor {
	// We are only looking at top-level nodes so no need to follow into children nodes
	switch node := n.(type) {
	case *ast.File:
		return v
	case *ast.GenDecl:
		for _, spec := range node.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					if typeSpec.Name.String() == v.ifName {
						// Here, we've found an interface with the same name as the interface
						// specified in the go generate command
						// so we abort the generation
						v.abort = true
					}
				}
			}
		}
	}

	return nil
}
