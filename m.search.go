package gev

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type ISearchModel interface {
	IModel
	GetCondition() ISearch
	Search(user IUserModel, condition ISearch) (interface{}, error)
}

// M.search Entity
type SearchModel struct {
	Model `xorm:"extends"`
}

func (this *SearchModel) GetCondition() ISearch {
	return &SearchBody{}
}

func GetSearchData2(bean interface{}, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	session := Db.NewSession()
	defer session.Close()
	sessionFunc(session)
	total, _ := session.Count(bean)
	sessionFunc(session)
	session.Limit(condition.GetSize(), condition.GetBegin())
	condition.GetOrderDefault(session, "-id")
	data := make([]interface{}, condition.GetSize())
	n := 0
	err := session.Iterate(bean, func(i int, item interface{}) error {
		data[i] = item
		n++
		return nil
	})
	return &SearchData{Content: data[:n], Total: total}, err
}

func GetSearchData(bean interface{}, user IUserModel, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	if data, ok := bean.(IDataSearch); ok {
		session := Db.NewSession()
		defer session.Close()
		sessionFunc(session)
		total, _ := session.Count(data)
		sessionFunc(session)
		session.Limit(condition.GetSize(), condition.GetBegin())
		condition.GetOrderDefault(session, "-id")
		data := make([]interface{}, condition.GetSize())
		n := 0
		err := session.Iterate(bean, func(i int, item interface{}) error {
			model := item.(IDataSearch)
			data[i] = model.GetSearch(user)
			n++
			return nil
		})
		return &SearchData{Content: data[:n], Total: total}, err
	}
	return GetSearchData2(bean, condition, sessionFunc)
}

func (this *SearchModel) Search(user IUserModel, condition ISearch) (interface{}, error) {
	bean := this.GetData()
	return GetSearchData(bean, user, condition, func(session *xorm.Session) {
		condition.MakeSession(user, session)
	})
}

func (this *SearchModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = this
	}
	this.Model.Bind(g, self)
	g.Info("搜索").Body(
		self.(ISearchModel).GetCondition(),
	).Data(
		NewSearchData(10, []interface{}{self.GetData()}),
	).POST("/search", func(c *gin.Context) {
		condition := this.Self().(ISearchModel).GetCondition()
		err := c.BindJSON(condition)
		if err != nil {
			Err(c, 1, err)
			return
		}
		if user, ok := c.Get("user"); ok {
			data, err := this.New().(ISearchModel).Search(user.(IUserModel), condition)
			Api(c, data, err)
		} else {
			data, err := this.New().(ISearchModel).Search(nil, condition)
			Api(c, data, err)
		}
	})
}
func NewSearchData(total int64, content []interface{}) map[string]interface{} {
	return map[string]interface{}{"content": content, "total": total}
}

type SearchData struct {
	Content []interface{} `json:"content" xorm:"" gev:"数据数组"`
	Total   int64         `json:"total" xorm:"" gev:"数据总量"`
	Ext     interface{}   `json:"ext,omitempty" xorm:"" gev:"附加数据"`
}
