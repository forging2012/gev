package gev

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ISearchModel interface {
	IModel
	GetCondition() interface{}
	Search(interface{}) (interface{}, error)
}

// M.search Entity
type SearchModel struct {
	Model `xorm:"extends"`
}

func (m *SearchModel) GetCondition() interface{} {
	return &SearchBody{}
}

func (m *SearchModel) Search(condition interface{}) (interface{}, error) {
	bean := m.Self().(ISearchModel)
	search, ok := condition.(*SearchBody)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%T.Condition格式不匹配"))
	}
	data := make([]interface{}, search.GetSize())
	n := 0
	session := Db.NewSession().Where(search.Where)
	total, _ := session.Count(bean)
	session.Cols(search.GetWhat()).Limit(search.GetSize(), search.GetBegin())
	err := session.Iterate(bean, func(i int, item interface{}) error {
		model := item.(IModel)
		model.SetSelf(model)
		data[i] = model.GetSearch()
		n++
		return nil
	})
	session.Close()
	return &SearchData{data[:n], total}, err
}

func (m *SearchModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.Model.Bind(g, self)
	g.Info("搜索", "").Body(
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

type SearchPage struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (s *SearchPage) GetSize() int {
	if s.Size < 1 {
		return 10
	}
	return s.Size
}

func (s *SearchPage) GetBegin() int {
	return s.Page * s.GetSize()
}

type SearchBody struct {
	SearchPage
	Where string `json:"where"`
	What  string `json:"what"`
}

func (s *SearchBody) GetWhat() string {
	if s.What == "" {
		return "*"
	}
	return s.What
}
