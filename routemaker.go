package gev

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/inu1255/annotation"
)

// function: convert path,query,body to param
type paramLoader func(*gin.Context) reflect.Value

// function: choose paramLoader for param
type paramManage func(param *Param)

type Param struct {
	Name   string
	In     string // param in path,query,body or file
	Desc   string
	Type   reflect.Type
	New    paramLoader
	Method *Method
}

type Method struct {
	Doc        map[string][]string // remment info
	Name       string
	StructName string // "" if not a struct function
	Func       reflect.Value
	Params     []*Param
	OutType    reflect.Type
}

// httpmethod from remment like below
// @method DELETE
// or return POST if contains body Param
func (this *Method) HttpMethod() string {
	if ss, ok := this.Doc["method"]; ok {
		return strings.ToUpper(ss[len(ss)-1])
	}
	for _, param := range this.Params {
		if param.In == "body" || param.In == "file" {
			return "POST"
		}
	}
	return "GET"
}

// param description from remment like below
// @param paramName paramDescs
func (this *Method) ParamDesc(paramName string) string {
	if ss, ok := this.Doc["param"]; ok {
		for _, s := range ss {
			ss := strings.Fields(s)
			if len(ss) > 0 && ss[0] == paramName {
				return strings.Join(ss[1:], " ")
			}
		}
	}
	return ""
}

// route path from remment like below
// @path /foo/bar
// or return path made by methodName and pathParams
func (this *Method) Path() (path string) {
	if ss, ok := this.Doc["path"]; ok {
		path = ss[len(ss)-1]
	} else {
		//  FooBar --> /foo/bar
		re := regexp.MustCompile(`([a-z]|^)[A-Z]`)
		path = "/" + re.ReplaceAllStringFunc(this.Name, nameToRoute)
	}
	for _, param := range this.Params {
		if param.In == "path" {
			path += "/:" + param.Name
		}
	}
	return
}

// return struct type if this is a struct method
func (this *Method) RecvType() reflect.Type {
	if this.StructName == "" || len(this.Params) < 1 {
		return nil
	}
	return this.Params[0].Type
}

type RouteMaker struct {
	Methods []*Method

	calls []paramManage
	fm    FM
}

// add paramManage function LIFO
// the last add paramManage match condition will take effect
func (this *RouteMaker) AddParamManage(fn paramManage) {
	this.calls = append(this.calls, fn)
}

// read route information from one function
func (this *RouteMaker) AddMethod(f interface{}) *Method {
	if f == nil {
		return nil
	}
	mFunc := funcValue(f)
	typ := mFunc.Type()
	no := typ.NumOut()
	if no != 2 || typ.Out(1).Kind() != reflect.Interface || !typ.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil
	}
	info := this.fm.ReadFunc(f)
	method := new(Method)
	method.Doc = annotation.Doc2Map(info.Doc)
	method.Name = info.Name
	method.StructName = info.StructName
	method.Func = mFunc
	ni := typ.NumIn()
	method.Params = make([]*Param, ni)
	param_start_index := 0
	if method.StructName != "" && ni > 0 {
		param_start_index = 1
		param := &Param{}
		param.Type = typ.In(0)
		param.New = newInstCall(param.Type.Elem())
		method.Params[0] = param
	}
	for i := param_start_index; i < ni; i++ {
		in := typ.In(i)
		param := &Param{}
		param.Name = info.Params[i-param_start_index]
		param.Type = in
		param.Desc = method.ParamDesc(param.Name)
		param.Method = method
		for j := len(this.calls) - 1; j >= 0; j-- {
			if this.calls[j](param); param.New != nil {
				break
			}
		}
		method.Params[i] = param
	}
	if no > 0 {
		method.OutType = typ.Out(0)
	}
	this.Methods = append(this.Methods, method)
	return method
}

