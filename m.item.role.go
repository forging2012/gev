package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// 自己和管理员可以编辑
type IItemRoleModel interface {
	IItemOwnerModel
}

type ItemRoleModel struct {
	ItemOwnerModel `xorm:"extends"`
}

func (o *ItemRoleModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (o *ItemRoleModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
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

func (i *ItemRoleModel) DeleteIds(user IUserRoleModel, ids []int) error {
	if len(ids) < 1 {
		return errors.New("数组长度不能为0")
	}
	if user.IsAdmin() {
		_, err := Db.In("id", ids).Delete(i.Self())
		return err
	} else {
		_, err := Db.In("id", ids).Where("owner_id=?", user.GetId()).Delete(i.Self())
		return err
	}
}

func (m *ItemRoleModel) Bind(g ISwagRouter, self IModel) {
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
		if user, ok := NeedAuth(c); ok {
			data, err = m.New().(IItemModel).GetInfo(user.(IUserModel), c.Param("id"))
			Api(c, data, err)
		}
	})
	g.Info("添加/修改", "用户可以添加或修改有写权限的东西").Body(
		self.GetBody(),
	).Data(
		self,
	).POST("/save", func(c *gin.Context) {
		// 获取当前登录用户
		var data interface{}
		var err error
		model := m.New().(IItemModel)
		src := model.GetBody()
		if err := c.BindJSON(src); err != nil {
			Err(c, 1, err)
			return
		}
		if user, ok := NeedAuth(c); ok {
			data, err = model.Save(user.(IUserModel), src)
			Api(c, data, err)
		}
	})
	g.Info("删除", "用户可以删除有写权限的东西").Params(
		g.PathParam("id", "id"),
	).GET("/del/:id", func(c *gin.Context) {
		// 获取当前登录用户
		var err error
		if user, ok := NeedAuth(c); ok {
			err = m.New().(IItemModel).Delete(user.(IUserModel), c.Param("id"))
			Api(c, nil, err)
		}
	})
	g.Info("批量删除", "用户可以批量删除有写权限的东西").Body(
		[]int{},
	).POST("/del", func(c *gin.Context) {
		// 获取当前登录用户
		var ids []int
		err := c.BindJSON(&ids)
		if err != nil {
			Err(c, 1, errors.New("需要id数组"))
			return
		}
		if user, ok := NeedAuth(c); ok {
			err = m.New().(IItemRoleModel).DeleteIds(user.(IUserRoleModel), ids)
			Api(c, nil, err)
		}
	})
}
