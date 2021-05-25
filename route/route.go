package route

import (
	"filestore-server/handler"
	"github.com/gin-gonic/gin"
	//"net/http"
)
func Route() *gin.Engine{
   router:=gin.Default()
   router.Static("/static", "./static")
   //router.StaticFS("/static2", http.Dir("my_file_system"))
   //router.StaticFile("/favicon.ico", "./resources/favicon.ico")

   router.GET("/user/signup",handler.SignupHandler)
   router.POST("/user/signup",handler.DoSignupHandler)

	//http.HandleFunc("/user/signup",handler.SignupHandler)
	router.GET("/user/signin",handler.SignInHandler)
	router.POST("/user/signin",handler.DoSignInHandler)
    router.Use(handler.HTTPInterceptor())

	//分块上传
	router.POST("/file/mpupload/init",handler.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart",handler.UploadPartHandler)
	router.POST("/file/mpupload/complete",handler.CompleteUploadHandler)

	//http.Handle("/", http.FileServer(http.Dir("/static")))
	router.POST("/file/uploadprocess",handler.DoUploadHandler)
	router.GET("/file/uploadprocess",handler.UploadHandler)
	router.GET("/file/uploadprocess/suc",handler.UploadSucHandler)
	router.POST("/file/download",handler.DownloadHandler)
	router.POST("/file/meta",handler.GetFileMetaHandler)
	router.POST("/file/query",handler.FileQueryHandler)

	router.POST("/file/update",handler.DoFileMetaUpdateHandler)
	router.GET("/file/update",handler.FileMetaUpdateHandler)
	router.POST("/file/delete",handler.FileDeleteHandler)
	router.POST("/file/fastupload",handler.TryFastUploadHandler)

	router.POST("/user/info",handler.UserInfoHandler)
	router.POST("/file/downloadurl",handler.DownloadURLHandler)
   return router
}