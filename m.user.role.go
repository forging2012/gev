package gev

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

var (
	// 用于判断登录者类型
	user_roles = make([]IUserRoleModel, 0, 5)
)

type IUserRoleModel interface {
	IUserRegistModel
	IsAdmin() bool
	GetRole() string
}

type UserRoleModel struct {
	UserRegistModel `xorm:"extends"`
	Role            string `gev:"用户角色" json:"role,omitempty" xorm:"varchar(32) unique(telphone) not null default '普通用户'"`
}

func (u *UserRoleModel) BeforeInsert() {
	u.Role = u.Self().(IUserRoleModel).GetRole()
}
func (u *UserRoleModel) BeforeUpdate() {
	u.Role = ""
}

func (this *UserRoleModel) Search(user IUserModel, condition ISearch) (interface{}, error) {
	bean := this.self
	return GetSearchData(bean, condition, func(session *xorm.Session) {
		session.Where("role=?", this.Self().(IUserRoleModel).GetRole())
		condition.MakeSession(session)
	})
}

func (u *UserRoleModel) GetRole() string {
	return u.Role
}

func (u *UserRoleModel) IsAdmin() bool {
	if u.Role == "管理员" {
		return true
	}
	return false
}

// 登录
func (u *UserRoleModel) Login(telphone, password string) (*LoginData, error) {
	bean := u.Self().(IUserRoleModel)
	// 通过手机号查用户
	ok, err := Db.Where("telphone=? and role=?", telphone, bean.GetRole()).Get(bean)
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
		return &LoginData{access, bean}, nil
	}
	return nil, errors.New("密码不正确")
}

func (this *UserRoleModel) Exist(telphone string) bool {
	if telphone == "" {
		return false
	}
	ok, _ := Db.Where("telphone=? and role=?", telphone, this.Self().(IUserRoleModel).GetRole()).Get(this.Self())
	return ok
}

func (u *UserRoleModel) ChangePassword(body interface{}) (*LoginData, error) {
	bean := u.Self().(IUserRoleModel)
	rbody := body.(*RegistorBody)
	if len(rbody.Password) < 6 || len(rbody.Password) > 32 {
		return nil, errors.New("请输入6~32位密码")
	}
	if err := UserVerify.New().(IVerifyModel).JudgeCode(rbody.Telphone, rbody.Code); err != nil {
		return nil, err
	}
	ok, _ := Db.Where("telphone=? and role=?", rbody.Telphone, bean.GetRole()).Get(bean)
	if !ok {
		return nil, errors.New("用户不存在")
	}
	u.Telphone = rbody.Telphone
	u.Password = bean.EncodePwd(rbody.Password)
	_, err := Db.ID(u.Id).Update(bean)
	// 生成Token
	access := NewAccessToken(u.Id)
	return &LoginData{access, bean}, err
}

func (u *UserRoleModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = u
	} else {
		// 加入用户类型类型
		user_roles = append(user_roles, self.(IUserRoleModel))
	}
	u.UserRegistModel.Bind(g, self)
}

func (u *UserRoleModel) MiddleWare(c *gin.Context) {
	// 当前登录用户数据
	token := c.Query("access_token")
	if token != "" {
		now := time.Now()
		user := &UserRoleModel{}
		ok, _ := Db.Cols("id", "role").Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
		if ok {
			// 判断登录者类型
			for _, item := range user_roles {
				if user.Role == item.GetRole() {
					bean := item.New()
					Db.ID(user.Id).Get(bean)
					c.Set("user", bean)
					break
				}
			}
		}
	}
}
