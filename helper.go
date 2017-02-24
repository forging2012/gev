package gev

import (
	"unsafe"

	"github.com/gin-gonic/gin"
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Ok(c *gin.Context, data interface{}) {
	c.IndentedJSON(200, gin.H{"code": 0, "data": data})
}
func Err(c *gin.Context, code int, err error) {
	msg := err.Error()
	if code == 0 {
		table := Str2bytes(msg)
		count := len(table)
		if count > 32 {
			count = 32
		}
		for i := 0; i < count; i++ {
			code += int(table[i])
		}
	}
	Log.Println("code:", code, "msg:", msg)
	c.IndentedJSON(200, gin.H{"code": code, "msg": msg})
}
func Api(c *gin.Context, data interface{}, err error) {
	if err != nil {
		Err(c, 0, err)
		return
	}
	Ok(c, data)
}
