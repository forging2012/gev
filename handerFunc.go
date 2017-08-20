package gev

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// 生成 IService 相关的 HanderFunc
func makeServiceHandlerFunc(m *Method, call []paramLoader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		n := len(call)
		params := make([]reflect.Value, n)
		params[0] = call[0](ctx)
		service := params[0].Interface().(IService)
		defer func() {
			err := recover()
			service.Finish(err)
		}()
		if service.Before(ctx) {
			for i := 1; i < n; i++ {
				params[i] = call[i](ctx)
			}
			Log.Println(params[0].Type(), m.Name, params[1:])
			out := m.Func.Call(params)
			data := out[0].Interface()
			msg := out[1].Interface()
			if msg == nil {
				service.After(data, nil)
			} else {
				service.After(data, msg.(error))
			}
		}
	}
}

// 生成 interface{} 相关的 HanderFunc
func makeHandlerFunc(m *Method, call []paramLoader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		n := len(call)
		params := make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			params[i] = call[i](ctx)
		}
		Log.Println(m.Name, m.Func.Type(), params)
		out := m.Func.Call(params)
		data := out[0].Interface()
		err := out[1].Interface()
		if err == nil {
			respOk(ctx, data)
		} else {
			respApi(ctx, data, err.(error))
		}
	}
}
