package gev

import (
	"github.com/gin-gonic/gin"
	// "github.com/golibs/uuid"

	"time"
)

type AccessToken struct {
	Model     `xorm:"extends"`
	Token     string    `json:"token,omitempty" xorm:"not null"`
	ExpiredAt time.Time `json:"expired_at,omitempty" xorm:"not null"`
	UserId    int       `json:"user_id,omitempty" xorm:"not null"`
	Ip        string    `json:"ip,omitempty" xorm:"not null"`
	UA        string    `json:"UA,omitempty" xorm:"not null"`
	Device    string    `json:"device,omitempty" xorm:"not null"`
	Uuid      string    `json:"uuid,omitempty" xorm:"not null"`
}

func (a *AccessToken) ReadContextInfo(c *gin.Context) {
	UA := c.Request.Header.Get("User-Agent")
	a.Ip = c.ClientIP()
	a.UA = UA
	a.Device = c.Request.Header.Get("X-DEVICE")
	a.Uuid = c.Request.Header.Get("X-UUID")
}

func NewAccessToken(user_id int) *AccessToken {
	return &AccessToken{
		UserId:    user_id,
		Token:     "uuid.Rand().Hex()",
		ExpiredAt: time.Now().Add(time.Duration(token_expire) * time.Second),
	}
}
