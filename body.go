package gev

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
