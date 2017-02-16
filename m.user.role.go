package gev

import (
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
	Role            string `gev:"用户角色" json:"role,omitempty" xorm:"not null default '普通用户'"`
}

func (u *UserRoleModel) BeforeInsert() {
	u.Role = u.Self().(IUserRoleModel).GetRole()
}

func (s *UserRoleModel) SearchSession(session *xorm.Session, condition ISearch) {
	session.Where("role=?", s.Self().(IUserRoleModel).GetRole())
	s.UserRegistModel.SearchSession(session, condition)
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