// read route information,save to this.Methods
func (this *RouteMaker) AddRoute(f interface{}) {
	typ := reflect.TypeOf(f)
	switch typ.Kind() {
	case reflect.Func:
		this.AddMethod(f)
	case reflect.Interface, reflect.Ptr:
		count := typ.NumMethod()
		for i := 0; i < count; i++ {
			this.AddMethod(typ.Method(i))
		}
	}
}

// add route to app
func (this *RouteMaker) RouteAll(app IRouter) {
	this.fm.Save()
	for _, method := range this.Methods {
		if desc, ok := method.Doc["desc"]; ok {
			app.Info(desc...)
		}
		app.Data(reflect.New(method.OutType).Interface())
		count := len(method.Params)
		calls := make([]paramLoader, count)
		for i := 0; i < count; i++ {
			param := method.Params[i]
			switch param.In {
			case "path":
				app.PathParam(param.Name, param.Desc)
			case "query":
				app.QueryParam(param.Name, param.Desc)
			case "body":
				app.Body(reflect.New(param.Type).Interface())
			case "file":
				app.FileParam(param.Name, param.Desc)
			}
			calls[i] = param.New
		}
		if method.StructName == "" {
			handler := makeHandlerFunc(method, calls)
			app.Handle(method.HttpMethod(), method.Path(), handler)
		} else {
			typ := method.RecvType()
			if typ != nil && typ.Implements(reflect.TypeOf(new(IService)).Elem()) {
				handler := makeServiceHandlerFunc(method, calls)
				app.Group(strings.ToLower(method.StructName)).Handle(method.HttpMethod(), method.Path(), handler)
			} else {
				handler := makeHandlerFunc(method, calls)
				app.Group(strings.ToLower(method.StructName)).Handle(method.HttpMethod(), method.Path(), handler)
			}
		}
	}
}

// aB --> a/b or A --> a
// help FooBar --> foo/bar
func nameToRoute(from string) string {
	if len(from) == 2 {
		return from[:1] + "/" + strings.ToLower(from[1:])
	}
	return strings.ToLower(from)
}

// convert func or reflect.Method to reflect.Value
func funcValue(f interface{}) reflect.Value {
	if m, ok := f.(reflect.Method); ok {
		return m.Func
	}
	return reflect.ValueOf(f)
}

// default paramManage
func DefaultManage(param *Param) {
	switch param.Type.Kind() {
	case reflect.String:
		param.In = "query"
		param.New = newQueryString(param.Name)
	case reflect.Int:
		param.In = "query"
		param.New = newQueryInt(param.Name)
	case reflect.Int64:
		param.In = "query"
		param.New = newQueryInt64(param.Name)
	case reflect.Float32:
		param.In = "query"
		param.New = newQueryFloat32(param.Name)
	case reflect.Float64:
		param.In = "query"
		param.New = newQueryFloat64(param.Name)
	case reflect.Struct, reflect.Ptr, reflect.Slice, reflect.Map:
		if param.Type.String() == "*multipart.FileHeader" {
			param.In = "file"
			param.New = newMultiFile(param.Name)
		} else {
			param.In = "body"
			param.New = newJsonCall(param.Type)
		}
	case reflect.Func:
		param.New = newNilCall(param.Type)
	}
	return
}

// self param for struct
func SelfManage(param *Param) {
	if param.Name != "self" {
		return
	}
	stype := param.Method.RecvType()
	if stype == nil {
		return
	}
	switch param.Type.Kind() {
	case reflect.Slice:
		if param.Type.Elem().Kind() == reflect.Interface {
			typ := reflect.SliceOf(stype)
			param.Type = typ
			param.New = newJsonArrayCall(typ, stype)
			param.In = "body"
		}
	case reflect.Interface:
		param.Type = stype
		param.New = newJsonCall(stype)
		param.In = "body"
	}
	return
}

func NewRouteMaker() *RouteMaker {
	m := new(RouteMaker)
	m.Methods = make([]*Method, 0, 20)
	m.calls = make([]paramManage, 0, 2)
	m.AddParamManage(DefaultManage)
	m.AddParamManage(SelfManage)
	m.fm = make(FM)
	m.fm.Restore()
	return m
}
