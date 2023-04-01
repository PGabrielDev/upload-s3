package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"sync"
)

var (
	s3Cliente *s3.S3
	bukectKey string
	wg        *sync.WaitGroup
)

func init() {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				"Sua acess key",
				"sua secret key",
				"",
			)})
	if err != nil {
		panic(err)
	}
	s3Cliente = s3.New(sess)
	bukectKey = "uploader-bucket-s3"
	wg = &sync.WaitGroup{}
}

func main() {
	controlUpload := make(chan struct{}, 50)
	errFile := make(chan string, 10)
	go func() {
		for {
			select {
			case filepath := <-errFile:
				wg.Add(1)
				controlUpload <- struct{}{}
				go uploadFile(filepath, controlUpload, errFile)
			}
		}
	}()
	dir, err := os.Open("../generator")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	for i := 0; i <= 5; i++ {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error ao ler o arquivo")
			continue
		}
		wg.Add(1)
		controlUpload <- struct{}{}
		go uploadFile(files[0].Name(), controlUpload, errFile)
	}

	wg.Wait()

}

func uploadFile(filename string, controlUpload <-chan struct{}, errFile chan<- string) {
	defer wg.Done()
	fmt.Printf("Iniciando upload do arquivo: %s no bucket: %s \n", filename, bukectKey)
	filePath := "../generator/" + filename
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo: %s\n", filename)
		fmt.Println(err.Error())
		<-controlUpload
		errFile <- filePath
		return
	}
	defer f.Close()
	_, err = s3Cliente.PutObject(&s3.PutObjectInput{Bucket: aws.String(bukectKey), Key: aws.String(filename), Body: f})
	if err != nil {
		fmt.Printf("Erro ao subir arquivo para o s3: %s", filename)
		<-controlUpload
		errFile <- filePath
		return
	}
	fmt.Printf("upload do arquivo: %s finalizadado ", filename)
	<-controlUpload
}
