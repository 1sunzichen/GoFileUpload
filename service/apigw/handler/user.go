package handler

import (
	"context"
	"filestore-server/service/account/proto/proto"
	"filestore-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"net/http"
)
var(
	userCli proto.UserService
)
func init(){
  service:=micro.NewService()
  //初始化 解析命令行参数
  service.Init()
  //
	userCli=proto.NewUserService("go.micro.service.user",service.Client())
}



func DoSignupHandler(c *gin.Context)  {

	username:=c.Request.FormValue("username")
	passwd:=c.Request.FormValue("password")
	fmt.Printf("user:%s,pass:%s ",username,passwd)
	//存在一个之后 不在成功的问题 ，就是唯一索引值  没有插入，插入之后 就可以刷新
	resp,err:=userCli.Signup(context.TODO(),&proto.ReqSignup{
		Username: username,
		Password: passwd,
	})
	if err!=nil{
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":resp.Code,
		"msg":resp.Message,
	})
}

//处理用户注册请求
func SignupHandler(c *gin.Context){
	c.Redirect(http.StatusFound,"/static/view/signup.html")
}

func SigninHandler(c *gin.Context){
	username:=c.Request.FormValue("username")
	password:=c.Request.FormValue("password")
	resp,err:=userCli.Signin(context.TODO(),&proto.ReqSignin{
		Username: username,
		Password: password,
	})
	if err!=nil{
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":resp.Code,
		"token":resp.Token,
		"msg":resp.Message,
	})
}

func UserInfoHandler(c *gin.Context){
	username:=c.Request.FormValue("username")
	resp,err:=userCli.UserInfo(context.TODO(),&proto.ReqUserInfo{
		Username: username,
	})
	//user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		//w.WriteHeader(http.StatusForbidden)
		c.JSON(http.StatusForbidden,gin.H{})
		return
	}
	//4.组装并且响应用户数据
	cliResp:=util.RespMsg{
		Msg: "ok",
		Code: 0,
		Data: gin.H{
			"Username":username,
			"SignupAt":resp.SignupAt,
			"LastActive":resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK,"application/json",cliResp.JSONBytes())
}