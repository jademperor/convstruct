package pkg

import (
	"go/ast"
)

// TypeSpec .
type TypeSpec struct {
	IsStruct bool
	Name     string
	Expr     ast.Expr
}

// GenStruct .
type GenStruct struct {
	Name   string
	Doc    string
	Fields []*GenField
}

// GenField .
type GenField struct {
	IsStar bool
	Name   string
	Expr   string
	Tag    string
}

// GenImport .
type GenImport struct {
	Name string
	Path string
}

// // StructType .
// type StructType struct {
// 	Name        string
// 	Node        *ast.StructType
// 	NeedImports []*ast.ImportSpec
// }

// // ParseToStructType .
// func ParseToStructType(st *ast.StructType) (*StructType, error) {
// 	newSt := &StructType{
// 		Name:        "",
// 		Node:        st,
// 		NeedImports: []*ast.ImportSpec{},
// 	}
// 	return newSt, nil
// }

// Package . to contains all package
type Package struct {
	Imports []*GenImport         // all path of packages those are imported
	Path    string               // self package path
	Name    string               // self package name
	Types   map[string]*TypeSpec // all package typs
}
