package handler

import (
	"filestore-server/util"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HTTPInterceptor() gin.HandlerFunc{
	return func(c *gin.Context){
		 //r.ParseForm()
		 username:=c.Request.FormValue("username")
		 token:=c.Request.FormValue("token")
		 if len(username)<3||!IsTokenVaild(token){
		 	//w.WriteHeader(http.StatusForbidden)
		 	c.Abort()
		 	resp:=util.NewRespMsg(
		 		999,
		 		"token无效",
		 		nil)
		 	c.JSON(http.StatusOK,resp)

		 	return
		 }
		 c.Next()
		}
}
