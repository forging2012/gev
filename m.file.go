package gev

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inu1255/gev/libs"
)

type IFileModel interface {
	IItemRoleModel
	Upload(filename string, file multipart.File, user IUserModel) (interface{}, error)
}

type FileModel struct {
	ItemRoleModel `xorm:"extends"`
	Ext           string `json:"ext,omitempty" xorm:"" gev:"文件后缀"`
	Place         string `json:"-" xorm:""`
	Filename      string `json:"-" xorm:"" gev:""`
	MD5           string `json:"-" xorm:"" gev:""`
	Url           string `json:"url" xorm:"" gev:"文件地址,需加上host,如http://www.tederen.com:8017/"`
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

func (f *FileModel) Upload(filename string, src multipart.File, user IUserModel) (interface{}, error) {
	var err error
	bean := f.Self().(IFileModel)
	// 创建用户文件夹
	uid := "0"
	if user != nil {
		uid = strconv.Itoa(user.GetId())
		f.OwnerId = user.GetId()
	}
	dir := strings.Join([]string{"upload", uid}, "/")
	err = os.MkdirAll(dir, 0755)

	bs, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write(bs)
	f.MD5 = hex.EncodeToString(h.Sum(nil))
	// 保存文件
	f.Place = strings.Join([]string{dir, "/", f.MD5}, "")
	if _, err = os.Stat(f.Place); err == nil {
		err = ioutil.WriteFile(f.Place, bs, 0644)
		if err != nil {
			return nil, err
		}
	}
	//  保存文件
	f.Ext = f.GetExt(filename)
	f.Filename = filename
	f.Url = f.Place
	if ok, _ := Db.Where("place=? and owner_id=?", f.Place, f.OwnerId).Get(bean); ok {
		return bean, nil
	}
	Db.InsertOne(bean)
	return bean, nil
}

func (m *FileModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.ItemRoleModel.Bind(g, self)
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
		data, err := m.New().(IFileModel).Upload(filename, file, user)
		Api(c, data, err)
	})
	g.Info("导出csv文件", "post一个二维数组").Body(
		[][]string{},
	).POST("/export/csv", func(c *gin.Context) {
		var tables [][]string
		if err := c.BindJSON(&tables); err != nil {
			Err(c, 0, err)
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=表格.csv")
		libs.SimpleWriteExcel(c.Writer, tables)
	})
}
