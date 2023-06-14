package aws

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"
	"github.com/google/uuid"
	"mime/multipart"
	"github.com/leemingen1227/couple-server/util"
)

func UploadImageToS3(sess *session.Session, file *multipart.FileHeader) (string, error) {
	config, err := util.LoadConfig(".")

    uploader := s3manager.NewUploader(sess)

	// Generate a new UUID for the image.
	imageId := uuid.New().String()

	// Open the file.
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(config.AwsBucketName),
        Key:    aws.String(imageId),
        Body:   src,
    })

    return imageId, err
}


func GetSignedURL(sess *session.Session, key string) (string, error) {
    config, err := util.LoadConfig(".")

    // Create S3 service client
    svc := s3.New(sess)

    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(config.AwsBucketName),
        Key:    aws.String(key),
    })

    urlStr, err := req.Presign(15 * time.Minute)

    if err != nil {
        return "", err
    }

    return urlStr, nil
}

