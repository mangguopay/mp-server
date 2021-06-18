package common

import (
	"a.a/mp-server/common/aws_s3"
	"a.a/mp-server/common/ss_struct"
)

var UploadS3 *aws_s3.UploadS3

// 初始化上传对象
func InitUploadS3(s3Conf ss_struct.Awss3Conf) {
	UploadS3 = aws_s3.NewUploadS3(s3Conf.AccessKeyId, s3Conf.SecretAccessKey, s3Conf.Region, s3Conf.Bucket)
}
