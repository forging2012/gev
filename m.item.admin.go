package gev

// 自己和管理员可以编辑
type IItemAdminModel interface {
	IItemOwnerModel
}

type ItemAdminModel struct {
	ItemOwnerModel `xorm:"extends"`
}

func (o *ItemAdminModel) CanRead(user IUserModel) bool {
	if u, ok := user.(IUserAdminModel); ok && u.IsAdmin(o.Self().(IModel)) {
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

func (o *ItemAdminModel) CanWrite(user IUserModel) bool {
	if u, ok := user.(IUserAdminModel); ok && u.IsAdmin(o.Self().(IModel)) {
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

func (m *ItemAdminModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.ItemOwnerModel.Bind(g, self)
}
