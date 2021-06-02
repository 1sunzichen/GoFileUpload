package api

import (
	"encoding/json"
	"filestore-server/common"
	cfg "filestore-server/config"
	dbplayer "filestore-server/db"
	"filestore-server/store/oss"
	"filestore-server/meta"
	"filestore-server/mq"
	"filestore-server/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
	"fmt"
	"log"

)

func DoUploadHandler(c *gin.Context){
	errCode:=0

	defer func(){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Methods","POST,OPTIONS")
		if errCode<0{
			c.JSON(http.StatusOK,gin.H{
				"code":errCode,
				"msg":"上传失败",
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"code":errCode,
				"msg":"上传成功",
			})
		}
	}()
	file,head,err:=c.Request.FormFile("file")
	if err!=nil{
		fmt.Printf("failed to get data,err:%s\n",err)
		return
	}
	defer file.Close()


	fileMeta:=meta.FileMeta{
		FileName: head.Filename,
		Location: "static/file/"+head.Filename,
		UploadAt: time.Now().Format("2016-01-02"),

	}
	//fileMeta.Location=cfg.TempLocalRootDir+fileMeta.FileSha1
	newFile,err:=os.Create(fileMeta.Location)
	if err!=nil{
		fmt.Printf("Failed to create file,err%s\n",err.Error())
		return
	}
	cur_offset,_:=newFile.Seek(0,1)
	//fmt.Printf("相对位置1:%d\n",cur_offset)
	defer newFile.Close()
	//cur_offset,_=newFile.Seek(0,1)
	//fmt.Printf("相对位置2:%d\n",cur_offset)
	fileMeta.FileSize,err=io.Copy(newFile,file)
	cur_offset,_=newFile.Seek(0,1)
	fmt.Printf("相对位置3:%d\n",cur_offset)
	if err!=nil{
		fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
		return
	}
	newFile.Seek(0,0)
	fileMeta.FileSha1=util.FileSha1(newFile)

	//以前是上传到本地   会更新
	newFile.Seek(0,0)
	//meta.UpdateFileMeta(fileMeta)
	//newFile2,err:=os.Open(fileMeta.Location)
	ossPath:="oss/"+fileMeta.FileSha1
	if !cfg.AsyncTransferEnable {
		// TODO: 设置oss中的文件名，方便指定文件名下载
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			log.Println(err.Error())
			//errCode = -5
			return
		}
		fileMeta.Location = ossPath
	} else {
		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   fileMeta.Location,
			DestLocation:  ossPath,
			DestStoreType: common.StoreOSS,
		}
		pubData, _ := json.Marshal(data)
		suc := mq.Publish(
			cfg.TransExchangeName,
			cfg.TransOSSRoutingKey,
			pubData)
		if !suc {
			//todo
			c.JSON(-1,
				"")
			//w.Write([]byte("mq更新失败"))
		}
	}
	//
	fmt.Println(fileMeta.Location)
	//更新文件信息
	_=meta.UpdateFileMetaDb(fileMeta)
	//更新用户表文件表记录

	username:=c.Request.FormValue("username")
	suc:=dbplayer.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
	if suc{
		errCode=0
	}else{

		errCode=-6

	}

}