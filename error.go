package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func NeedAuth(c *gin.Context) {
	Err(c, 1255, errors.New("需要登录"))
}
