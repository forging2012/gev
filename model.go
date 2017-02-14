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

type ISchemaBody interface {
	GetData() (IModel, error)
}

// 基础数据模型
// 父类可以通过Self()获取实例
type IModel interface {
	Class
	ISchemaBody
	GetDetail() interface{}
	GetSearch() interface{}
	GetBody() ISchemaBody
	GetId() int
	Bind(g ISwagRouter, model IModel)
}

type Model struct {
	self     IModel    `xorm:"-"`
	Id       int       `json:"id,omitempty" xorm:"pk autoincr"`
	CreateAt time.Time `json:"create_at,omitempty" xorm:"created"`
	UpdateAt time.Time `json:"update_at,omitempty" xorm:"updated"`
}

func (m *Model) Self() Class {
	return m.self
}
func (m *Model) SetSelf(self Class) {
	m.self = self.(IModel)
}
func (m *Model) New() Class {
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
func (m *Model) GetData() (IModel, error) {
	return m.Self().(IModel), nil
}

func (m *Model) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		log.Fatalln("model.Bind需要 self")
	}
	m.self = self
	Db.Sync(self)
}
