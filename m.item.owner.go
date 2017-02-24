package gev

import "errors"

type IItemOwnerModel interface {
	IItemModel
	DeleteIds(user IUserRoleModel, ids []int) error
}

type ItemOwnerModel struct {
	ItemModel `xorm:"extends"`
	OwnerId   int `gev:"-" json:"-" xorm:""`
}

func (o *ItemOwnerModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (o *ItemOwnerModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if o.Id < 1 && o.OwnerId == 0 {
		o.OwnerId = user.GetId()
		return true
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (i *ItemOwnerModel) DeleteIds(user IUserRoleModel, ids []int) error {
	if len(ids) < 1 {
		return errors.New("数组长度不能为0")
	}
	_, err := Db.In("id", ids).Where("owner_id=?", user.GetId()).Delete(i.Self())
	return err
}

func (m *ItemOwnerModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.ItemModel.Bind(g, self)
}
