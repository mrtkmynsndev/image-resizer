package registry

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	S3Uploader = SetupS3Uploder()
)

func SetupS3Uploder() *s3manager.Uploader {
	session := session.Must(session.NewSession())

	return s3manager.NewUploader(session)
}
