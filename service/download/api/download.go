package api

import (
	"filestore-server/meta"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"filestore-server/common"
	dbplayer "filestore-server/db"
	"filestore-server/store/oss"
	// dlcfg "filestore-server/service/download/config"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找记录
	row, err := dbplayer.GetFileMeta(filehash)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
		return
	}

	//tblFile := dbcli.ToTableFile(dbResp.Data)
	//TODO
	u,_ := url.Parse(row.FileAddr.String)
	q := u.Query()
	u.RawQuery = q.Encode()
	//
	signedURL:=oss.DownloadURL(row.FileAddr.String,row.FileName.String)
	c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))

}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	//username := c.Request.FormValue("username")
	// TODO: 处理异常情况
	fResp := meta.GetFileMeta(fsha1)
	f,err:=os.Open(fResp.Location)
	//ufResp, uferr := dbplayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
		return
	}
	//r.ParseForm()
	defer f.Close()
	//小文件
	data,err:=ioutil.ReadAll(f)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,
		})
		return
	}
	c.Header("content-type", "application/octect-stream")
	c.Header("content-disposition", "attachment; filename=\""+fResp.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)
}
