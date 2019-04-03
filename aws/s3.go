package aws

import (
	"bytes"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
	dbHelper "github.com/jpnewman/artifactory-s3-export/dbs"
	"github.com/jpnewman/artifactory-s3-export/models"
	"github.com/samonzeweb/godb"
	"github.com/spf13/viper"
)

// GetS3Objects - Get S3 Objects
func GetS3Objects(s *session.Session, sqliteDb *godb.DB, repo string) {
	glog.Info("Updating S3 Objects in database")
	s3Key := path.Join(viper.GetString("aws.s3_key"), repo)

	s3svc := s3.New(s)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(viper.GetString("aws.s3_bucket")),
		Prefix: aws.String(s3Key),
	}

	err := s3svc.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, value := range page.Contents {
			var s3Obj models.S3Object
			s3Obj.Key = *value.Key
			s3Obj.Size = uint64(*value.Size)

			dbHelper.InsertOrUpdate(sqliteDb, &s3Obj)
		}

		return true
	})

	if err != nil {
		panic(err.Error())
	}
}

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
func UploadFileToS3(s *session.Session, filePath string, s3Key string) (models.S3Object, error) {
	glog.Infof("Uploading file to S3: %s\n", filePath)
	var s3Obj models.S3Object

	file, err := os.Open(filePath)
	if err != nil {
		return s3Obj, err
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

	s3Obj.Key = s3Key
	s3Obj.Size = uint64(size)

	if err != nil {
		s3Obj.UploadError = err.Error()
	}

	return s3Obj, err
}
