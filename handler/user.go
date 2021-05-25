package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	//"io/ioutil"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

const (
	pwd_salt="*#890"
)

func DoSignupHandler(c *gin.Context)  {

	username:=c.Request.FormValue("username")
	passwd:=c.Request.FormValue("password")
	//存在一个之后 不在成功的问题 ，就是唯一索引值  没有插入，插入之后 就可以刷新
	phone:=c.Request.FormValue("phone")
	if len(username)<3||len(passwd)<5{
		c.JSON(http.StatusOK,gin.H{
			"msg":"无效参数",
			"code":-1,

		})
		return
		//w.Write([]byte("Invaild parameter"))
		//return
	}
	enc_passwd:=util.Sha1([]byte(passwd+pwd_salt))
	fmt.Println(enc_passwd,username)
	suc:=dblayer.UserSignUp(username,enc_passwd,phone)
	if suc{
		c.JSON(http.StatusOK,gin.H{
			"msg":"注册成功",
			"code":0,

		})
		//w.Write([]byte("success"))
	}else{
		c.JSON(http.StatusOK,gin.H{
			"msg":"注册失败",
			"code":-2,

		})
		//w.Write([]byte("Fail"))
	}
}
//处理用户注册请求
func SignupHandler(c *gin.Context){
	c.Redirect(http.StatusFound,"/static/view/signup.html")
}
type result struct {
	url string
}
func SignInHandler(c *gin.Context){
	c.Redirect(http.StatusFound,"/static/view/signin.html")
}
//DoSignInHandler post登录接口
func DoSignInHandler(c *gin.Context){
	//r.ParseForm()
	//校验用户名及密码
	username:=c.Request.FormValue("username")
	password:=c.Request.FormValue("password")
	encPasswd:=util.Sha1([]byte(password+pwd_salt))
	pwdChecked:=dblayer.UserSignin(username,encPasswd)
	if !pwdChecked{
		c.JSON(http.StatusOK,gin.H{
			"msg":"登录失败",
			"code":-2,

		})
		//w.Write([]byte("FAILED"))
		return
	}
	//生成访问凭证
	token:=GenToken(username)
	upRes:=dblayer.UpdateToken(username,token)
	if !upRes{
		//w.Write([]byte("FAILED"))
		c.JSON(http.StatusOK,gin.H{
			"msg":"登录失败",
			"code":-2,

		})
		return
	}
	//登录成功后 重定向到首页
	//url:=result{url:"http://"+r.Host+"/static/view/home.html"}
	//http.Redirect(w,r,"http://"+r.Host+"/static/view/home.html",http.StatusFound)
	//w.Write([]byte(`{"code":401,"msg":"http://`+r.Host+`/static/view/home.html"}`))
		resp:=util.RespMsg{
			Code:0,
			Msg:"ok",
			Data:struct{
				Location string
				Username string
				Token string
			}{
				Location:"/static/view/home.html",
				Username:username,
				Token:token,
			},
		}
		c.Data(http.StatusOK,"application",resp.JSONBytes())
		//w.Write(resp.JSONBytes())
}
func UserInfoHandler(c *gin.Context){

	//1.解析请求参数
	//r.ParseForm()
	username:=c.Request.FormValue("username")
	//token:=r.Form.Get("token")
	////2.验证token是否有效
	//isVaildToken:=IsTokenVaild(token)
	//if !isVaildToken{
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息
	user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		//w.WriteHeader(http.StatusForbidden)
		c.JSON(http.StatusForbidden,gin.H{})
		return
	}
	//4.组装并且响应用户数据
	//resp:=util.RespMsg{
	//	Code: 0,
	//	Msg: "OK",
	//	Data:user,
	//}
	//w.Write(resp.JSONBytes())
	c.JSON(http.StatusOK,gin.H{
		"msg":"ok",
		"code":0,
		"Data":user,
	})
}

func IsTokenVaild(token string)bool{
	//TODO 判断token时效性
	if len(token)!=40{
		return false
	}
	//从数据表tbl_user
	return true
}
func GenToken(username string)string{
	//40位字符
	ts:=fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix:=util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix+ts[:8]
}
