package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)
//HandlerFunc type是一个适配器，通过类型转换让我们可以将普通的函数作为HTTP处理器使用。
//如果f是一个具有适当签名的函数，HandlerFunc(f)通过调用f实现了Handler接口。

func main(){
	//FileServer返回一个使用FileSystem接口root提供文件访问服务的HTTP处理器 "/"。
	//要使用操作系统的FileSystem接口实现，可使用http.Dir：
	fs := http.FileServer(http.Dir("static/"))
	//StripPrefix返回一个处理器，该处理器会将请求的URL.Path字段中给定前缀prefix去除后再交由h处理。
	//StripPrefix会向URL.Path字段中没有给定前缀的请求回复404 page not found。
	http.Handle("/static/", http.StripPrefix("/static", fs))

	//http.Handle("/", http.FileServer(http.Dir("/static")))
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Printf("Filed to start server",err.Error())
	}

}