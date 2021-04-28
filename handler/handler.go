package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
		if r.Method=="GET"{
			//返回上传html页面
			data,err:=ioutil.ReadFile("./static/view/upload.html")
			if err !=nil{
				io.WriteString(w,"interError")
				return
			}
			io.WriteString(w,string(data))
		}else if r.Method=="POST"{
            file,head,err:=r.FormFile("file")
            if err!=nil{
				fmt.Printf("failed to get data,err:%s\n",err)
				return
			}
            defer file.Close()
            fileMeta:=meta.FileMeta{
            	FileName: head.Filename,
            	Location:"./upload/"+head.Filename,
            	UploadAt: time.Now().Format("2021-06-02 15:04:05"),

			}
            newFile,err:=os.Create(fileMeta.Location)
            if err!=nil{
            	fmt.Printf("Failed to create file,err%s\n",err.Error())
            	return
			}
			defer newFile.Close()
            fileMeta.FileSize,err=io.Copy(newFile,file)
            if err!=nil{
            	fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
            	return
			}
			newFile.Seek(0,0)
            fileMeta.FileSha1=util.FileSha1(newFile)
            meta.UpdateFileMeta(fileMeta)
			http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
		}
}
//上传已完成
func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished")
}
// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	fMeta:=meta.GetFileMeta(filehash);
	data,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//下载信息
func DownloadHandler(w http.ResponseWriter,r *http.Request)  {
	
}