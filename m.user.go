package gev

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type IUserModel interface {
	ISearchModel
	Login(telphone, password string) (*LoginData, error)
	EncodePwd(string) string
}

type UserModel struct {
	SearchModel `xorm:"extends"`
	Nickname    string `json:"nickname,omitempty" xorm:"not null"`
	Telphone    string `json:"telphone,omitempty" xorm:"not null"`
	Password    string `json:"password,omitempty" xorm:"not null"`
}

// 登录返回数据结构
type LoginData struct {
	Access *AccessToken `json:"access,omitempty" xorm:"not null"`
	User   interface{}  `json:"user,omitempty" xorm:"not null"`
}

// 不返回密码
func (u *UserModel) GetSearch() interface{} {
	u.Password = ""
	return u.Model.GetSearch()
}

// 不返回密码
func (u *UserModel) GetDetail() interface{} {
	u.Password = ""
	return u.Model.GetDetail()
}

//  save时对密码进行加密
func (r *UserModel) GetData() (IModel, error) {
	r.Password = r.EncodePwd(r.Password)
	return r.Model.GetData()
}

// 登录
func (u *UserModel) Login(telphone, password string) (*LoginData, error) {
	bean := u.Self().(IUserModel)
	// 通过手机号查用户
	ok, err := Db.Where("telphone=?", telphone).Get(bean)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("用户不存在")
	}
	// 匹配密码
	if u.Password == bean.EncodePwd(password) {
		// 生成Token
		access := NewAccessToken(u.Id)
		return &LoginData{access, bean.GetDetail()}, nil
	}
	return nil, errors.New("密码不正确")
}

func (u *UserModel) EncodePwd(password string) string {
	h := md5.New()
	h.Write([]byte(u.Telphone + password))
	hexText := make([]byte, 32)
	hex.Encode(hexText, h.Sum(nil))
	return string(hexText)
}

func (u *UserModel) Bind(g ISwagRouter, self IModel) {
	Db.Sync(new(AccessToken))
	if self == nil {
		u.SearchModel.Bind(g, u)
	} else {
		u.SearchModel.Bind(g, self)
	}
	g.POST("/login", func(c *gin.Context) {
		if err := c.BindJSON(u); err != nil {
			Err(c, 1, errors.New("JSON解析出错"))
			return
		}
		data, err := u.New().(IUserModel).Login(u.Telphone, u.Password)
		if data != nil {
			c.SetCookie("X-AUTH-TOKEN", data.Access.Token, token_expire, "", "", false, false)
			data.Access.ReadContextInfo(c)
			Db.InsertOne(data.Access)
			Ok(c, data)
		} else {
			Err(c, 52, err)
		}
	})
	g.GET("/mine/info", func(c *gin.Context) {
		if user, ok := c.Get("user"); ok {
			Err(c, 1255, errors.New("需要登录"))
		} else {
			Ok(c, user.(IUserModel).GetDetail())
		}
	})
}

func (u *UserModel) MiddleWare(c *gin.Context) {
	// 当前登录用户数据
	token := c.Query("access_token")
	if token != "" {
		now := time.Now()
		user := u.Self().(IModel).New()
		ok, _ := Db.Where("id in (select user_id from access_token where token=? and expired_at<?)", token, now).Get(user)
		if ok {
			c.Set("user", user)
		}
	}
}
