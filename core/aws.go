package core

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"

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
		Key:    aws.String(fmt.Sprintf("%v/%v.zip", name, name)),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return "", err
	}

	return *result.ETag, nil
}

func (self *AwsFetcher) GetAsset(path string, name string) error {
	file, err := os.Create(fmt.Sprintf("%v/%v/%v.zip", path, name, name))
	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(self.session)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(self.bucket),
			Key:    aws.String(fmt.Sprintf("%v/%v.zip", name, name)),
		})
	if err != nil {
		return err
	}

	// Unzip asset
	self.unzip(fmt.Sprintf("%v/%v", path, name), fmt.Sprintf("%v/%v/%v.zip", path, name, name))

	return nil
}

func (self *AwsFetcher) unzip(path string, name string) error {
	reader, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {

		fpath := filepath.Join(path, file.Name)

		// Handle directories
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
	}

	return nil
}
