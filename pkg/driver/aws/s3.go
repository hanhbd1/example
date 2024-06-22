package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(conn *Connection) (client *s3.Client, presignClient *s3.PresignClient, err error) {
	var cfg aws.Config
	if conn == nil {
		cfg, err = DefaultAwsConfig()
		if err != nil {
			return
		}
	} else {
		cfg, err = NewAwsConfig(conn)
		if err != nil {
			return
		}
	}
	client = s3.NewFromConfig(cfg)

	presignClient = s3.NewPresignClient(client)
	return
}
