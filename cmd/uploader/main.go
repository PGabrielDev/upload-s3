package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Cliente *s3.S3
	bukectKey string
)

func init() {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				"exemplo",
				"exemplo",
				"",
			)})
	if err != nil {
		panic(err)
	}
	s3Cliente = s3.New(sess)
	bukectKey = "uploader-bukect-s3"
}

func main() {

}
