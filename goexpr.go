package goexpr

import (
	"fmt"
	"go/ast"
	// "go/format"
	"errors"
	"go/parser"
	// "go/token"
	"reflect"
)

type EvaluateScope struct {
	Variables map[string]interface{}
}

func NewEvaluateScope() *EvaluateScope {
	scope := &EvaluateScope{}
	scope.Variables = make(map[string]interface{})
	return scope
}

func (scope *EvaluateScope) Evaluate(expr string) (interface{}, error) {
	// fmt.Println("(goexpr) Evaluating", expr)
	// fset := token.NewFileSet()
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}

	// ast.Print(fset, node)

	result, err := scope.walk(node)

	return result, err
}

func (scope *EvaluateScope) walk(node ast.Expr) (interface{}, error) {

	if ident, ok := node.(*ast.Ident); ok {
		if scope.Variables[ident.Name] == nil {
			return nil, errors.New(fmt.Sprint("variable not found:", ident.Name))
		}
		return scope.Variables[ident.Name], nil
	}

	if selExpr, ok := node.(*ast.SelectorExpr); ok {
		obj, err := scope.walk(selExpr.X)
		if err != nil {
			return nil, err
		}
		selName := selExpr.Sel.Name

		objVal := reflect.ValueOf(obj)
		objValKind := objVal.Kind()

		var fieldVal reflect.Value
		if objValKind == reflect.Struct {
			fieldVal = objVal.FieldByName(selName)
		} else if objValKind == reflect.Ptr {
			fieldVal = objVal.Elem().FieldByName(selName)
		} else {
			return nil, errors.New(fmt.Sprint("not supported obj kind:", objValKind))
		}

		if fieldVal.IsValid() {
			if fieldVal.CanInterface() {
				return fieldVal.Interface(), nil
			} else {
				fmt.Println("Field", selName, "can not get interface, so return string")
				return fieldVal.String(), nil
			}
		} else {
			return nil, errors.New(fmt.Sprint("Field not valid:", selName))
		}
	}

	return nil, errors.New(fmt.Sprint("not supported node:", node))
}
