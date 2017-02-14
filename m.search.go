package gev

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type ISearchModel interface {
	IModel
	GetCondition() ISearch
	SearchSession(session *xorm.Session, condition ISearch) error
	Search(ISearch) (interface{}, error)
}

// M.search Entity
type SearchModel struct {
	Model `xorm:"extends"`
}

func (m *SearchModel) GetCondition() ISearch {
	return &SearchBody{}
}

func (s *SearchModel) SearchSession(session *xorm.Session, condition ISearch) error {
	search := condition.(*SearchBody)
	session.Where(search.Where).Cols(search.GetWhat())
	return nil
}

func (m *SearchModel) Search(condition ISearch) (interface{}, error) {
	bean := m.Self().(ISearchModel)
	data := make([]interface{}, condition.GetSize())
	n := 0
	session := Db.NewSession()
	defer session.Close()
	bean.SearchSession(session, condition)
	total, _ := session.Count(bean)
	bean.SearchSession(session, condition)
	session.Limit(condition.GetSize(), condition.GetBegin())
	err := session.Iterate(bean, func(i int, item interface{}) error {
		model := item.(IModel)
		model.SetSelf(model)
		data[i] = model.GetSearch()
		n++
		return nil
	})
	return &SearchData{data[:n], total}, err
}

func (m *SearchModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.Model.Bind(g, self)
	g.Info("搜索", "所有人都可以查询").Body(
		self.(ISearchModel).GetCondition(),
	).Data(
		NewSearchData(10, []interface{}{self.(ISearchModel).GetSearch()}),
	).POST("/search", func(c *gin.Context) {
		condition := m.Self().(ISearchModel).GetCondition()
		c.BindJSON(condition)
		data, err := m.New().(ISearchModel).Search(condition)
		Api(c, data, err)
	})
}
func NewSearchData(total int64, content []interface{}) map[string]interface{} {
	return map[string]interface{}{"content": content, "total": total}
}

type SearchData struct {
	Content []interface{} `json:"content" xorm:"not null"`
	Total   int64         `json:"total" xorm:"not null"`
}
