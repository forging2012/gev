package gev

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

type IUserRegistModel interface {
	IUserModel
	GetRegistorBody() interface{}
	RegistorJudge(body interface{}) error
	Registor(body interface{}) (*LoginData, error)
}

type UserRegistModel struct {
	UserModel `xorm:"extends"`
}

// 默认注册数据结构
type RegistorBody struct {
	Code     string `json:"code,omitempty" xorm:"not null"`
	Telphone string `json:"telphone,omitempty" xorm:"not null"`
	Password string `json:"password,omitempty" xorm:"not null"`
}

// 注册信息验证
func (u *UserRegistModel) GetRegistorBody() interface{} {
	return &RegistorBody{}
}

func (u *UserRegistModel) RegistorJudge(body interface{}) error {
	bean := u.Self()
	rbody := body.(*RegistorBody)
	ok, _ := Db.Where("telphone=?", rbody.Telphone).Get(bean)
	if ok {
		return errors.New("账号已注册")
	}
	if len(rbody.Password) < 6 || len(rbody.Password) > 32 {
		return errors.New("请输入6~32位密码")
	}
	return UserVerify.New().(IVerifyModel).JudgeCode(rbody.Telphone, rbody.Code)
}

func (u *UserRegistModel) Registor(body interface{}) (*LoginData, error) {
	bean := u.Self().(IUserRegistModel)
	if err := bean.RegistorJudge(body); err != nil {
		return nil, err
	}
	rbody := body.(*RegistorBody)
	u.Telphone = rbody.Telphone
	u.Password = bean.EncodePwd(rbody.Password)
	Db.InsertOne(bean)
	// 生成Token
	access := NewAccessToken(u.Id)
	return &LoginData{access, bean.GetDetail()}, nil
}

func (u *UserRegistModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		u.UserModel.Bind(g, u)
	} else {
		u.UserModel.Bind(g, self)
	}
	if UserVerify == nil {
		log.Println("userRegistmodel,没有设置UserVerify,忽略注册模块")
	} else {
		g.POST("/register", func(c *gin.Context) {
			user := u.New().(IUserRegistModel)
			body := user.GetRegistorBody()
			if err := c.BindJSON(body); err != nil {
				Err(c, 1, errors.New("JSON解析出错"))
			} else {
				data, err := user.Registor(body)
				if data != nil {
					c.SetCookie("X-AUTH-TOKEN", data.Access.Token, token_expire, "", "", false, false)
					data.Access.ReadContextInfo(c)
					Db.InsertOne(data.Access)
					Ok(c, data)
				} else {
					Err(c, 52, err)
				}
			}
		})
	}
}
