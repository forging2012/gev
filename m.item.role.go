package gev

import "github.com/gin-gonic/gin"

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
	if o.OwnerId == 0 {
		o.OwnerId = user.GetId()
		return true
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (m *ItemRoleModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	//
	m.SearchModel.Bind(g, self)
	g.Info("详情", "用户可以查看有读权限删除的东西").Params(
		g.PathParam("id", "id"),
	).Data(
		self.GetDetail(),
	).GET("/info/:id", func(c *gin.Context) {
		// 获取当前登录用户
		var data interface{}
		var err error
		if user, ok := NeedAuth(c); ok {
			data, err = m.New().(IItemModel).GetInfo(user.(IUserModel), c.Param("id"))
		}
		Api(c, data, err)
	})
	g.Info("添加/修改", "用户可以添加或修改有写权限的东西").Body(
		self.GetBody(),
	).Data(
		self.GetDetail(),
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
		}
		Api(c, data, err)
	})
	g.Info("删除", "用户可以删除有写权限的东西").Params(
		g.PathParam("id", "id"),
	).GET("/del/:id", func(c *gin.Context) {
		// 获取当前登录用户
		var err error
		if user, ok := NeedAuth(c); ok {
			err = m.New().(IItemModel).Delete(user.(IUserModel), c.Param("id"))
		}
		Api(c, nil, err)
	})
}
