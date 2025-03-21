package lib

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Uploader struct {
	logger     Logger
	BucketName string
	Region     string
	AccessKey  string
	SecretKey  string
	downloader *s3manager.Downloader
}

func NewS3Uploader(logger Logger, env Env) S3Uploader {
	return S3Uploader{
		logger:     logger,
		BucketName: env.AWSBucketName,
		Region:     env.AWSRegion,
		AccessKey:  env.AWSAccessKeyID,
		SecretKey:  env.AWSSecretAccessKey,
		downloader: s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(env.AWSRegion)})),
	}
}

func (u *S3Uploader) GeneratePreSignedURL(fileName string) (string, error) {
	// Create a downloader with the session and default options
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(u.Region),
		Credentials: credentials.NewStaticCredentials(
			u.AccessKey,
			u.SecretKey,
			"",
		),
		LogLevel: aws.LogLevel(aws.LogDebugWithRequestErrors),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	objectKey := fileName

	u.logger.Debug("fileName: ", fileName)
	u.logger.Debug("objectKey: ", objectKey)
	u.logger.Debug("s3BucketName: ", u.BucketName)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(objectKey),
		ACL:    aws.String("public-read"),
	})

	str, err := req.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("failed to sign request", err)
		return "", err
	}

	u.logger.Debug("Pre-signed URL:", str)
	return str, nil
}

// UploadFile uploads a local file to S3
func (u *S3Uploader) UploadFileFromPath(filePath string) (string, error) {
	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get the file name to use as the S3 key
	fileName := filepath.Base(filePath)

	// Use the existing Upload method to upload the file
	fileURL, err := u.Upload(file, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return fileURL, nil
}

func (u *S3Uploader) Upload(file io.Reader, fileName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(u.Region),
		Credentials: credentials.NewStaticCredentials(
			u.AccessKey,
			u.SecretKey,
			"",
		),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	// Detect the content type of the file
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(u.BucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.BucketName, u.Region, fileName)
	return fileURL, nil
}

// DownloadFile downloads a file from S3
func (u *S3Uploader) DownloadFile(fileURL string) (io.Reader, error) {
	// Parse the file URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Failed to parse file URL: %v", err))
		return nil, err
	}
	s3ObjectKey := parsedURL.Path[1:]

	// Debug: log the parsed S3 object key
	u.logger.Debug(fmt.Sprintf("Parsed S3 object key: %v", s3ObjectKey))

	// Create a buffer to write the S3 object contents to.
	buff := &aws.WriteAtBuffer{}

	// Write the contents of the S3 object to the buffer
	_, err = u.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(s3ObjectKey),
	})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Failed to download file from S3: %v", err))
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	// Debugging: Log the downloaded file size
	u.logger.Debug("Downloaded file size: ", len(buff.Bytes()))

	// Check if the file is empty after downloading
	if len(buff.Bytes()) == 0 {
		u.logger.Warn("Downloaded file is empty.")
		return nil, fmt.Errorf("downloaded file is empty")
	}

	// Return the file content as a reader
	return bytes.NewReader(buff.Bytes()), nil
}


// DownloadFile downloads a file from Google Drive
func (u *S3Uploader) DownloadImageFromGoogleDrive(url, filename string) error {
	// Manipulate the Google Drive URL if necessary
	if strings.Contains(url, "drive.google.com") {
		url = strings.Replace(url, "open?id=", "uc?export=download&id=", 1)
		url = strings.Replace(url, "file/d/", "uc?export=download&id=", 1)
		url = strings.Replace(url, "/view?usp=sharing", "", 1)
		url = strings.Replace(url, "/view", "", 1)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
