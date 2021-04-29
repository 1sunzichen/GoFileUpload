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

     dbplayer "filestore-server/db"
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
            	Location: "static/file/"+head.Filename,
            	UploadAt: time.Now().Format("2016-01-02"),

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
            //meta.UpdateFileMeta(fileMeta)
            _=meta.UpdateFileMetaDb(fileMeta)
            //更新用户表文件表记录
            r.ParseForm()
            username:=r.Form.Get("username")
            suc:=dbplayer.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
            if suc{
				http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
			}else{
				w.Write([]byte("上传失败"))
			}
		}
}
//上传已完成
func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished")
}
//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	//对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form。
	//解析结果中，POST或PUT请求主体要优先于URL查询字符串（同名变量，主体的值在查询字符串的值前面）。
	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	//fMeta:=meta.GetFileMeta(filehash)
	fMeta,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(d)
}


func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1:=r.Form.Get("filehash")
	fm:=meta.GetFileMeta(fsha1)
	f,err:=os.Open(fm.Location)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	//小文件
	data,err:=ioutil.ReadAll(f)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//http.
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-disposition","attachment;filename=\""+fm.FileName)
	w.Write(data)
}
func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	opType:=r.Form.Get("op")
	fileSha1:=r.Form.Get("filehash")
	newFileName:=r.Form.Get("filename")
	if opType!="0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method!="POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta:=meta.GetFileMeta(fileSha1)
	curFileMeta.FileName=newFileName
	meta.UpdateFileMeta(curFileMeta)
	data,err:=json.Marshal(curFileMeta)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
//删除接口
func FileDeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filesha1:=r.Form.Get("filehash")
	fmeta:=meta.GetFileMeta(filesha1)
	os.Remove(fmeta.Location)
	meta.RemoveFileMeta(filesha1)

	w.WriteHeader(http.StatusOK)
}