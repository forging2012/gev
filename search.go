package gev

import (
	"strings"

	"github.com/go-xorm/xorm"
)

type ISearch interface {
	GetBegin() int
	GetSize() int
	SetDefaultOrder(o string)
	GetOrder(session *xorm.Session)
	GetOrderDefault(session *xorm.Session, default_order string)
	MakeSession(user IUserModel, session *xorm.Session)
}

// 分页查询
type SearchPage struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	OrderBy string `json:"order_by,omitempty" gev:"排序规则:-id"`
}

func (this *SearchPage) GetSize() int {
	if this.Size < 1 {
		return 10
	}
	return this.Size
}

func (this *SearchPage) GetBegin() int {
	return this.Page * this.GetSize()
}

func (this *SearchPage) SetDefaultOrder(o string) {
	if this.OrderBy == "" {
		this.OrderBy = o
	}
}

func (this *SearchPage) GetOrder(session *xorm.Session) {
	if this.OrderBy != "" {
		orders := strings.Split(this.OrderBy, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	} else {
		session.Desc("id")
	}
}

func (this *SearchPage) GetOrderDefault(session *xorm.Session, default_order string) {
	if this.OrderBy != "" {
		default_order = this.OrderBy
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

func (this *SearchPage) MakeSession(user IUserModel, session *xorm.Session) {
}

// 通用查询
type SearchBody struct {
	SearchPage
	Where string `json:"where" gev:"要查的内容 id,name,telphone"`
	What  string `json:"what" gev:"查询条件 name='abc' and telphone='xxx'"`
}

func (this *SearchBody) GetWhat() string {
	if this.What == "" {
		return "*"
	}
	return this.What
}

func (this *SearchBody) MakeSession(user IUserModel, session *xorm.Session) {
	session.Where(this.Where).Cols(this.GetWhat())
}

type SearchKeyword struct {
	SearchPage
	Keyword string `json:"keyword,omitempty" gev:"关键词"`
}

func (this *SearchKeyword) GetWordLike() string {
	return WordLike(this.Keyword)
}

func (this *SearchKeyword) GetCharLike() string {
	return CharLike(this.Keyword)
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
