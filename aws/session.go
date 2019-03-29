package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

func InitAWSSession() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws.aws_region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.access_key"), viper.GetString("aws.secret_key"), "")})
	if err != nil {
		log.Fatal(err)
	}

	return s
}
