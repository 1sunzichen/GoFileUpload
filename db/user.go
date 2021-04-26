package db
import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

func UserSignUp(username string,encpasswd string,phone string)bool {
	stmt,err:=mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`,`phone`) " +
		"values(?,?,?)")
	if err!=nil{
		fmt.Println("Failed to insert,err:"+err.Error())
		return false
	}
	defer stmt.Close()
	ret,err2:=stmt.Exec(username,encpasswd,phone)
	if err2!=nil{
		fmt.Println("Failed to insert,err:"+err.Error())
		return false
	}
	if rf,err3:=ret.RowsAffected();nil==err3{
		if rf<=0{
			fmt.Print("此账号注册过")
			return false
		}
		return true
	}

	return false
}

func UserSignin(username string,encpwd string)bool{
	stmt,err:=mydb.DBConn().Prepare("select user_pwd from tbl_user where user_name=? limit 1 ")
	if err !=nil{
		fmt.Println(err.Error())
		return false
	}
	rows,err:=stmt.Query(username)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	} else if rows==nil{
		fmt.Println("username not found "+username)
		return false
	}



	for rows.Next() {
		var (user_pwd string)
		if ers := rows.Scan(&user_pwd); ers != nil {
			fmt.Println(ers.Error())
			return false
		}
		fmt.Printf("user_pwd2 %s %s\n", user_pwd,encpwd)
		if user_pwd==encpwd{
			return true
		}
		return false
	}
		//list = append(list, item)


	//pRows:=mydb.Parses(rows)
	//if len(pRows)>0&&string(pRows[0]["user_pwd"].([]byte))==encpwd{
	//	return true
	//}
	return false


}
func UpdateToken(username string,token string)bool{
	stmt,err:=mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`)values(?,?)")
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err=stmt.Exec(username,token)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}
