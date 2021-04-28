package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main(){
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Printf("Filed to start server",err.Error())
	}
}