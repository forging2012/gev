package gev

import "errors"

type IBody interface {
	CopyTo(user IUserModel, data interface{}) error
	GetId() int
	IsNew() bool
}

type ModelBody struct {
	Id int `json:"id"`
}

func (m *ModelBody) IsNew() bool {
	return m.Id < 1
}

func (m *ModelBody) GetId() int {
	return m.Id
}

func (m *ModelBody) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*Model)
	data.Id = m.Id
	return nil
}

type UserModelBody struct {
	ModelBody
	Telphone string `json:"telphone,omitempty" xorm:"" gev:"用户账号，需要管理员权限"`
	Password string `json:"password,omitempty" xorm:"" gev:"用户密码，需要管理员权限"`
	Nickname string `gev:"用户昵称" json:"nickname" xorm:""`
}

func (this *UserModelBody) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*UserModel)
	if this.IsNew() {
		if role, ok := user.(IUserRoleModel); ok && role.IsAdmin() {
			if this.Telphone == "" || this.Password == "" {
				return errors.New("账号/密码不能为空")
			}
			data.Telphone = this.Telphone
			data.Password = data.EncodePwd(this.Password)
		} else {
			return errors.New("需要id")
		}
	}
	data.Nickname = this.Nickname
	return this.ModelBody.CopyTo(user, &data.Model)
}
