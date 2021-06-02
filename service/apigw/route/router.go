package route

import (
	"filestore-server/service/apigw/handler"
	"github.com/gin-gonic/gin"
)

func Router()* gin.Engine {
	router:=gin.Default()
	router.Static("/static/","./static")
	router.GET("/user/signup",handler.SignupHandler)
	router.POST("/user/signup",handler.DoSignupHandler)
	router.POST("/user/signin",handler.SigninHandler)
	router.Use(handler.Authorize())
	router.POST("/user/info",handler.UserInfoHandler)
	// 用户文件查询
	router.POST("/file/query", handler.FileQueryHandler)
	return router
}
