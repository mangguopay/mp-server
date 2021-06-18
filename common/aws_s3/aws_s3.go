package aws_s3

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	Private_Img_Dir = "simg" // 不能直接访问的图片目录
	Public_Img_Dir  = "img"  // 可以直接访问的图片目录
	App_Dir         = "app"  // 存放app下载目录

	Xlsx_Dir = "biztransfer" // 存放商家批量转账的xlsx文件
	Misc_Dir = "misc"        // 存放一下杂项零碎的文件

	Encrypt_Key = "oUwd5mxg1PDyj8JWFRntQz7AiC4ITfbG"
)

var headSalt = []byte("GPCoO2q7YKMeSsic")

type UploadS3 struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Region          string
}

func NewUploadS3(accessKeyID, secretAccessKey, region, bucket string) *UploadS3 {
	return &UploadS3{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Bucket:          bucket,
	}
}

// 获取上传的session
func (u *UploadS3) GetSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(u.Region),
		Credentials: credentials.NewStaticCredentials(u.AccessKeyID, u.SecretAccessKey, ""),
	})
}

// 获取对象
//
// fileName 对象名称
//
// result, err := uploadS3.GetObject("simg/b182eef2fdeed7aa9b743bcf8a9101ff.jpeg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-12
func (u *UploadS3) GetObject(fileName string) (*s3.GetObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(fileName),
	}

	return s3.New(sess).GetObject(input)
}

// 上传文件到aws的S3
//
// fileName 待上传的文件
// specifyName 指定上传到S3的文件名
// isPublic 目标对象是否可以公开访问
//
// fileName := "C:\\Users\\Administrator\\Pictures\\scenery-03.jpg"
// result, err := UploadFile(fileName, "img/scenery-03.jpg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) UploadFile(fileName string, specifyName string, isPublic bool) (*s3.PutObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, ferr := file.Stat()
	if ferr != nil {
		return nil, ferr
	}

	// 将文件读入buffer
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(u.Bucket),
		Key:           aws.String(specifyName),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	}

	if isPublic { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	return s3.New(sess).PutObject(input)
}

// 通过FileHeader上传
//
// fileHeader
// specifyName 指定上传到S3的文件名
// isPublic 目标对象是否可以公开访问
//
// file, _ := c.FormFile("file") // c为*gin.Context
// result, err := UploadByMultipartFileHeader(file, "img/aaa.jpeg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) UploadByMultipartFileHeader(fileHeader *multipart.FileHeader, specifyName string, isPublic bool) (*s3.PutObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	// 获取文件内容
	fileContent, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	bytesContent, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(u.Bucket),
		Key:           aws.String(specifyName),
		Body:          bytes.NewReader(bytesContent),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(http.DetectContentType(bytesContent)),
	}

	if isPublic { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	return s3.New(sess).PutObject(input)
}

// 通过文件内容直接上传
//
// fileHeader
// fileName 指定上传到S3的文件名
// isPublic 目标对象是否可以公开访问
//
// result, err := UploadByContent(bytesContent, "img/aaa.jpeg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) UploadByContent(content []byte, fileName string, isPublic bool) (*s3.PutObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(u.Bucket),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(content),
		ContentLength: aws.Int64(int64(len(content))),
		ContentType:   aws.String(http.DetectContentType(content)),
	}

	if isPublic { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	return s3.New(sess).PutObject(input)
}

// 在slice的头部添加salt
// @auth xiaoyanchun 2020-05-13
func addHeadSalt(data []byte) []byte {
	// 创建一个大的slice，来装data和salt
	c := make([]byte, len(data)+len(headSalt))

	// 将salt装入
	copy(c, headSalt)

	// 在装入salt后面，再装入data
	copy(c[len(headSalt):], data)
	return c
}

// 去掉通过addHeadSalt函数添加的头部salt
// @auth xiaoyanchun 2020-05-13
func removeHeadSalt(data []byte) []byte {
	dLen := len(data)
	sLen := len(headSalt)

	if dLen < sLen { // 内容的长度比salt还短,不处理
		return data
	}

	// 创建一个排除salt的长度的slice
	c := make([]byte, dLen-sLen)

	// 将salt后面的数据装入
	copy(c, data[sLen:])
	return c
}

// 异或加密
// @auth xiaoyanchun 2020-05-12
func xorEncode(data []byte, key string) []byte {
	keyLen := len(key)
	for i, b := range data {
		data[i] = byte((key[i%keyLen]) ^ b)
	}
	return data
}

