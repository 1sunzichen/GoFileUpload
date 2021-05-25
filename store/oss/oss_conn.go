package oss
import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	cfg "filestore-server/config"
)
var ossCli *oss.Client

func Client() *oss.Client{
	if ossCli!=nil{
		return ossCli
	}
	ossCli,err:=oss.New(cfg.OSSEndpoint,
		cfg.OSSAccessKeyID,cfg.OSSAccessKeySecret)
	if err!=nil{
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}
//Bucket 获取bucket存储空间
func Bucket() *oss.Bucket{
	cli:=Client()
	if cli!=nil{
		bucket,err:=cli.Bucket(cfg.OSSBucket)
		if err!=nil{
			fmt.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}

//DownloadURL 临时授权下载url
func DownloadURL(objName ,filename string)string{
	fmt.Printf("路径%v 文件名%v",objName,filename)

	//err:=Bucket().GetObjectToFile(objName,filename)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return false
	//}else{
	//	return  true
	//}
 signedURL,err:=Bucket().SignURL(objName,oss.HTTPGet,3600)
 if err!=nil{
 	fmt.Println(err.Error())
 	return ""
 }
 return signedURL
}