package route

import (
	"filestore-server/service/upload/api"
	"github.com/gin-gonic/gin"
)


func Route() *gin.Engine{
	router:=gin.Default()
	router.Static("/static", "./static")
	//router.StaticFS("/static2", http.Dir("my_file_system"))

	//http.Handle("/", http.FileServer(http.Dir("/static")))
	router.POST("/file/upload",api.DoUploadHandler)
    router.OPTIONS("/file/upload", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Methods","POST,OPTIONS")
		c.Status(204)
	})
	return router
}