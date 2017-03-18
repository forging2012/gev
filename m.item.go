package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type IItemModel interface {
	ISearchModel
	CanRead(user IUserModel) bool
	CanWrite(user IUserModel) bool
	GetInfo(user IUserModel, id interface{}) (interface{}, error)
	Save(user IUserModel, schema IBody) (interface{}, error)
	Delete(user IUserModel, id interface{}) error
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

func (m *ItemModel) GetInfo(user IUserModel, id interface{}) (interface{}, error) {
	data := m.Self().(IModel).GetData()
	ok, err := Db.Table(m.Self()).Where("id=?", id).Get(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("不存在")
	}
	if bean, ok := data.(IItemModel); ok && !bean.CanRead(user) {
		return nil, errors.New("没有权限")
	}
	if bean, ok := data.(IDataDetail); ok {
		return bean.GetDetail(user), err
	}
	return data, nil
}
func (m *ItemModel) Save(user IUserModel, schema IBody) (interface{}, error) {
	bean := m.New().(IModel)
	var err error
	if err = schema.CopyTo(user, bean); err != nil {
		return nil, err
	}
	// 更新或插入
	if bean.IsNew() {
		if !bean.(IItemModel).CanWrite(user) {
			return nil, errors.New("没有权限")
		}
		_, err = Db.InsertOne(bean)
	} else {
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
	}
	return m.self.(IItemModel).GetInfo(user, m.Id)
}
func (m *ItemModel) Delete(user IUserModel, id interface{}) error {
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
	_, err = Db.ID(id).Delete(bean.New())
	return err
}

func (m *ItemModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.SearchModel.Bind(g, self)
	g.Info("详情", "用户可以查看有读权限的东西").Params(
		g.PathParam("id", "id"),
	).Data(
		self.GetData(),
	).GET("/info/:id", func(c *gin.Context) {
		// 获取当前登录用户
		var data interface{}
		var err error
		if user, ok := c.Get("user"); ok {
			data, err = m.New().(IItemModel).GetInfo(user.(IUserModel), c.Param("id"))
		} else {
			data, err = m.New().(IItemModel).GetInfo(nil, c.Param("id"))
		}
		Api(c, data, err)
	})
	g.Info("添加/修改", "用户可以添加或修改有写权限的东西").Body(
		self.GetBody(),
	).Data(
		self.GetBody(),
	).POST("/save", func(c *gin.Context) {
		// 获取当前登录用户
		var data interface{}
		var err error
		user, ok := c.Get("user")
		model := m.New().(IItemModel)
		src := model.GetBody()
		if err := c.BindJSON(src); err != nil {
			Err(c, 1, err)
			return
		}
		if ok {
			data, err = model.Save(user.(IUserModel), src)
		} else {
			data, err = model.Save(nil, src)
		}
		Api(c, data, err)
	})
	g.Info("删除", "用户可以删除有写权限的东西").Params(
		g.PathParam("id", "id"),
	).GET("/del/:id", func(c *gin.Context) {
		// 获取当前登录用户
		var err error
		if user, ok := c.Get("user"); ok {
			err = m.New().(IItemModel).Delete(user.(IUserModel), c.Param("id"))
		} else {
			err = m.New().(IItemModel).Delete(nil, c.Param("id"))
		}
		Api(c, nil, err)
	})
}
