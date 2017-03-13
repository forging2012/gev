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

func (this *UserRoleModel) BeforeInsert() {
	this.Role = this.Self().(IUserRoleModel).GetRole()
}
func (this *UserRoleModel) BeforeUpdate() {
	this.Role = ""
}

func (this *UserRoleModel) Search(user IUserModel, condition ISearch) (interface{}, error) {
	bean := this.self
	return GetSearchData(bean, user, condition, func(session *xorm.Session) {
		session.Where("role=?", this.Self().(IUserRoleModel).GetRole())
		condition.MakeSession(user, session)
	})
}

func (this *UserRoleModel) GetRole() string {
	return this.Role
}

func (this *UserRoleModel) IsAdmin() bool {
	if this.Role == "管理员" {
		return true
	}
	return false
}

// 登录
func (this *UserRoleModel) Login(telphone, password string) (*LoginData, error) {
	bean := this.Self().(IUserRoleModel)
	// 通过手机号查用户
	ok, err := Db.Where("telphone=? and role=?", telphone, bean.GetRole()).Get(bean)
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

func (this *UserRoleModel) Exist(telphone string) bool {
	if telphone == "" {
		return false
	}
	ok, _ := Db.Where("telphone=? and role=?", telphone, this.Self().(IUserRoleModel).GetRole()).Get(this.Self())
	return ok
}

func (this *UserRoleModel) GetByTelphone(telphone string) bool {
	if telphone == "" {
		return false
	}
	bean := this.Self().(IUserModel)
	ok, _ := Db.Where("telphone=? and role=?", telphone, this.Self().(IUserRoleModel).GetRole()).Get(bean)
	return ok
}
func (this *UserRoleModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = this
	} else {
		// 加入用户类型类型
		user_roles = append(user_roles, self.(IUserRoleModel))
	}
	this.UserRegistModel.Bind(g, self)
}

func (this *UserRoleModel) MiddleWare(c *gin.Context) {
	// 当前登录用户数据
	token := c.Query("access_token")
	if token != "" {
		now := time.Now()
		user := &UserRoleModel{}
		ok, _ := Db.Cols("id", "role").Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
		if ok {
			// Log.Println(len(user_roles))
			// 判断登录者类型
			for _, item := range user_roles {
				// Log.Println(item.GetRole())
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
