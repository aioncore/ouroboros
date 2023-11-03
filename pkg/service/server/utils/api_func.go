package utils

import (
	"reflect"
	"strings"
)

type APIFunc struct {
	F        reflect.Value  // underlying api function
	Args     []reflect.Type // type of each function arg
	Returns  []reflect.Type // type of each return arg
	ArgNames []string       // name of each argument
}

// return a function's argument types
func funcArgTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumIn()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.In(i)
	}
	return typez
}

// return a function's return types
func funcReturnTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumOut()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.Out(i)
	}
	return typez
}

func NewAPIFunc(f interface{}, args string) *APIFunc {
	return newAPIFunc(f, args)
}

func newAPIFunc(f interface{}, args string) *APIFunc {
	var argNames []string
	if args != "" {
		argNames = strings.Split(args, ",")
	}

	r := &APIFunc{
		F:        reflect.ValueOf(f),
		Args:     funcArgTypes(f),
		Returns:  funcReturnTypes(f),
		ArgNames: argNames,
	}
	return r
}
