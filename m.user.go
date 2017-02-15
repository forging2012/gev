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
	Nickname    string `gev:"用户昵称" json:"nickname,omitempty" xorm:""`
	Telphone    string `gev:"电话号码" json:"telphone,omitempty" xorm:""`
	Password    string `gev:"密码" json:"password,omitempty" xorm:""`
}

func (u *UserModel) TableName() string {
	return "user"
}

// 登录返回数据结构
type LoginData struct {
	Access *AccessToken `json:"access,omitempty" xorm:""`
	User   interface{}  `json:"user,omitempty" xorm:""`
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
		self = u
	}
	u.SearchModel.Bind(g, self)
	g.Info("登录", "在header中加入以下两项以作统计\n`X-DEVICE`:ios/android/web\n`X-UUID`: 设备唯一标识 \n登录成功后返回用户信息和token").Body(
		map[string]interface{}{"telphone": "", "password": ""},
	).Data(
		&LoginData{User: self},
	).POST("/login", func(c *gin.Context) {
		if err := c.BindJSON(u); err != nil {
			Err(c, 1, errors.New("JSON解析出错"))
			return
		}
		data, err := u.New().(IUserModel).Login(u.Telphone, u.Password)
		if data != nil {
			c.SetCookie("X-AUTH-TOKEN", data.Access.Token, token_expire, "", "", false, false)
			data.Access.Logined(c)
			Ok(c, data)
		} else {
			Err(c, 0, err)
		}
	})
	g.Info("我的信息", "").Data(
		self.GetDetail(),
	).GET("/mine/info", func(c *gin.Context) {
		if user, ok := c.Get("user"); !ok {
			NeedAuth(c)
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
		ok, _ := Db.Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
		if ok {
			c.Set("user", user)
		}
	}
}
