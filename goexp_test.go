package goexpr

import (
	"fmt"
	"testing"
)

type MyType struct {
	StringField     string
	NestField       *MyType
	unexportedField string
	SliceField      []string
	MapField        map[string]MyType
}

func (m MyType) MyFunc(strA string, intA int, fA float64) string {
	fmt.Println("Calling MyType.MyFunc", strA, intA, fA)
	return m.StringField
}

func Test(t *testing.T) {
	// expr := `myVar.NestField.StringField`
	// expr := `len(myVar.unexportedField)`
	// expr := `myVar.MyFunc("str args", 4, 3.2)`
	// expr := `myVar.SliceField[1]`
	expr := `myVar.MapField["strKey"]`

	myVar := MyType{}
	myVar.StringField = "MyStringFieldValue"
	myVar.NestField = &MyType{}
	myVar.NestField.StringField = "MyNestStringValue"
	myVar.unexportedField = "hehe"
	myVar.SliceField = []string{"Hello", "World"}
	myVar.MapField = map[string]MyType{
		"strKey": myVar,
	}
	Debug = true

	scope := NewEvaluateScope()
	scope.Variables["myVar"] = myVar
	result, err := scope.Evaluate(expr)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(result)
}
