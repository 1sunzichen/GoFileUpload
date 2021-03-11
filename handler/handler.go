package handler

import (
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
            	Location: "./static/file/"+head.Filename,
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
            meta.UpdateFileMeta(fileMeta)
			http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
		}
}
//上传已完成
func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished")
}