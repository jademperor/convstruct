package pkg

import (
	"go/ast"
)

// CollectTypes .
func CollectTypes(node ast.Node) map[string]*TypeSpec {
	typs := make(map[string]*TypeSpec)
	collectFunc := func(n ast.Node) bool {
		// var (
		// 	expr ast.Expr
		// 	name string
		// )

		switch x := n.(type) {
		case *ast.TypeSpec:
			name := x.Name.Name
			expr := x.Type
			iss := isStruct(expr)
			DebugF("collect type: Name-[%s], Expr-[%v], IsStruct-[%v]", name, expr, iss)
			typs[name] = &TypeSpec{
				Name:     name,
				IsStruct: iss,
				Expr:     expr,
			}
		}

		return true
	}

	ast.Inspect(node, collectFunc)
	return typs
}

func isStruct(expr ast.Expr) bool {
	_, ok := expr.(*ast.StructType)
	return ok
}

// // CollectStructs collects and maps structType nodes to their positions
// func CollectStructs(node ast.Node) map[string]*StructType {
// 	structs := make(map[string]*StructType, 0)
// 	collectStructs := func(n ast.Node) bool {
// 		var t ast.Expr
// 		var structName string

// 		switch x := n.(type) {
// 		case *ast.TypeSpec:
// 			if x.Type == nil {
// 				return true
// 			}
// 			structName = x.Name.Name
// 			t = x.Type
// 		case *ast.CompositeLit:
// 			t = x.Type
// 		case *ast.ValueSpec:
// 			structName = x.Names[0].Name
// 			t = x.Type
// 		}

// 		x, ok := t.(*ast.StructType)
// 		if !ok {
// 			return true
// 		}

// 		structs[structName] = &StructType{
// 			Name: structName,
// 			Node: x,
// 		}
// 		return true
// 	}
// 	ast.Inspect(node, collectStructs)
// 	return structs
// }
