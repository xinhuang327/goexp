package goexpr

import (
	"fmt"
	"go/ast"
	// "go/format"
	"errors"
	"go/parser"
	"go/token"
	"reflect"
)

var _ = token.NewFileSet
var _ = ast.Print

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
	fset := token.NewFileSet()
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}

	ast.Print(fset, node)

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

	if funExpr, ok := node.(*ast.CallExpr); ok {
		var directFuncName string
		var methodVal reflect.Value
		if funIdent, ok := funExpr.Fun.(*ast.Ident); ok {
			directFuncName = funIdent.Name
		} else if funSel, ok := funExpr.Fun.(*ast.SelectorExpr); ok {
			methodValResult, err := scope.walk(funSel)
			if err != nil {
				return nil, err
			}
			methodVal = methodValResult.(reflect.Value)
		}
		var args []reflect.Value
		for _, a := range funExpr.Args {
			result, err := scope.walk(a)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(result))
		}
		fmt.Println("Args:", args)
		if directFuncName != "" {
			if directFuncName == "len" {
				if len(args) > 0 {
					return args[0].Len(), nil
				} else {
					return nil, errors.New("No args provided.")
				}
			}
		} else if methodVal.IsValid() {
			returnVals := methodVal.Call(args)
			var returnObjs []interface{}
			for _, rv := range returnVals {
				returnObjs = append(returnObjs, rv.Interface())
			}
			if len(returnObjs) > 1 {
				return returnObjs, nil
			} else if len(returnObjs) == 1 {
				return returnObjs[0], nil
			} else {
				return nil, nil
			}
		} else {
			return nil, errors.New("Cannot find any method to call.")
		}
	}

	if selExpr, ok := node.(*ast.SelectorExpr); ok {
		obj, err := scope.walk(selExpr.X)
		if err != nil {
			return nil, err
		}
		selName := selExpr.Sel.Name

		objVal := reflect.ValueOf(obj)
		objValKind := objVal.Kind()

		var memberVal reflect.Value
		if objValKind == reflect.Struct {
			memberVal = objVal.FieldByName(selName)
			if !memberVal.IsValid() {
				memberVal = objVal.MethodByName(selName)
			}
		} else if objValKind == reflect.Ptr {
			memberVal = objVal.Elem().FieldByName(selName)
			if !memberVal.IsValid() {
				memberVal = objVal.Elem().MethodByName(selName)
			}
		} else {
			return nil, errors.New(fmt.Sprint("not supported obj kind:", objValKind))
		}

		if memberVal.IsValid() {
			if memberVal.Kind() != reflect.Func {
				if memberVal.CanInterface() {
					return memberVal.Interface(), nil
				} else {
					fmt.Println("Field", selName, "can not get interface, so return reflect.Value")
					return memberVal, nil
				}
			} else {
				return memberVal, nil // return func value
			}
		} else {
			return nil, errors.New(fmt.Sprint("Field not valid:", selName))
		}
	}

	return nil, errors.New(fmt.Sprint("not supported node:", node))
}
