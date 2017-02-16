package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func NeedAuth(c *gin.Context) (interface{}, bool) {
	if user, ok := c.Get("user"); ok {
		return user, true
	} else {
		Err(c, 1255, errors.New("需要登录"))
		return nil, false
	}
}

func NeedAuthRole(c *gin.Context, role string) (interface{}, bool) {
	if user, ok := c.Get("user"); ok {
		if v, ok := user.(IUserRoleModel); ok && v.GetRole() == role {
			return user, true
		}
		Err(c, 1256, errors.New("需要"+role+"权限"))
		return nil, false
	} else {
		Err(c, 1255, errors.New("需要登录"))
		return nil, false
	}
}
