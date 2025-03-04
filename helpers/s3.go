package helpers

import (
	configs "LoganXav/sori/configs"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWS Session
var sess *session.Session

// AWS S3 Service
var svc *s3.S3

// Define a struct for the upload result to maintain type safety
type UploadResult struct {
    AWSUrl   string
    Filename string
    Mimetype string
    Size     int64
}

func StartAwsSession() error {
	var err error
	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String(configs.GetEnv("AWS_DEFAULT_REGION")),
		Credentials: credentials.NewStaticCredentials(configs.GetEnv("AWS_ACCESS_KEY_ID"), configs.GetEnv("AWS_SECRET_ACCESS_KEY"), ""),
	})

	if err != nil {
		return err
	}
	svc = s3.New(sess)

	return nil
}

func GetPresignAWSS3(fileKey string) (string, error) {

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:    aws.String(fileKey),
	})
	return req.Presign(15 * time.Minute)
}

func DeleteFromAWSS3(fileKey string) (bool, error) {

	// Define the parameters for the delete object operation
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:    aws.String(fileKey),
	}

	// Delete the object
	_, errD := svc.DeleteObject(params)
	if errD != nil {
		return false, fmt.Errorf("failed to delete file: %w", errD)
	}
	return true, nil
}

func DownloadFromS3(fileKey, destPath string) error {
	// Create the file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", destPath, err)
	}
	defer file.Close()

	// Get object from S3
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return fmt.Errorf("failed to download file from S3: %v", err)
	}

	defer output.Body.Close()


	// Write object content to file
	_, errCopy := io.Copy(file, output.Body)
	if errCopy != nil {
		return fmt.Errorf("failed to write file content: %v", errCopy)
	}


	return nil
}

func UploadToS3(filePath, fileKey string) (*UploadResult, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Read file content
	fileInfo, _ := file.Stat()
	fileBuffer := make([]byte, fileInfo.Size())
	file.Read(fileBuffer)

	contentType := http.DetectContentType(fileBuffer)

	// Upload to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:                  aws.String(fileKey),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileInfo.Size()),
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate a pre-signed URL for the uploaded file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:    aws.String(fileKey),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
        AWSUrl:   urlStr,
        Filename: fileKey,
        Mimetype: contentType,
        Size:     fileInfo.Size(),
    }, nil
}

