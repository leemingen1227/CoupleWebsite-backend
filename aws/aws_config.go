package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/leemingen1227/couple-server/util"
)

// type Config struct {
// 	AwsAccessKey string `mapstructure:"AWS_ACCESS_KEY_ID"`
// 	AwsSecretKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
// 	AwsRegion    string `mapstructure:"AWS_REGION"`
// }

func InitAWS() (*session.Session, error) {
	config, err := util.LoadConfig(".")
    
    // Create a new session
    sess, err := session.NewSession(&aws.Config{
	    Region:      aws.String(config.AwsBucketRegion),
	    Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyID, config.AwsSecretKey, ""),
    })
    if err != nil {
	    return nil, err
    }
    
    return sess, nil
}
