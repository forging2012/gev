package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type IItemModel interface {
	ISearchModel
	CanRead(user IUserModel) bool
	CanWrite(user IUserModel) bool
	GetInfo(user IUserModel, id string) (interface{}, error)
	Save(user IUserModel, schema ISchemaBody) (interface{}, error)
	Delete(user IUserModel, id string) error
}

type ItemModel struct {
	SearchModel `xorm:"extends"`
}

func (o *ItemModel) CanRead(user IUserModel) bool {
	return true
}
func (o *ItemModel) CanWrite(user IUserModel) bool {
	return true
}

func (m *ItemModel) GetInfo(user IUserModel, id string) (interface{}, error) {
	bean := m.Self().(IItemModel)
	ok, err := Db.Id(id).Get(bean)
	if !ok {
		return nil, errors.New("不存在")
	}
	if !bean.CanRead(user) {
		return nil, errors.New("没有权限")
	}
	return bean.GetDetail(), err
}
func (m *ItemModel) Save(user IUserModel, schema ISchemaBody) (interface{}, error) {
	bean, err := schema.GetData()
	if err != nil {
		return nil, err
	}
	// 更新或插入
	if !bean.(IItemModel).CanWrite(user) {
		return nil, errors.New("没有权限")
	}
	if bean.GetId() > 0 {
		item := bean.New()
		var ok bool
		ok, err = Db.Id(bean.GetId()).Get(item)
		if !ok {
			return nil, errors.New("不存在")
		}
		if err != nil {
			return nil, err
		}
		if !item.(IItemModel).CanWrite(user) {
			return nil, errors.New("没有修改权限")
		}
		_, err = Db.ID(bean.GetId()).Update(bean)
	} else {
		_, err = Db.InsertOne(bean)
	}
	return bean, err
}
func (m *ItemModel) Delete(user IUserModel, id string) error {
	bean := m.Self()
	ok, err := Db.Id(id).Get(bean)
	if !ok {
		return errors.New("不存在")
	}
	if err != nil {
		return err
	}
	if !bean.(IItemModel).CanWrite(user) {
		return errors.New("没有权限")
	}
	_, err = Db.ID(m.Id).Delete(bean)
	return err
}

func (m *ItemModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.SearchModel.Bind(g, self)
	g.Info("详情", "用户可以查看有读权限删除的东西").Params(
		g.PathParam("id", "id"),
	).Data(
		self.GetDetail(),
	).GET("/info/:id", func(c *gin.Context) {
		// 获取当前登录用户
		user, _ := c.Get("user")
		data, err := m.New().(IItemModel).GetInfo(user.(IUserModel), c.Param("id"))
		Api(c, data, err)
	})
	g.Info("添加/修改", "用户可以添加或修改有写权限的东西").Body(
		self.GetBody(),
	).Data(
		self.GetDetail(),
	).POST("/save", func(c *gin.Context) {
		// 获取当前登录用户
		user, _ := c.Get("user")
		model := m.New().(IItemModel)
		src := model.GetBody()
		if err := c.BindJSON(src); err != nil {
			Err(c, 1, err)
			return
		}
		data, err := model.Save(user.(IUserModel), src)
		Api(c, data, err)
	})
	g.Info("删除", "用户可以删除有写权限的东西").Params(
		g.PathParam("id", "id"),
	).GET("/del/:id", func(c *gin.Context) {
		// 获取当前登录用户
		user, _ := c.Get("user")
		err := m.New().(IItemModel).Delete(user.(IUserModel), c.Param("id"))
		Api(c, nil, err)
	})
}
