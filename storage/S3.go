package storage

import (
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

)


func UploadFileToS3(file multipart.File, path string) (string, error) {

	s3_region := os.Getenv("AWS_S3_REGION")
	aws_access_key := os.Getenv("AWS_ACCESS_KEY_ID")
	aws_secret_access_key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	s3_bucket := os.Getenv("AWS_BUCKET_NAME")

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(s3_region),
        Credentials: credentials.NewStaticCredentials(aws_access_key, aws_secret_access_key, ""),
    })
    if err != nil {
        return "", err
    }

    uploader := s3manager.NewUploader(sess)

    // Upload input parameters
    prams := &s3manager.UploadInput{
        Bucket: aws.String(s3_bucket),
        Key:    aws.String(path),
        Body:   file,
    }

    // Perform an upload.
    result, err := uploader.Upload(prams)
    if err != nil {
        return "", err
    }

    return result.Location, nil
}