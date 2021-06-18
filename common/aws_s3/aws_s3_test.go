package aws_s3

import (
	"fmt"
	"testing"
)

const (
	S3_Region = "ap-southeast-1"

	// user-dev 用户
	S3_AccessKeyID     = "AKIAVEYE56XZ4FTUNOWI"
	S3_SecretAccessKey = "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"
	S3_Bucket          = "modernpay-dev"
)

func TestListObjects(t *testing.T) {
	uploadS3 := NewUploadS3(S3_AccessKeyID, S3_SecretAccessKey, S3_Region, S3_Bucket)
	result, err := uploadS3.ListObjects("img", 10, "img/scenery-01.jpg")
	fmt.Println("err:", err)
	fmt.Println("result:", result)
}

func TestUploadFile(t *testing.T) {
	uploadS3 := NewUploadS3(S3_AccessKeyID, S3_SecretAccessKey, S3_Region, S3_Bucket)
	fileName := "C:\\Users\\Administrator\\Pictures\\scenery-05.jpg"
	result, err := uploadS3.UploadFile(fileName, "img/scenery-05.jpg", true)
	fmt.Println("err:", err)
	fmt.Println("result:", result)
}

func TestDeleteOne(t *testing.T) {
	uploadS3 := NewUploadS3(S3_AccessKeyID, S3_SecretAccessKey, S3_Region, S3_Bucket)
	result, err := uploadS3.DeleteOne("img/scenery-01.jpg")
	fmt.Println("err:", err)
	fmt.Println("result:", result)
}

func TestDeleteMulti(t *testing.T) {
	uploadS3 := NewUploadS3(S3_AccessKeyID, S3_SecretAccessKey, S3_Region, S3_Bucket)
	fileNames := []string{"img/scenery-02.jpg", "img/scenery-03.jpg"}
	result, err := uploadS3.DeleteMulti(fileNames)
	fmt.Println("err:", err)
	fmt.Println("result.Errors:", result.Errors)
}
