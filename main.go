package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)
// HandlerFunc is an adapter type that allows us to use ordinary functions as HTTP handlers through type conversion.
// If f is a function with the appropriate signature, HandlerFunc(f) implements the Handler interface by calling f.

func main(){
	// FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at "/".
	// To use the operating system's file system implementation, use http.Dir:
	fs := http.FileServer(http.Dir("static/"))
	// StripPrefix returns a handler that removes the given prefix from the request URL.Path before forwarding to h.
	// StripPrefix returns a 404 page not found reply for requests whose URL.Path does not begin with prefix.
	http.Handle("/static/", http.StripPrefix("/static", fs))
	// Multipart upload
	http.HandleFunc("/file/mpupload/init",handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",handler.HTTPInterceptor(handler.CompleteUploadHandler))

	//http.Handle("/", http.FileServer(http.Dir("/static")))
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
<<<<<<< HEAD
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
=======
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/query",handler.FileQueryHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/user/signin",handler.SignInHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUploadHandler))
	http.HandleFunc("/user/signup",handler.SignupHandler)
	http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
>>>>>>> part5-2

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Printf("Filed to start server",err.Error())
	}

}