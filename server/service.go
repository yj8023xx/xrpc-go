package server

import (
	"go/ast"
	"reflect"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

type serviceMethod struct {
	method    reflect.Method
	argsType  reflect.Type
	replyType reflect.Type
}

type service struct {
	name      string
	typ       reflect.Type
	rcvr      reflect.Value
	methodMap map[string]*serviceMethod
}

func NewService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.typ = reflect.TypeOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	if !ast.IsExported(s.name) {

	}
	s.registerMethods()
	return s
}

func (s *service) registerMethods() {
	s.methodMap = make(map[string]*serviceMethod)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		argsType, replyType := mType.In(1), mType.In(2)
		if argsType.Kind() != reflect.Ptr {
			continue
		}
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		if !isExportedOrBuiltinType(argsType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		if returnType := mType.Out(0); returnType != errorType {
			continue
		}
		s.methodMap[method.Name] = &serviceMethod{
			method:    method,
			argsType:  argsType.Elem(),
			replyType: replyType.Elem(),
		}
	}
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *service) call(m *serviceMethod, argv, replyv reflect.Value) error {
	f := m.method.Func
	ret := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if err := ret[0].Interface(); err != nil {
		return err.(error)
	}
	return nil
}
