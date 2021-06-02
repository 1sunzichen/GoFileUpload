package handler

import (
	"filestore-server/util"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	rPool "filestore-server/cache/redis"
	"strings"
	"time"
	dblayer "filestore-server/db"
	"github.com/gin-gonic/gin"
)
type MuitpartUploadInfo struct{
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}
func InitialMultipartUploadHandler(c *gin.Context){
	//r.ParseForm()
	username:=c.Request.FormValue("username")
	filehash:=c.Request.FormValue("filehash")
	filesize,err:=strconv.Atoi(c.Request.FormValue("filesize"))
	if err!=nil{
		//w.Write(util.NewRespMsg(-1,"params invaild",nil).JSONBytes())
		c.JSON(http.StatusOK,util.NewRespMsg(-1,"params invaild",nil))
		return
	}
	rConn:=rPool.RedisPool().Get()
	defer rConn.Close()
	upinfo:=MuitpartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		UploadID:username+fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize:5*1024*1024,
		ChunkCount:int(math.Ceil(float64(filesize)/5*1024*1024)),
	}

	rConn.Do("HSET","MP_"+upinfo.UploadID,"chunkcount",upinfo.ChunkCount)
	rConn.Do("HSET","MP_"+upinfo.UploadID,"filehash",upinfo.FileHash)
	rConn.Do("HSET","MP_"+upinfo.UploadID,"filesize",upinfo.FileSize)


	fmt.Println(upinfo)
	c.JSON(http.StatusOK,util.NewRespMsg(0,"ok",upinfo))
	//w.Write(util.NewRespMsg(0,"ok",upinfo).JSONBytes())
}


    //CompleteUploadHandler(w http.Response) 通知上传合并
func CompleteUploadHandler(c *gin.Context){
	//r.ParseForm()
	uploadid:=c.Request.FormValue("uploadid")
	username:=c.Request.FormValue("username")
	filehash:=c.Request.FormValue("filehash")
	filesize:=c.Request.FormValue("filesize")
	filename:=c.Request.FormValue("filename")
	totalCount:=0
	chunkCount:=0
	rConn:=rPool.RedisPool().Get()
	defer rConn.Close()
	data,err:=redis.Values(rConn.Do("HGETALL","MP_"+uploadid))
     if err!=nil{
     	c.JSON(http.StatusOK,util.NewRespMsg(-1,"complete upload failed",nil))
     	//w.Write(util.NewRespMsg(-1,"complete upload failed",nil).JSONBytes())
     	return
     }


	 for i:=0;i<len(data);i+=2{
	 	k:=string(data[i].([]byte))
	 	v:=string(data[i+1].([]byte))
	 	if k=="chunkcount"{
	 		totalCount,_=strconv.Atoi(v)//理论值
	 	}else if strings.HasPrefix(k,"chkidx_")&&v=="1"{
	 		chunkCount++//实际值
	 	}
	 }
	 if totalCount!=chunkCount{
	 	c.JSON(http.StatusOK,util.NewRespMsg(-2,"invaild request",nil))
	 	//w.Write(util.NewRespMsg(-2,"invaild request",nil).JSONBytes())
	 	return
	 }
	 //TODO 合并分块
    fsize,_:=strconv.Atoi(filesize)
    dblayer.OnFileUploadFinished(filehash,filename,int64(fsize),"")
    dblayer.OnUserFileUploadFinished(username,filehash,filename,int64(fsize))
    c.JSON(http.StatusOK,util.NewRespMsg(0,"ok",nil))
    //w.Write(util.NewRespMsg(0,"ok",nil).JSONBytes())


}

func UploadPartHandler(c *gin.Context)  {
	//1.解析用户请求参数
	//r.ParseForm()
	 //username:=r.Form.Get("username")
	 uploadID:=c.Request.FormValue("uploadid")
	chunkIndex:=c.Request.FormValue("index")
	//2.获得redis 链接池中的一个链接
	rConn:=rPool.RedisPool().Get()
	defer rConn.Close()
	//3.
	fpath:="/data/"+uploadID+"/"+chunkIndex
	os.MkdirAll(path.Dir(fpath),0744)
	fd,err:=os.Create(fpath)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":"Upload part failed",
			"code":-1,
		})
		//w.Write(util.NewRespMsg(-1,"Upload part failed",nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf:=make([]byte,1024*1024)
    for{
    	n,err:=c.Request.Body.Read(buf)
    	fd.Write(buf[:n])
    	if err!=nil{
    		break
		}

	}

	rConn.Do("HSET","MP_"+uploadID,"chkidx_"+chunkIndex,1)
	//w.Write(util.NewRespMsg(0,"ok",nil).JSONBytes() )
	c.JSON(http.StatusOK,gin.H{
		"msg":"ok",
		"code":0,
	})
}