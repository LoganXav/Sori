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
	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write file content: %v", err)
	}

	return nil
}

func UploadToS3(filePath, fileKey string) (any, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Read file content
	fileInfo, _ := file.Stat()
	fileBuffer := make([]byte, fileInfo.Size())
	file.Read(fileBuffer)

	contentType := http.DetectContentType(fileBuffer)

	// Upload to S3
	s3PutObjectOutput, err := svc.PutObject(&s3.PutObjectInput{
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
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Convert the interface to a map
	resultMap := make(map[string]interface{})
	resultMap["BucketKeyEnabled"] = s3PutObjectOutput.BucketKeyEnabled
	resultMap["ChecksumCRC32"] = s3PutObjectOutput.ChecksumCRC32
	resultMap["ChecksumCRC32C"] = s3PutObjectOutput.ChecksumCRC32C
	resultMap["ChecksumSHA1"] = s3PutObjectOutput.ChecksumSHA1
	resultMap["ChecksumSHA256"] = s3PutObjectOutput.ChecksumSHA256
	resultMap["ETag"] = s3PutObjectOutput.ETag
	resultMap["Expiration"] = s3PutObjectOutput.Expiration
	resultMap["RequestCharged"] = s3PutObjectOutput.RequestCharged
	resultMap["SSECustomerAlgorithm"] = s3PutObjectOutput.SSECustomerAlgorithm
	resultMap["SSECustomerKeyMD5"] = s3PutObjectOutput.SSECustomerKeyMD5
	resultMap["SSEKMSEncryptionContext"] = s3PutObjectOutput.SSEKMSEncryptionContext
	resultMap["SSEKMSKeyId"] = s3PutObjectOutput.SSEKMSKeyId
	resultMap["ServerSideEncryption"] = s3PutObjectOutput.ServerSideEncryption
	resultMap["VersionId"] = s3PutObjectOutput.VersionId

	// Generate a pre-signed URL for the uploaded file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configs.GetEnv("AWS_BUCKET")),
		Key:    aws.String(fileKey),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}

	// assign presigned url
	// Add a new attribute to the map
	resultMap["AWSUrl"] = urlStr
	resultMap["filename"] = fileKey
	resultMap["mimetype"] = contentType
	resultMap["size"] = fileInfo.Size()

	return resultMap, nil

	// Generate S3 URL
	// s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", configs.GetEnv("AWS_BUCKET"), configs.GetEnv("AWS_DEFAULT_REGION"), fileKey)
	// return s3URL, nil
}

