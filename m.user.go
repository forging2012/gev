package gev

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type IUserModel interface {
	IItemRoleModel
	Login(telphone, password string) (*LoginData, error)
	ChangePassword(body interface{}) (*LoginData, error)
	JudgeChpwdCode(code string) error
	GetByTelphone(telphone string) bool
	EncodePwd(string) string
	Exist(telphone string) bool
}

type LoginBody struct {
	Telphone string `gev:"电话号码" json:"telphone" xorm:"varchar(32) unique(telphone) not null"`
	Password string `gev:"密码" json:"password" xorm:"varchar(64)"`
}

type UserModel struct {
	ItemRoleModel `xorm:"extends"`
	Nickname      string `gev:"用户昵称" json:"nickname" xorm:""`
	Telphone      string `gev:"电话号码" json:"telphone" xorm:"varchar(32) unique(telphone) not null"`
	Password      string `gev:"密码" json:"-" xorm:""`
}

func (this *UserModel) TableName() string {
	return "user"
}

func (this *UserModel) CanRead(user IUserModel) bool {
	if this.Id == user.GetId() {
		return true
	}
	return this.ItemRoleModel.CanWrite(user)
}
func (this *UserModel) CanWrite(user IUserModel) bool {
	if this.Id == user.GetId() {
		return true
	}
	return this.ItemRoleModel.CanWrite(user)
}

//  save时对密码进行加密
func (this *UserModel) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*UserModel)
	data.Password = this.Self().(IUserModel).EncodePwd(this.Password)
	return nil
}

func (this *UserModel) EncodePwd(password string) string {
	h := md5.New()
	h.Write([]byte(this.Telphone + password))
	hexText := make([]byte, 32)
	hex.Encode(hexText, h.Sum(nil))
	return string(hexText)
}

func (this *UserModel) Exist(telphone string) bool {
	if telphone == "" {
		return false
	}
	ok, _ := Db.Where("telphone=?", telphone).Get(this.Self())
	return ok
}

func (this *UserModel) GetByTelphone(telphone string) bool {
	if telphone == "" {
		return false
	}
	bean := this.Self().(IUserModel)
	ok, _ := Db.Where("telphone=?", telphone).Get(bean)
	return ok
}

func (this *UserModel) JudgeChpwdCode(code string) error {
	if this.Password != this.Self().(IUserModel).EncodePwd(code) {
		return errors.New("旧密码错误")
	}
	return nil
}

// 登录返回数据结构
type LoginData struct {
	Access *AccessToken `json:"access,omitempty" xorm:""`
	User   interface{}  `json:"user,omitempty" xorm:""`
}

// 登录
func (this *UserModel) Login(telphone, password string) (*LoginData, error) {
	bean := this.Self().(IUserModel)
	// 通过手机号查用户
	ok, err := Db.Where("telphone=?", telphone).Get(bean)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("用户不存在")
	}
	// 匹配密码
	if this.Password == bean.EncodePwd(password) {
		// 生成Token
		access := NewAccessToken(this.Id)
		data, _ := bean.GetInfo(bean, this.Id)
		return &LoginData{access, data}, nil
	}
	return nil, errors.New("密码不正确")
}

func (this *UserModel) ChangePassword(body interface{}) (*LoginData, error) {
	bean := this.Self().(IUserModel)
	rbody := body.(*RegistorBody)
	if len(rbody.Password) < 6 || len(rbody.Password) > 32 {
		return nil, errors.New("请输入6~32位密码")
	}
	if ok := bean.GetByTelphone(rbody.Telphone); !ok {
		return nil, errors.New("用户不存在")
	}
	if err := this.Self().(IUserModel).JudgeChpwdCode(rbody.Code); err != nil {
		return nil, err
	}
	this.Password = bean.EncodePwd(rbody.Password)
	_, err := Db.ID(this.Id).Update(bean)
	// 生成Token
	access := NewAccessToken(this.Id)
	return &LoginData{access, bean}, err
}

func (this *UserModel) Bind(g ISwagRouter, self IModel) {
	Db.Sync(new(AccessToken))
	if self == nil {
		self = this
	}
	this.ItemRoleModel.Bind(g, self)
	g.Info("登录", "在header中加入以下两项以作统计\n`X-DEVICE`:ios/android/web\n`X-UUID`: 设备唯一标识 \n登录成功后返回用户信息和token").Body(
		map[string]interface{}{"telphone": "", "password": ""},
	).Data(
		&LoginData{User: self.GetData()},
	).POST("/login", func(c *gin.Context) {
		body := &LoginBody{}
		if err := c.BindJSON(body); err != nil {
			Err(c, 1, errors.New("JSON解析出错"))
			return
		}
		this.Telphone = body.Telphone
		this.Password = body.Password
		data, err := this.New().(IUserModel).Login(this.Telphone, this.Password)
		if data != nil {
			c.SetCookie("X-AUTH-TOKEN", data.Access.Token, token_expire, "", "", false, false)
			data.Access.Logined(c)
			Ok(c, data)
		} else {
			Err(c, 0, err)
		}
	})
	g.Info("验证码修改密码", "").Body(
		self.(IUserRegistModel).GetRegistorBody(),
	).Data(
		&LoginData{User: self},
	).POST("/change/password", func(c *gin.Context) {
		user := this.New().(IUserRegistModel)
		body := user.GetRegistorBody()
		if err := c.BindJSON(body); err != nil {
			Err(c, 1, errors.New("JSON解析出错"))
		} else {
			data, err := user.ChangePassword(body)
			if data != nil {
				c.SetCookie("X-AUTH-TOKEN", data.Access.Token, token_expire, "", "", false, false)
				data.Access.PasswordChanged(c)
				Ok(c, data)
			} else {
				Err(c, 0, err)
			}
		}
	})
	g.Info("我的信息", "").Data(
		self.GetData(),
	).GET("/mine/info", func(c *gin.Context) {
		if user, ok := NeedAuth(c); ok {
			data := user.(IUserModel).GetData()
			if bean, ok := data.(IData); ok {
				Ok(c, bean.GetDetail(user.(IUserModel)))
				return
			}
			Ok(c, user)
		}
	})
}

func (this *UserModel) MiddleWare(c *gin.Context) {
	// 当前登录用户数据
	token := c.Query("access_token")
	if token != "" {
		now := time.Now()
		user := this.Self().(IModel).New()
		ok, _ := Db.Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
		if ok {
			c.Set("user", user)
		}
	}
}
