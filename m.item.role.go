package gev

// 自己和管理员可以编辑
type IItemRoleModel interface {
	IItemOwnerModel
}

type ItemRoleModel struct {
	ItemOwnerModel `xorm:"extends"`
}

func (o *ItemRoleModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
	}
	if user == nil {
		return false
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (o *ItemRoleModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
	}
	if user == nil {
		return false
	}
	if o.OwnerId == 0 {
		o.OwnerId = user.GetId()
		return true
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (m *ItemRoleModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.ItemOwnerModel.Bind(g, self)
}
