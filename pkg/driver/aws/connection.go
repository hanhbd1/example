package aws

import (
	"context"
	"time"

	gocoreConfig "example/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Connection struct {
	Region                 string
	SharedConfigProfile    string
	SharedConfigFiles      []string // SDK defaults to ~/.aws/config
	SharedCredentialsFiles []string // SDK defaults to ~/.aws/credentials
	BaseEndpoint           string
	UseFips                bool
	UseDualStack           bool
	MaxAttempts            int
	MaxBackoffDelay        time.Duration
}

func BuildDefaultConnection() *Connection {
	return &Connection{
		Region:                 gocoreConfig.GetString("aws.region"),
		SharedConfigProfile:    gocoreConfig.GetString("aws.configProfile"),
		SharedConfigFiles:      gocoreConfig.GetStringSlice("aws.configFiles"),
		SharedCredentialsFiles: gocoreConfig.GetStringSlice("aws.credentialsFiles"),
		BaseEndpoint:           gocoreConfig.GetString("aws.baseEndpoint"),
		UseFips:                gocoreConfig.GetBool("aws.useFips"),
		UseDualStack:           gocoreConfig.GetBool("aws.useDualStack"),
		MaxAttempts:            gocoreConfig.GetInt("aws.maxAttempts"),
		MaxBackoffDelay:        gocoreConfig.GetDuration("aws.maxBackoffDelay"),
	}
}

func DefaultAwsConfig() (aws.Config, error) {
	optFns := [](func(*config.LoadOptions) error){
		config.WithRetryer(func() aws.Retryer {
			retryer := retry.NewStandard()
			if gocoreConfig.GetInt("aws.maxAttempts") != 0 {
				retry.AddWithMaxAttempts(retryer, gocoreConfig.GetInt("aws.maxAttempts"))
			}
			if gocoreConfig.GetDuration("aws.maxBackoffDelay") != time.Duration(0) {
				retry.AddWithMaxBackoffDelay(retryer, gocoreConfig.GetDuration("aws.maxBackoffDelay"))
			}
			return retryer
		}),
	}
	if gocoreConfig.GetString("aws.region") != "" {
		optFns = append(optFns, config.WithRegion(gocoreConfig.GetString("aws.region")))
	}
	if gocoreConfig.GetString("aws.configProfile") != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(gocoreConfig.GetString("aws.configProfile")))
	}
	if len(gocoreConfig.GetStringSlice("aws.configFiles")) > 0 {
		optFns = append(optFns, config.WithSharedConfigFiles(gocoreConfig.GetStringSlice("aws.configFiles")))
	}
	if len(gocoreConfig.GetStringSlice("aws.credentialsFiles")) > 0 {
		optFns = append(optFns, config.WithSharedCredentialsFiles(gocoreConfig.GetStringSlice("aws.credentialsFiles")))
	}
	return config.LoadDefaultConfig(context.Background(), optFns...)
}

func NewAwsConfig(conn *Connection) (aws.Config, error) {
	optFns := [](func(*config.LoadOptions) error){
		config.WithRetryer(func() aws.Retryer {
			retryer := retry.NewStandard()
			if conn.MaxAttempts != 0 {
				retry.AddWithMaxAttempts(retryer, conn.MaxAttempts)
			}
			if conn.MaxBackoffDelay != time.Duration(0) {
				retry.AddWithMaxBackoffDelay(retryer, conn.MaxBackoffDelay)
			}
			return retryer
		}),
	}
	if conn.Region != "" {
		optFns = append(optFns, config.WithRegion(conn.Region))
	}
	if conn.SharedConfigProfile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(conn.SharedConfigProfile))
	}
	if len(conn.SharedConfigFiles) > 0 {
		optFns = append(optFns, config.WithSharedConfigFiles(conn.SharedConfigFiles))
	}
	if len(conn.SharedCredentialsFiles) > 0 {
		optFns = append(optFns, config.WithSharedCredentialsFiles(conn.SharedCredentialsFiles))
	}
	return config.LoadDefaultConfig(context.Background(), optFns...)
}
