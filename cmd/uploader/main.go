package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
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

func uploadFile(filename string) {
	fmt.Printf("Iniciando upload do arquivo: %s no bucket: ", filename, bukectKey)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo: %s", filename)
		return
	}
	defer f.Close()
	_, err = s3Cliente.PutObject(&s3.PutObjectInput{Bucket: aws.String(bukectKey), Key: aws.String(filename), Body: f})
	if err != nil {
		fmt.Printf("Erro ao subir arquivo para o s3: %s", filename)
		return
	}
	fmt.Printf("upload do arquivo: %s finalizadado ", filename)
}
