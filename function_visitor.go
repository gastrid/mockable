package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type functionVisitor struct {
	ifName        string
	structName    string
	fset          *token.FileSet
	varName       string
	functions     []functionData
	cmap          ast.CommentMap
	oldGenComment string
	newGenComment string
}

type functionData struct {
	Name    string
	Params  *ast.FieldList
	Results *ast.FieldList
	Doc     *ast.CommentGroup
}

func (v *functionVisitor) Visit(n ast.Node) ast.Visitor {
	// If the declaration is a function and not a method
	switch node := n.(type) {
	case *ast.File:
		return v
	case *ast.FuncDecl:
		if node.Recv == nil {
			// We add the function signature to our list
			oldNode := node
			fn := functionData{
				Name:    node.Name.Name,
				Params:  node.Type.Params,
				Results: node.Type.Results,
				Doc:     node.Doc,
			}
			v.functions = append(v.functions, fn)

			// then we transform our function to a method
			node.Recv = &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("mck"),
						},
						Type: &ast.StarExpr{
							X: ast.NewIdent(v.structName),
						},
					},
				},
			}

			v.cmap.Update(oldNode, n)
		}
	}

	return nil
}

func (v *functionVisitor) generateMockable() []ast.Decl {
	decls := []ast.Decl{
		v.createInterfaceDeclaration(),
		v.createStructDeclaration(),
		v.createVarDeclaration(),
		v.createSetFunction(),
	}

	for _, fn := range v.functions {
		decls = append(decls, v.createFunctionLayer(fn))
	}

	return decls
}

// createFunctionLayer will create functions that will call the struct's associated method
func (v *functionVisitor) createFunctionLayer(fn functionData) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fn.Name),
		Type: &ast.FuncType{
			Params:  fn.Params,
			Results: fn.Results,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(v.varName),
								Sel: ast.NewIdent(fn.Name),
							},
							Args: paramNames(fn.Params),
						},
					},
				},
			},
		},
	}
}

func paramNames(fl *ast.FieldList) []ast.Expr {
	params := make([]ast.Expr, len(fl.List))

	for i := 0; i < len(fl.List); i++ {
		// TODO What if there are more names? When are there more names?
		params[i] = ast.NewIdent(fl.List[i].Names[0].Name)
	}

	return params
}

// We're generating the interface with all the methods we've registered in v.Visit
func (v *functionVisitor) createInterfaceDeclaration() *ast.GenDecl {
	methods := []*ast.Field{}

	for _, fn := range v.functions {

		m := &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(fn.Name),
			},
			Type: &ast.FuncType{
				Params:  fn.Params,
				Results: fn.Results,
			},
		}
		methods = append(methods, m)
	}

	ifDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(v.ifName),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: methods,
					},
				},
			},
		},
	}

	return ifDecl
}

// We create a struct
func (v *functionVisitor) createStructDeclaration() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name:   ast.NewIdent(v.structName),
				Assign: token.NoPos,
				Type: &ast.StructType{
					Fields:     &ast.FieldList{},
					Incomplete: false,
				},
			},
		},
	}
}

// We create a variable of type interfaceName that is a struct
func (v *functionVisitor) createVarDeclaration() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent(v.varName),
				},
				Type: ast.NewIdent(v.ifName),
				Values: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: ast.NewIdent(v.structName),
						},
					},
				},
			},
		},
	}
}

// This function is to easily replace the struct with a mock one
func (v *functionVisitor) createSetFunction() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("Set%s", v.ifName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("newS"),
						},
						Type: ast.NewIdent(v.ifName),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: ast.NewIdent(v.ifName),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("old"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						ast.NewIdent(v.varName),
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(v.varName),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("newS"),
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("old"),
					},
				},
			},
		},
	}
}
