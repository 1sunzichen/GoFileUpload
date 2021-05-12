package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
<<<<<<< HEAD
	"time"
=======
	"strconv"
	"time"

     dbplayer "filestore-server/db"
>>>>>>> part5-2
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
		if r.Method=="GET"{
			//返回上传html页面
			data,err:=ioutil.ReadFile("./static/view/upload.html")
			if err !=nil{
				io.WriteString(w,"interError")
				return
			}
			io.WriteString(w,string(data))
		}else if r.Method=="POST"{
            file,head,err:=r.FormFile("file")
            if err!=nil{
				fmt.Printf("failed to get data,err:%s\n",err)
				return
			}
            defer file.Close()
<<<<<<< HEAD
            fileMeta:=meta.FileMeta{
            	FileName: head.Filename,
            	Location:"./upload/"+head.Filename,
            	UploadAt: time.Now().Format("2021-06-02 15:04:05"),
=======


            fileMeta:=meta.FileMeta{
            	FileName: head.Filename,
            	Location: "static/file/"+head.Filename,
            	UploadAt: time.Now().Format("2016-01-02"),
>>>>>>> part5-2

			}
            newFile,err:=os.Create(fileMeta.Location)
            if err!=nil{
            	fmt.Printf("Failed to create file,err%s\n",err.Error())
            	return
			}
			defer newFile.Close()
            fileMeta.FileSize,err=io.Copy(newFile,file)
            if err!=nil{
            	fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
            	return
			}
			newFile.Seek(0,0)
            fileMeta.FileSha1=util.FileSha1(newFile)
<<<<<<< HEAD
            meta.UpdateFileMeta(fileMeta)
			http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
=======
            //meta.UpdateFileMeta(fileMeta)
            _=meta.UpdateFileMetaDb(fileMeta)
            //更新用户表文件表记录
            r.ParseForm()
            username:=r.Form.Get("username")
            suc:=dbplayer.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
            if suc{
				resp:=util.RespMsg{
					Code:0,
					Msg:"ok",
					Data:struct{
						Location string
						Username string
						//Token string
					}{
						Location:"http://"+r.Host+"/static/view/home.html",
						Username:username,
						//Token:token,
					},
				}
				w.Write(resp.JSONBytes())
				//http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
			}else{
				w.Write([]byte("上传失败"))
			}
>>>>>>> part5-2
		}
}
//上传已完成
func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished")
}
<<<<<<< HEAD
// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	fMeta:=meta.GetFileMeta(filehash);
	data,err:=json.Marshal(fMeta)
=======
//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	//对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form。
	//解析结果中，POST或PUT请求主体要优先于URL查询字符串（同名变量，主体的值在查询字符串的值前面）。
	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	//fMeta:=meta.GetFileMeta(filehash)
	fMeta,err:=meta.GetFileMetaDB(filehash)
>>>>>>> part5-2
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
<<<<<<< HEAD
	w.Write(data)
}

//下载信息
func DownloadHandler(w http.ResponseWriter,r *http.Request)  {
	
=======
	d,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(d)
}
//批量获取文件元信息
func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	//对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form。
	//解析结果中，POST或PUT请求主体要优先于URL查询字符串（同名变量，主体的值在查询字符串的值前面）。
	r.ParseForm()
	//字符串转数组
	limitCnt,_:=strconv.Atoi(r.Form.Get("limit"))
	username:=r.Form.Get("username");
	fmt.Println(limitCnt,username)

	fileMetas,err:=dbplayer.QueryUserFileMatas(username,limitCnt)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:=json.Marshal(fileMetas);

	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1:=r.Form.Get("filehash")
	fm:=meta.GetFileMeta(fsha1)
	f,err:=os.Open(fm.Location)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	//小文件
	data,err:=ioutil.ReadAll(f)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//http.
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-disposition","attachment;filename=\""+fm.FileName)
	w.Write(data)
}
func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	opType:=r.Form.Get("op")
	fileSha1:=r.Form.Get("filehash")
	newFileName:=r.Form.Get("filename")
	if opType!="0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method!="POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta:=meta.GetFileMeta(fileSha1)
	curFileMeta.FileName=newFileName
	meta.UpdateFileMeta(curFileMeta)
	data,err:=json.Marshal(curFileMeta)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
//删除接口
func FileDeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filesha1:=r.Form.Get("filehash")
	fmeta:=meta.GetFileMeta(filesha1)
	os.Remove(fmeta.Location)
	meta.RemoveFileMeta(filesha1)

	w.WriteHeader(http.StatusOK)
}
//
func TryFastUploadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	//
	username:=r.Form.Get("username")
	filehash:=r.Form.Get("filehash")
	filename:=r.Form.Get("filename")
	//filesize:=r.Form.Get("filesize")
	filesize,_:=strconv.ParseInt(r.Form.Get("filesize"),10,64)
	//
	fileMeta,err:=meta.GetFileMetaDB(filehash)
	if err!=nil{
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
   fileMetaStatus,_:=json.Marshal(fileMeta)
	if string(fileMetaStatus) == "{}"{
		resp:=util.RespMsg{
			Code:-1,
			Msg: "秒传失败,请访问普通上传接口",


		}
		w.Write(resp.JSONBytes())
		return
	}

	suc:=dbplayer.OnUserFileUploadFinished(username,filehash,filename,filesize)


    if suc{
    	resp:=util.RespMsg{
    		Code: 0,
    		Msg: "秒传成功",


		}
		w.Write(resp.JSONBytes())
	}else{
		resp:=util.RespMsg{
			Code: -2,
			Msg: "秒传失败",


		}
		w.Write(resp.JSONBytes())
	}

>>>>>>> part5-2
}