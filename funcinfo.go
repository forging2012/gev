package gev

import (
	"encoding/json"
	"go/ast"
	"io/ioutil"
	"reflect"

	"github.com/inu1255/annotation"
)

type FuncInfo struct {
	StructName string
	Name       string
	Doc        string
	Params     []string
}

type FM map[string]*FuncInfo

func (this FM) Restore() bool {
	if data, err := ioutil.ReadFile("fm.json"); err == nil {
		err = json.Unmarshal(data, this)
		return err == nil
	}
	return false
}

func (this FM) Save() bool {
	if data, err := json.Marshal(this); err == nil {
		err = ioutil.WriteFile("fm.json", data, 0644)
		return err == nil
	}
	return false
}

func (this FM) ReadFunc(f interface{}) *FuncInfo {
	if f == nil {
		return nil
	}
	var funcDecl *ast.FuncDecl
	var methodName, structPkg, structName string
	switch v := f.(type) {
	case reflect.Method:
		typ := isStruct(v.Type.In(0))
		methodName, structPkg, structName = v.Name, typ.PkgPath(), typ.Name()
		funcDecl = annotation.GetFuncByMethod(v)
	default:
		methodName, structPkg, structName = annotation.GetFuncInfo(f)
		funcDecl = annotation.FindFunc(methodName, structPkg, structName)
	}
	key := structPkg + "." + structName + "." + methodName
	// fmt.Println(key, funcDecl)
	if funcDecl == nil {
		return this[key]
	}
	info := &FuncInfo{
		StructName: structName,
		Name:       methodName,
		Doc:        funcDecl.Doc.Text(),
		Params:     annotation.GetParams(funcDecl),
	}
	this[key] = info
	return info
}

func isStruct(typ reflect.Type) reflect.Type {
	if typ == nil {
		return nil
	}
	switch typ.Kind() {
	case reflect.Interface, reflect.Ptr:
		return isStruct(typ.Elem())
	case reflect.Struct:
		return typ
	}
	return nil
}
