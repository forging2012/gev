package gev

import (
	"strings"
)

type ISearch interface {
	GetBegin() int
	GetSize() int
}

// 分页查询
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

// 通用查询
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
