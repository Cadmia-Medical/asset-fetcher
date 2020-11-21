package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"os"
)

type AwsFetcher struct {
	session *session.Session
	bucket  string
}

func CreateAwsFetcher(region string, bucket string) (*AwsFetcher, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	return &AwsFetcher{sess, bucket}, nil

}

func (self *AwsFetcher) GetAssetTag(name string) (string, error) {
	svc := s3.New(self.session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(self.bucket),
		Key:    aws.String(name),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return "", err
	}

	return *result.ETag, nil
}

func (self *AwsFetcher) GetAsset(path string, name string) error {
	file, err := os.Create(path)
	if err != nil {
		return nil
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(self.session)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(self.bucket),
			Key:    aws.String(name),
		})
	if err != nil {
		return err
	}

	return nil
}
