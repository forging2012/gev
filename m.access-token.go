package gev

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golibs/uuid"
)

type AccessToken struct {
	Model     `xorm:"extends"`
	Token     string    `gev:"身份密钥" json:"token,omitempty" xorm:"index"`
	ExpiredAt time.Time `gev:"过期时间" json:"expired_at,omitempty" xorm:"index"`
	UserId    int       `json:"-" xorm:""`
	Ip        string    `json:"-" xorm:""`
	UA        string    `json:"-" xorm:""`
	Device    string    `json:"-" xorm:""`
	Uuid      string    `json:"-" xorm:""`
	Action    string    `json:"-" xorm:""`
}

func (a *AccessToken) ReadContextInfo(c *gin.Context) {
	UA := c.Request.Header.Get("User-Agent")
	a.Ip = c.ClientIP()
	a.UA = UA
	a.Device = c.Request.Header.Get("X-DEVICE")
	a.Uuid = c.Request.Header.Get("X-UUID")
}

func (a *AccessToken) Save(c *gin.Context) {
	a.ReadContextInfo(c)
	Db.InsertOne(a)
}

func (a *AccessToken) Logined(c *gin.Context) {
	a.Action = "login"
	a.Save(c)
	if a.Id > 0 {
		Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=? and device=?", a.Id, a.UserId, a.Device)
	}
}

func (a *AccessToken) PasswordChanged(c *gin.Context) {
	a.Action = "password"
	a.Save(c)
	if a.Id > 0 {
		Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=?", a.Id, a.UserId)
	}
}

func Logout(user_id int) {
	if user_id > 0 {
		Db.Exec("update access_token set expired_at='1993-03-07' where user_id=?", user_id)
	}
}

func NewAccessToken(user_id int) *AccessToken {
	return &AccessToken{
		UserId:    user_id,
		Token:     uuid.Rand().Hex(),
		ExpiredAt: time.Now().Add(time.Duration(token_expire) * time.Second),
	}
}
