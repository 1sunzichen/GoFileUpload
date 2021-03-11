package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
            newFile,err:=os.Create("/tmp/"+head.Filename)
            if err!=nil{
            	fmt.Printf("Failed to create file,err%s\n",err.Error())
            	return
			}
			defer newFile.Close()
            _,err=io.Copy(newFile,file)
            if err!=nil{
            	fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
            	return
			}
			http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
		}
}
//上传已完成
func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished")
}