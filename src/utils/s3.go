package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"mime/multipart"
)

func UploadToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})
	if err != nil {
		return "", err
	}

	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), fileHeader.Filename)

	svc := s3.New(sess)
	key := "/" + uniqueFileName
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("cdn.animoshi.com"),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://cdn.animoshi.com%s", key), nil
}