// 通过文件内容直接上传-内容加密
//
// content 待上传的内容
// fileName 指定上传到S3的文件名
//
// result, err := UploadByContentEncrypt([]byte("aaaabbbccc"), "img/bbbbb.jpeg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-12
func (u *UploadS3) UploadByContentEncrypt(content []byte, fileName string) (*s3.PutObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	// 在头部添加salt后，再进行异或加密
	encryptCon := addHeadSalt(xorEncode(content, Encrypt_Key))

	return s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(u.Bucket),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(encryptCon),
		ContentLength: aws.Int64(int64(len(encryptCon))),
	})
}

// 获取加密过的对象(内部会自动解码)
//
// fileName 对象名称
//
// bytes, err := GetObjectEncrypt("img/bbbbb.jpeg")
// fmt.Println("err:", err)
//
// @auth xiaoyanchun 2020-05-12
func (u *UploadS3) GetObjectEncrypt(fileName string) ([]byte, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	result, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return nil, err
	}

	// 读取body内容
	bytesContent, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	// 去掉头部的salt后,再进行解码返回
	return xorEncode(removeHeadSalt(bytesContent), Encrypt_Key), nil
}

// 复制对象
//
// srcName 源对象
// destName 目标对象
// isPublic 目标对象是否可以公开访问
//
// srcName := "img/scenery-01.jpg"
// destName := "img/scenery-latest.jpg"
// result, err := uploadS3.CopyObject(srcName, destName)
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-11
func (u *UploadS3) CopyObject(srcName string, destName string, isPublic bool) (*s3.CopyObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(u.Bucket),
		CopySource: aws.String("/" + u.Bucket + "/" + srcName),
		Key:        aws.String(destName),
	}

	if isPublic { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	return s3.New(sess).CopyObject(input)
}

// 列出对象
//
// dirName 目录名称
// pageSize 取多少条记录
// startAfter 从那个记录之后开始
//
// result, err := ListObjects("img", 10, "img/scenery-01.jpg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) ListObjects(dirName string, pageSize int64, startAfter string) (*s3.ListObjectsV2Output, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsV2Input{
		Bucket:     aws.String(u.Bucket),
		Prefix:     aws.String(dirName),
		MaxKeys:    aws.Int64(pageSize),
		StartAfter: aws.String(startAfter),
	}

	return s3.New(sess).ListObjectsV2(input)
}

// 删除对象-单个
//
// fileName 待删除的文件
// dirName 目录名称
//
// result, err := DeleteOne("img/scenery-01.jpg")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) DeleteOne(fileName string) (*s3.DeleteObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	return s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(fileName),
	})
}

// 删除对象-多个
//
// fileNames 待删除的文件(支持多个)
// dirName 目录名称
//
// fileNames := []string{"img/scenery-01.jpg", "img/scenery-02.jpg"}
// result, err := DeleteMulti(fileNames)
// fmt.Println("err:", err)
// fmt.Println("result.Errors:", result.Errors)
//
// @auth xiaoyanchun 2020-05-09
func (u *UploadS3) DeleteMulti(fileNames []string) (*s3.DeleteObjectsOutput, error) {
	// 获取session
	sess, err := u.GetSession()

	if err != nil {
		return nil, err
	}
	if len(fileNames) == 0 {
		return nil, errors.New("fileNames不能为空")
	}

	objList := []*s3.ObjectIdentifier{}
	for _, v := range fileNames {
		objList = append(objList, &s3.ObjectIdentifier{Key: aws.String(v)})
	}

	return s3.New(sess).DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(u.Bucket),
		Delete: &s3.Delete{
			Objects: objList,
		},
	})
}

// 上传.apk文件到aws的S3
//
// fileName 待上传的文件
// specifyName 指定上传到S3的文件名
// isPublic 目标对象是否可以公开访问
//
// fileName := "C:\\Users\\Administrator\\Pictures\\scenery-03.apk"
// result, err := UploadFile(fileName, "app/scenery-03.apk")
// fmt.Println("err:", err)
// fmt.Println("result:", result)
//
// @auth xiaoyanchun 2020-07-13
func (u *UploadS3) UploadAPKFile(fileName string, specifyName string, isPublic bool) (*s3.PutObjectOutput, error) {
	// 获取session
	sess, err := u.GetSession()
	if err != nil {
		return nil, err
	}

	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, ferr := file.Stat()
	if ferr != nil {
		return nil, ferr
	}

	// 将文件读入buffer
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(u.Bucket),
		Key:           aws.String(specifyName),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String("application/vnd.android.package-archive"),
	}

	if isPublic { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	return s3.New(sess).PutObject(input)
}
