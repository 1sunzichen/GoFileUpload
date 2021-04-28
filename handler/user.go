package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt="*#890"
)
//处理用户注册请求
func SignupHandler(w http.ResponseWriter,r *http.Request){
	if r.Method==http.MethodGet{
		data,err:=ioutil.ReadFile("./static/view/signup.html")
		if err!=nil{

		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username:=r.Form.Get("username")
	passwd:=r.Form.Get("password")
	//存在一个之后 不在成功的问题 ，就是唯一索引值  没有插入，插入之后 就可以刷新
	phone:=r.Form.Get("phone")
	if len(username)<3||len(passwd)<5{
		w.Write([]byte("Invaild parameter"))
		return
	}
	enc_passwd:=util.Sha1([]byte(passwd+pwd_salt))
	fmt.Println(enc_passwd,username)
	suc:=dblayer.UserSignUp(username,enc_passwd,phone)
	if suc{
		w.Write([]byte("success"))
	}else{
		w.Write([]byte("Fail"))
	}
}
type result struct {
	url string
}
//SignInHandler 登录接口
func SignInHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	//校验用户名及密码
	username:=r.Form.Get("username")
	password:=r.Form.Get("password")
	encPasswd:=util.Sha1([]byte(password+pwd_salt))
	pwdChecked:=dblayer.UserSignin(username,encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		return
	}
	//生成访问凭证
	token:=GenToken(username)
	upRes:=dblayer.UpdateToken(username,token)
	if !upRes{
		w.Write([]byte("FAILED"))
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
				Location:"http://"+r.Host+"/static/view/home.html",
				Username:username,
				Token:token,
			},
		}
		w.Write(resp.JSONBytes())
}
func UserInfoHandler(w http.ResponseWriter,r *http.Request){
	//1.解析请求参数
	r.ParseForm()
	username:=r.Form.Get("username")
	token:=r.Form.Get("token")
	//2.验证token是否有效
	isVaildToken:=IsTokenVaild(token)
	if !isVaildToken{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//3.查询用户信息
	user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp:=util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenVaild(token string)bool{
	//TODO 判断token时效性
	//从数据表tbl_user
	return true
}
func GenToken(username string)string{
	//40位字符
	ts:=fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix:=util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix+ts[:8]
}
