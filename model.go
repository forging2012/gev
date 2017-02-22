package gev

import (
	"log"
	"reflect"
	"time"
)

type Class interface {
	Self() Class
	SetSelf(Class)
	New() Class
}

// 基础数据模型
// 父类可以通过Self()获取实例
type IModel interface {
	Class
	ISchemaBody
	GetDetail() interface{}
	GetSearch() interface{}
	GetBody() ISchemaBody
	Bind(g ISwagRouter, model IModel)
}

type Model struct {
	self     IModel    `xorm:"-"`
	Id       int       `json:"id,omitempty" xorm:"pk autoincr"`
	CreateAt time.Time `json:"-" xorm:"created"`
	UpdateAt time.Time `json:"-" xorm:"updated"`
}

func (m *Model) Self() Class {
	// if m.self == nil {
	// 	return m
	// }
	return m.self
}
func (m *Model) SetSelf(self Class) {
	m.self = self.(IModel)
}
func (m *Model) IsNew() bool {
	return m.Id < 1
}
func (m *Model) New() Class {
	// if m.self == nil {
	// 	m.self = m
	// }
	model := reflect.New(reflect.TypeOf(m.self).Elem()).Interface().(Class)
	model.SetSelf(model)
	return model
}
func (m *Model) GetDetail() interface{} {
	return m.Self()
}
func (m *Model) GetSearch() interface{} {
	return m.Self()
}
func (m *Model) GetBody() ISchemaBody {
	return m.Self().(ISchemaBody)
}
func (m *Model) GetId() int {
	return m.Id
}
func (m *Model) GetData(IUserModel) (IModel, error) {
	return m.Self().(IModel), nil
}

func (m *Model) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		log.Fatalln("model.Bind需要 self")
	}
	m.self = self
	err := Db.Sync(self)
	if err != nil {
		log.Printf("%T-->%v\n", self, err)
	}
}
