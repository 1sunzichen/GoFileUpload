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
)
type MuitpartUploadInfo struct{
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}
func InitialMultipartUploadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	username:=r.Form.Get("username")
	filehash:=r.Form.Get("filehash")
	filesize,err:=strconv.Atoi(r.Form.Get("filesize"))
	if err!=nil{
		w.Write(util.NewRespMsg(-1,"params invaild",nil).JSONBytes())
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
	w.Write(util.NewRespMsg(0,"ok",upinfo).JSONBytes())
}


    // CompleteUploadHandler notifies to merge upload chunks
func CompleteUploadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	uploadid:=r.Form.Get("uploadid")
	username:=r.Form.Get("username")
	filehash:=r.Form.Get("filehash")
	filesize:=r.Form.Get("filesize")
	filename:=r.Form.Get("filename")
	totalCount:=0
	chunkCount:=0
	rConn:=rPool.RedisPool().Get()
	defer rConn.Close()
	data,err:=redis.Values(rConn.Do("HGETALL","MP_"+uploadid))
     if err!=nil{
     	w.Write(util.NewRespMsg(-1,"complete upload failed",nil).JSONBytes())
     	return
     }


	 for i:=0;i<len(data);i+=2{
	 	k:=string(data[i].([]byte))
	 	v:=string(data[i+1].([]byte))
	 	if k=="chunkcount"{
	 		totalCount,_=strconv.Atoi(v)// expected value
	 	}else if strings.HasPrefix(k,"chkidx_")&&v=="1"{
	 		chunkCount++// actual value
	 	}
	 }
	 if totalCount!=chunkCount{
	 	w.Write(util.NewRespMsg(-2,"invaild request",nil).JSONBytes())
	 	return
	 }
	// TODO Merge chunks
    fsize,_:=strconv.Atoi(filesize)
    dblayer.OnFileUploadFinished(filehash,filename,int64(fsize),"")
    dblayer.OnUserFileUploadFinished(username,filehash,filename,int64(fsize))
    w.Write(util.NewRespMsg(0,"ok",nil).JSONBytes())


}

func UploadPartHandler(w http.ResponseWriter,r *http.Request)  {
	// 1. Parse user request parameters
	r.ParseForm()
	 //username:=r.Form.Get("username")
	 uploadID:=r.Form.Get("uploadid")
	chunkIndex:=r.Form.Get("index")
	// 2. Get a connection from the Redis pool
	rConn:=rPool.RedisPool().Get()
	defer rConn.Close()
	//3.
	fpath:="/data/"+uploadID+"/"+chunkIndex
	os.MkdirAll(path.Dir(fpath),0744)
	fd,err:=os.Create(fpath)
	if err!=nil{
		w.Write(util.NewRespMsg(-1,"Upload part failed",nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf:=make([]byte,1024*1024)
    for{
    	n,err:=r.Body.Read(buf)
    	fd.Write(buf[:n])
    	if err!=nil{
    		break
		}

	}

	rConn.Do("HSET","MP_"+uploadID,"chkidx_"+chunkIndex,1)
	w.Write(util.NewRespMsg(0,"ok",nil).JSONBytes() )
}