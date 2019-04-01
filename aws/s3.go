package aws

import (
	"bytes"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jpnewman/artifactory-s3-export/models"
	"github.com/spf13/viper"
)

// CheckS3ObjectExists - Check if S3 Object Exists.
func CheckS3ObjectExists(s *session.Session, s3Key string) bool {

	_, err := s3.New(s).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(viper.GetString("aws.s3_bucket")),
		Key:    aws.String(s3Key),
	})

	if err == nil {
		return true
	}

	return false
}

// UploadFileToS3 - Upload File to S3.
func UploadFileToS3(s *session.Session, node *models.Node) error {
	fileDir := path.Join(viper.GetString("repo.filestore_path"), node.RepoFilePath)
	s3Key := path.Join(viper.GetString("aws.s3_key"), node.Repo, node.NodeName)

	if CheckS3ObjectExists(s, s3Key) {
		node.Uploaded = true
		return nil
	}

	// fmt.Printf("Uploading file to S3: %s", fileDir)

	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(viper.GetString("aws.s3_bucket")),
		Key:                  aws.String(s3Key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	node.Uploaded = true

	return err
}
