package gev

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type IFileModel interface {
	IItemAdminModel
	Upload(filename string, file multipart.File, user IUserModel) (interface{}, error)
	GetUrl() string
}

type FileModel struct {
	ItemAdminModel `xorm:"extends"`
	Ext            string `gev:"文件后缀" json:"ext,omitempty" xorm:""`
	Place          string `json:"-" xorm:""`
	Url            string `gev:"文件地址" json:"url,omitempty" xorm:""`
}

func (f *FileModel) TableName() string {
	return "file"
}

func (f *FileModel) GetExt(filename string) string {
	index := strings.LastIndex(filename, ".")
	if index >= 0 {
		return strings.ToLower(filename[index+1:])
	}
	return ""
}

func (f *FileModel) GetUrl() string {
	// uri := ""
	// if len(f.Place) > 6 {
	// 	uri = f.Place[6:]
	// }
	return strings.Join([]string{"http://", Host, "/", f.Place}, "")
}

func (f *FileModel) Upload(filename string, src multipart.File, user IUserModel) (interface{}, error) {
	var err error
	bean := f.Self().(IFileModel)
	// 创建用户文件夹
	uid := "0"
	if user != nil {
		uid = strconv.Itoa(user.GetId())
		f.OwnerId = user.GetId()
	}
	now := time.Now()
	dir := strings.Join([]string{"upload", uid, now.Format("2006-01-02")}, "/")
	err = os.MkdirAll(dir, 0755)

	// 保存文件
	f.Place = strings.Join([]string{dir, "/", now.Format("03:04:05"), "-", filename}, "")
	dst, err := os.OpenFile(f.Place, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	//  保存文件
	f.Ext = f.GetExt(filename)
	f.Url = bean.GetUrl()
	// if ok, _ := Db.Where("place=? and owner_id=?", f.Place, f.OwnerId).Get(bean); ok {
	// 	return bean, nil
	// }
	Db.InsertOne(bean)
	return bean, nil
}

func (m *FileModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.ItemAdminModel.Bind(g, self)
	g.Info("上传文件", "").Data(self).Params(
		g.FileParam("file", "文件"),
	).POST("/upload", func(c *gin.Context) {
		var user IUserModel
		// 上传者
		if u, ok := c.Get("user"); ok {
			user = u.(IUserModel)
		}
		// 上传的文件
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			Err(c, 2, err)
			return
		}
		// 文件名
		var filename string
		if header.Filename != "" {
			filename = header.Filename
		} else {
			filename = strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		// 保存文件
		m.New().(IFileModel).Upload(filename, file, user)
	})
}
