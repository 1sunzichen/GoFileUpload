package ceph
import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)
var cephConn *s3.S3
func GetCephConnection() *s3.S3 {
	//1.
	if(cephConn!=nil){
      return cephConn
	}
	auth:=aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}
	curRegion:=aws.Region{
		Name:                 "default",
		EC2Endpoint:          "",
		S3Endpoint:           "",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	return s3.New(auth,curRegion)
}
//获得bucket 对象
func GetCephBucket(bucket string)*s3.Bucket{
	conn:=GetCephConnection()
	return conn.Bucket(bucket)
}
