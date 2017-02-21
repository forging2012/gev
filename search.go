package gev

import (
	"strings"

	"github.com/go-xorm/xorm"
)

type ISearch interface {
	GetBegin() int
	GetSize() int
	GetOrder(session *xorm.Session)
	GetOrderDefault(session *xorm.Session, default_order string)
	MakeSession(session *xorm.Session)
}

// 分页查询
type SearchPage struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	OrderBy string `json:"order_by,omitempty" gev:"排序规则:-id"`
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

func (s *SearchPage) GetOrder(session *xorm.Session) {
	if s.OrderBy != "" {
		orders := strings.Split(s.OrderBy, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	}
}

func (s *SearchPage) GetOrderDefault(session *xorm.Session, default_order string) {
	if s.OrderBy != "" {
		default_order = s.OrderBy
	}
	if default_order != "" {
		orders := strings.Split(default_order, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	}
}

func (s *SearchPage) MakeSession(session *xorm.Session) {
}

// 通用查询
type SearchBody struct {
	SearchPage
	Where string `json:"where" gev:"要查的内容 id,name,telphone"`
	What  string `json:"what" gev:"查询条件 name='abc' and telphone='xxx'"`
}

func (s *SearchBody) GetWhat() string {
	if s.What == "" {
		return "*"
	}
	return s.What
}

func (s *SearchBody) MakeSession(session *xorm.Session) {
	session.Where(s.Where).Cols(s.GetWhat())
}

type SearchKeyword struct {
	SearchPage
	Keyword string `json:"keyword,omitempty" gev:"关键词"`
}

func (s *SearchKeyword) GetWordLike() string {
	return WordLike(s.Keyword)
}

func (s *SearchKeyword) GetCharLike() string {
	return CharLike(s.Keyword)
}

// 查找地地
type SearchAddress struct {
	SearchPage
	Keyword  string `json:"keyword,omitempty" gev:"地地 如:'北京|'"`
	ParentId int    `json:"parent_id,omitempty" gev:"父地址id"`
}

func WordLike(key string) string {
	return strings.Join([]string{"%", key, "%"}, "")
}

func CharLike(key string) string {
	s := strings.Split(key, "")
	return WordLike(strings.Join(s, "%"))
}
