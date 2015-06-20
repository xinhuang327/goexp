package goexpr

import (
	"fmt"
	"testing"
)

type MyType struct {
	StringField     string
	NestField       *MyType
	unexportedField string
}

func Test(t *testing.T) {
	// src := `myVar.MethodName(1,"strVal", varName)`
	// expr := `myVar.NestField.StringField`
	expr := `myVar.unexportedField`

	myVar := MyType{}
	myVar.StringField = "MyStringFieldValue"
	myVar.NestField = &MyType{}
	myVar.NestField.StringField = "MyNestStringValue"
	myVar.unexportedField = "hehe"

	scope := NewEvaluateScope()
	scope.Variables["myVar"] = myVar
	result, err := scope.Evaluate(expr)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(result)
}
