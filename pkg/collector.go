package pkg

import (
	"go/ast"
)

// CollectTypes .
func CollectTypes(node ast.Node) map[string]*TypeSpec {
	typs := make(map[string]*TypeSpec)
	collectFunc := func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			name := x.Name.Name
			expr := x.Type
			iss := isStruct(expr)
			DebugF("[CollectTypes] collect type: Name-[%s], Expr-[%v], IsStruct-[%v]", name, expr, iss)
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
