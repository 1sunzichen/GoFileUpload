package handler

import (
	"encoding/json"
	"filestore-server/common"
	cfg "filestore-server/config"
	"filestore-server/meta"
	"filestore-server/mq"
	"filestore-server/store/oss"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	dbplayer "filestore-server/db"
)
func DoUploadHandler(c *gin.Context){
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
		resp:=util.RespMsg{
			Code:0,
			Msg:"ok",
			Data:struct{
				Location string
				Username string
				//Token string
			}{
				Location:"/static/view/home.html",
				Username:username,
				//Token:token,
			},
		}
		c.JSON(http.StatusOK,gin.H{
			"msg":"ok",
			"code":0,
			"data":resp,
		})
		//c.Data(http.StatusOK,"application",resp.JSONBytes())
		//w.Write(resp.JSONBytes())
		//http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}else{
		c.JSON(http.StatusOK,gin.H{
			"msg":"上传失败",
			"code":-1,

		})
		//w.Write([]byte("上传失败"))
	}

}
func UploadHandler(c *gin.Context){
	c.Redirect(http.StatusFound,"/static/view/upload.html")
}
//上传已完成
func UploadSucHandler(c *gin.Context){
	//io.WriteString(w,"Upload finished")
	c.JSON(http.StatusOK,gin.H{
		"msg":"上传成功",
		"code":0,
	})
}
//获取文件元信息
func GetFileMetaHandler(c *gin.Context){
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	//对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form。
	//解析结果中，POST或PUT请求主体要优先于URL查询字符串（同名变量，主体的值在查询字符串的值前面）。
	//r.ParseForm()
	filehash:=c.Request.FormValue("filehash")
	//fMeta:=meta.GetFileMeta(filehash)
	fMeta,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":"获取失败",
			"code":-1,
		})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d,err:=json.Marshal(fMeta)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,
		})

		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "application/json", d)
	//w.Write(d)
}
//批量获取文件元信息
func FileQueryHandler(c *gin.Context){
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	//对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form。
	//解析结果中，POST或PUT请求主体要优先于URL查询字符串（同名变量，主体的值在查询字符串的值前面）。
	//r.ParseForm()
	//字符串转数组
	limitCnt,_:=strconv.Atoi(c.Request.FormValue("limit"))
	username:=c.Request.FormValue("username");
	fmt.Println(limitCnt,username)

	fileMetas,err:=dbplayer.QueryUserFileMatas(username,limitCnt)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,

		})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:=json.Marshal(fileMetas);

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,
		})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
    //c.JSON(http.StatusOK,gin.H{
    //	"msg":"ok",
    //	"code":0,
    //	"data":data,
	//})
	c.Data(http.StatusOK, "application/json", data)
	//w.Write(data)
}

func DownloadHandler(c *gin.Context){
	//r.ParseForm()
	fsha1:=c.Request.FormValue("filehash")
	fm:=meta.GetFileMeta(fsha1)
	f,err:=os.Open(fm.Location)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,
		})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	//小文件
	data,err:=ioutil.ReadAll(f)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
			"code":-1,
		})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//http.
	//w.Header().Set("Content-Type","application/octect-stream")
	c.Header("content-type", "application/octect-stream")
	c.Header("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)

	//w.Header().Set("content-disposition","attachment;filename=\""+fm.FileName)
	//w.Write(data)
}
func FileMetaUpdateHandler(c *gin.Context){
	c.JSON(http.StatusMethodNotAllowed,gin.H{})
}
func DoFileMetaUpdateHandler(c *gin.Context){
	//r.ParseForm()
	opType:=c.Request.FormValue("op")
	fileSha1:=c.Request.FormValue("filehash")
	newFileName:=c.Request.FormValue("filename")
	if opType!="0"{
		c.JSON(http.StatusForbidden,gin.H{
			"msg":"更新失败",
			"code":-1,
		})
		//w.WriteHeader(http.StatusForbidden)
		return
	}


	curFileMeta:=meta.GetFileMeta(fileSha1)
	curFileMeta.FileName=newFileName
	meta.UpdateFileMeta(curFileMeta)
	data,err:=json.Marshal(curFileMeta)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"code":-1})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK,"application/json",data)
	//w.WriteHeader(http.StatusOK)
	//w.Write(data)

}
//删除接口
func FileDeleteHandler(c *gin.Context){
	//r.ParseForm()
	filesha1:=c.Request.FormValue("filehash")
	fmeta:=meta.GetFileMeta(filesha1)
	os.Remove(fmeta.Location)
	meta.RemoveFileMeta(filesha1)
    c.JSON(http.StatusOK,gin.H{})
	//w.WriteHeader(http.StatusOK)
}
//
func TryFastUploadHandler(c *gin.Context){
	//r.ParseForm()
	//
	username:=c.Request.FormValue("username")
	filehash:=c.Request.FormValue("filehash")
	filename:=c.Request.FormValue("filename")
	//filesize:=r.Form.Get("filesize")
	filesize,_:=strconv.ParseInt(c.Request.FormValue("filesize"),10,64)
	//
	fileMeta,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError,gin.H{})
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
   fileMetaStatus,_:=json.Marshal(fileMeta)
	if string(fileMetaStatus) == "{}"{
		//resp:=util.RespMsg{
		//	Code:-1,
		//	Msg: "秒传失败,请访问普通上传接口",
		//
		//
		//}
		//c.Data(http.StatusOK,"application/json",resp.JSONBytes())
		//w.Write(resp.JSONBytes())
		c.JSON(http.StatusForbidden,gin.H{
			"msg":"秒传失败,请访问普通上传接口",
			"code":-1,
		})
		return
	}

	suc:=dbplayer.OnUserFileUploadFinished(username,filehash,filename,filesize)


    if suc{
    	//resp:=util.RespMsg{
    	//	Code: 0,
    	//	Msg: "秒传成功",
		//
		//
		//}
		c.JSON(http.StatusOK,gin.H{
			"msg":"秒传成功",
			"code":0,
		})
		//w.Write(resp.JSONBytes())
	}else{
		//resp:=util.RespMsg{
		//	Code: -2,
		//	Msg: "秒传失败",
		//
		//
		//}
		//w.Write(resp.JSONBytes())
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":"秒传失败",
			"code":-2,
		})
	}

}
//生成下载地址接口
func DownloadURLHandler(c *gin.Context){
  //r.ParseForm()
  filehash:=c.Request.FormValue("filehash")
  //文件表查找记录
  row,_:=dbplayer.GetFileMeta(filehash)
  //TODO
	u,_ := url.Parse(row.FileAddr.String)
	q := u.Query()
	u.RawQuery = q.Encode()
 //
	sighedURL:=oss.DownloadURL(row.FileAddr.String,row.FileName.String)
	//c.Data(http.StatusOK,"application/json",[]byte(sighedURL))
    //w.Write([]byte(sighedURL))
    c.JSON(http.StatusOK,gin.H{
    	"data":sighedURL,
    	"msg":"ok",
    	"code":0,

	})

}
//BuildLifecycleRule