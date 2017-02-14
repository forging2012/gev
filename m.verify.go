package gev

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type IVerifyModel interface {
	IModel
	NewVerifyCode(title string) error
	SendCode(title, code string) error
	RandCode() string
	JudgeCode(title, code string) error
}

// 验证码模型
type VerifyModel struct {
	Model `xorm:"extends"`
	Title string `json:"title,omitempty" xorm:"not null"`
	Code  string `json:"code,omitempty" xorm:"unique not null"`
	// 剩余验证次数
	Rest int `json:"rest,omitempty" xorm:"not null default 10"`
}

func (v *VerifyModel) Disable() {
	if v.Id > 0 {
		v.Rest = 0
		Db.ID(v.Id).Cols("rest").Update(v.Self())
	}
}

func (v *VerifyModel) NewVerifyCode(title string) error {
	bean := v.Self().(IVerifyModel)
	ok, _ := Db.Where("title=?", title).Get(bean)
	if v.UpdateAt.Add(time.Minute).After(time.Now()) {
		return errors.New("发送太频繁")
	}
	v.Title = title
	v.Code = bean.RandCode()
	v.Rest = 10
	err := bean.SendCode(title, v.Code)
	if err != nil {
		return err
	}
	if ok {
		Db.Where("title=?", title).Update(bean)
	} else {
		Db.InsertOne(bean)
	}
	return nil
}

func (v *VerifyModel) JudgeCode(title, code string) error {
	bean := v.Self().(IVerifyModel)
	ok, _ := Db.Where("title=?", title).Get(bean)
	if !ok {
		return errors.New("尚未发送验证码")
	}
	if v.UpdateAt.Add(10 * time.Minute).Before(time.Now()) {
		return errors.New("验证码已过期")
	}
	if v.Rest < 1 {
		return errors.New("验证码已失效")
	}
	if v.Code != code {
		v.Rest--
		Db.ID(v.Id).Cols("rest").Update(bean)
		return errors.New("验证码错误")
	}
	return nil
}

func (v *VerifyModel) SendCode(title, code string) error {
	fmt.Println(title, "=>", code)
	return nil
}

func (v *VerifyModel) RandCode() string {
	code := make([]byte, 4)
	for i := 0; i < 4; i++ {
		code[i] = byte('0' + rand.Intn(10))
	}
	return string(code)
}

func (v *VerifyModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		v.Model.Bind(g, v)
	} else {
		v.Model.Bind(g, self)
	}
	g.Info("发送验证码", "测试环境默认验证码1024").Params(
		g.PathParam("title", "手机号/邮箱"),
	).GET("/send/:title", func(c *gin.Context) {
		title := c.Param("title")
		err := v.New().(IVerifyModel).NewVerifyCode(title)
		Api(c, nil, err)
	})
}
