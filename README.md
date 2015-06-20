# goexpr
Evaluate golang expression dynamically


```go
func Test(t *testing.T) {
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
```
