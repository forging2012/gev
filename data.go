package gev

import "time"

type IDataSearch interface {
	GetSearch(user IUserModel) interface{}
}
type IDataDetail interface {
	GetDetail(user IUserModel) interface{}
}

type IData interface {
	IDataDetail
	IDataSearch
}

type ModelData struct {
	Id       int       `json:"id,omitempty"`
	CreateAt time.Time `json:"create_at,omitempty"`
}

type UserModelData struct {
	Id       int    `json:"id,omitempty"`
	Nickname string `gev:"用户昵称" json:"nickname" xorm:""`
	Telphone string `gev:"电话号码" json:"telphone" xorm:"varchar(32) unique(telphone) not null"`
	Password string `gev:"密码" json:"-" xorm:""`
}

func (this *UserModelData) TableName() string {
	return "user"
}
