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
	IBind
	GetId() int
	IsNew() bool
	GetBody() IBody
	GetData() interface{}
}

type Model struct {
	self     IModel    `xorm:"-"`
	Id       int       `json:"id,omitempty" xorm:"pk autoincr"`
	CreateAt time.Time `json:"create_at,omitempty" xorm:"created"`
	UpdateAt time.Time `json:"-" xorm:"updated"`
}

// Class 接口
func (this *Model) Self() Class {
	// if this.self == nil {
	// 	return m
	// }
	return this.self
}
func (this *Model) SetSelf(self Class) {
	this.self = self.(IModel)
}
func (this *Model) New() Class {
	// if this.self == nil {
	// 	this.self = m
	// }
	model := reflect.New(reflect.TypeOf(this.self).Elem()).Interface().(Class)
	model.SetSelf(model)
	return model
}

func (this *Model) GetId() int {
	return this.Id
}
func (this *Model) IsNew() bool {
	return this.Id < 1
}
func (this *Model) GetBody() IBody {
	return this.Self().(IBody)
}
func (this *Model) GetData() interface{} {
	return this.self
}

// IBody 接口
func (this *Model) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*Model)
	data.Id = this.Id
	return nil
}

func (this *Model) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		log.Fatalln("model.Bind需要 self")
	}
	this.self = self
	err := Db.Sync(self)
	if err != nil {
		log.Printf("%T-->%v\n", self, err)
	}
}
