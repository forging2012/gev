package gev

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// only new
func newInstCall(t reflect.Type) paramLoader {
	return func(*gin.Context) reflect.Value {
		return reflect.New(t)
	}
}

// new and copy fields
func copyInstCall(src interface{}) paramLoader {
	return func(*gin.Context) reflect.Value {
		return CloneValue(reflect.ValueOf(src))
	}
}
func CloneValue(src reflect.Value) reflect.Value {
	var dst reflect.Value
	if src.Kind() == reflect.Ptr || src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	dst = reflect.New(src.Type())
	cloneValue(src, dst.Elem())
	return dst
}
func cloneValue(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		n := src.NumField()
		for i := 0; i < n; i++ {
			f := src.Field(i)
			cloneValue(f, dst.Field(i))
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
}

func newString(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		return reflect.ValueOf(ctx.Param(index))
	}
}
func newInt(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.Atoi(ctx.Param(index)); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newInt64(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseInt(ctx.Param(index), 10, 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newFloat32(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Param(index), 32); e == nil {
			return reflect.ValueOf(float32(r))
		}
		return reflect.ValueOf(float32(0))
	}
}
func newFloat64(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Param(index), 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0.0)
	}
}
func newQueryString(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		return reflect.ValueOf(ctx.Query(index))
	}
}
func newQueryInt(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.Atoi(ctx.Query(index)); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newQueryInt64(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseInt(ctx.Query(index), 10, 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newQueryFloat32(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Query(index), 32); e == nil {
			return reflect.ValueOf(float32(r))
		}
		return reflect.ValueOf(float32(0))
	}
}
func newQueryFloat64(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Query(index), 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0.0)
	}
}
func newMultiFile(index string) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		_, header, _ := ctx.Request.FormFile(index)
		return reflect.ValueOf(header)
	}
}
func newJsonCall(t reflect.Type) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		v := reflect.New(t)
		if err := ctx.BindJSON(v.Interface()); err == nil {
			return v.Elem()
		} else {
			return reflect.New(t.Elem())
		}
	}
}
func newJsonArrayCall(t reflect.Type, t2 reflect.Type) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		v := reflect.New(t)
		if err := ctx.BindJSON(v.Interface()); err == nil {
			v = v.Elem()
			v2 := reflect.MakeSlice(t2, v.Len(), v.Cap())
			for i := 0; i < v.Len(); i++ {
				v2.Index(i).Set(v.Index(i))
			}
			return v2
		} else {
			panic(err)
		}
	}
}
func newNilCall(t reflect.Type) paramLoader {
	return func(ctx *gin.Context) reflect.Value {
		return reflect.Zero(t)
	}
}
