package gev

import "errors"

type ISchemaBody interface {
	GetData(user IUserModel) (IModel, error)
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

func (m *ModelBody) GetData(user IUserModel) (IModel, error) {
	data := &Model{}
	data.Id = m.Id
	return data, nil
}

type UserModelBody struct {
	ModelBody
	Nickname string `gev:"用户昵称" json:"nickname" xorm:""`
}

func (b *UserModelBody) GetData(user IUserModel) (IModel, error) {
	data := &UserModel{}
	if b.IsNew() {
		return nil, errors.New("需要id")
	}
	model, err := b.ModelBody.GetData(user)
	data.Model = *(model.(*Model))
	if err != nil {
		return nil, err
	}
	data.Nickname = b.Nickname
	return data, nil
}
