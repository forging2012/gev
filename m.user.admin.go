package gev

type IUserAdminModel interface {
	IUserModel
	IsAdmin(model IModel) bool
}

type UserAdminModel struct {
	UserModel `xorm:"extends"`
	Role      string `json:"role,omitempty" xorm:"not null default '普通用户'"`
}

func (u *UserAdminModel) IsAdmin(model IModel) bool {
	if u.Role == "管理员" {
		return true
	}
	return false
}

func (u *UserAdminModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		u.UserModel.Bind(g, u)
	} else {
		u.UserModel.Bind(g, self)
	}
}
