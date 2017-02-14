package gev

type IUserAdminModel interface {
	IUserRegistModel
	IsAdmin(model IModel) bool
}

type UserAdminModel struct {
	UserRegistModel `xorm:"extends"`
	Role            string `gev:"用户角色" json:"role,omitempty" xorm:"not null default '普通用户'"`
}

func (u *UserAdminModel) IsAdmin(model IModel) bool {
	if u.Role == "管理员" {
		return true
	}
	return false
}

func (u *UserAdminModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = u
	}
	u.UserRegistModel.Bind(g, self)
}
