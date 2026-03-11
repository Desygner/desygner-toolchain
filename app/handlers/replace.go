package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	s4 "github.com/aws/aws-sdk-go-v2/service/s3"
	types2 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/desygner/desygner-toolchain/app/models"
	facades3 "github.com/desygner/desygner-toolchain/packages/aws/facades"
	"github.com/desygner/desygner-toolchain/packages/slack/extension"
	facades2 "github.com/desygner/desygner-toolchain/packages/slack/facades"
	"github.com/gookit/goutil/syncs"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
	"github.com/slack-go/slack"
	"io"
	"net/url"
	"os"
	path2 "path"
	_ "strconv"
	"strings"
	_ "time"
)

type Replace struct {
	CDN
}

type ReplaceInput struct {
	extension.Input
	Source string
	Target string
}

func GetBucketFromS3URL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(u.Hostname(), ".")
	pathParts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")

	// Virtual-hosted style:
	//   my-bucket.s3.amazonaws.com
	//   my-bucket.s3.us-east-1.amazonaws.com
	//   my-bucket.s3-accelerate.amazonaws.com
	if hostParts[0] != "s3" && !strings.HasPrefix(hostParts[0], "s3-") {
		return hostParts[0], nil
	}

	// Path-style:
	//   s3.amazonaws.com/my-bucket/key
	//   s3.us-east-1.amazonaws.com/my-bucket/key
	if len(pathParts) > 0 && pathParts[0] != "" {
		return pathParts[0], nil
	}

	return "", fmt.Errorf("could not determine bucket from URL: %s", rawURL)
}

func (r *Replace) Handle(ctx context.Context, input *ReplaceInput) (err error) {
	var (
		logger              = facades.Log()
		pathUrl             = input.Target
		sourceUrl           = input.Source
		parentTS            = input.SlackTimestamp()
		awsCli              = facades3.Aws()
		s3                  = awsCli.S3()
		slackCli            = facades2.Slack().Client()
		s3Object            *s4.GetObjectOutput
		backupContentLength int64
		mimeType            string
	)
	parsed, _ := url.Parse(pathUrl)
	basePath := path2.Base(parsed.Path)
	key := strings.TrimLeft(parsed.EscapedPath(), "/")
	tempBackup, err := os.CreateTemp("", basePath+"-backup")
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", path2.Base(parsed.Path)+"-replace")
	if err != nil {
		return err
	}

	defer func() {
		tempFile.Close()
		tempBackup.Close()
		if err != nil {
			input.SendMessage(ctx, fmt.Sprintf(":x: Failed to replace %s. Reason: %s", pathUrl, err.Error()))
		}
	}()
	input.SendMessage(ctx, ":spin: Looking up the location for "+pathUrl)

	logger.Infof("resolving %s", pathUrl)
	// try to resolve the bucket first
	bucketName, _ := GetBucketFromS3URL(pathUrl)
	bucketConfig := &models.CDNMap{
		Bucket: bucketName,
	}
	if bucketConfig.Bucket == "" {
		bucketConfig, err = r.ResolveBucket(ctx, pathUrl)
	} else {
		key = strings.ReplaceAll(key, bucketConfig.Bucket, "")
		key = strings.TrimLeft(key, "/")
	}
	logger.Infof("resolving %s", pathUrl)
	if err != nil {
		return err
	}
	logger.Infof("resolved %s at %s - %s", pathUrl, bucketConfig.Bucket, key)
	logger.Infof("downloading %s and %s", sourceUrl, pathUrl)
	input.SendMessage(ctx, fmt.Sprintf(":spin: Location found at %s. Download the original item as backup", bucketConfig.Bucket))
	errsync := syncs.NewErrGroup()
	errsync.Go(func() error {
		logger.Infof("downloading %s", sourceUrl)
		err2 := slackCli.GetFileContext(ctx, sourceUrl, tempFile)
		if err2 != nil {
			logger.Errorf("failed to download slack file. error: %s. ", err2.Error())
			return err2
		}
		mimeType, err2 = file.MimeType(tempFile.Name())
		return err2
	})
	errsync.Go(func() (err2 error) {
		lg := logger.With(map[string]any{"key": key, "bucket": bucketConfig.Bucket})
		lg.Infof("downloading %s", pathUrl)
		s3Object, err2 = s3.GetObject(ctx, &s4.GetObjectInput{
			Bucket: aws.String(bucketConfig.Bucket),
			Key:    aws.String(key),
		})

		if err2 != nil {
			lg.Error("failed to download s3 object. " + err2.Error())
			return
		}
		backupContentLength, err2 = io.Copy(tempBackup, s3Object.Body)
		return
	})
	err = errsync.Wait()
	if err != nil {
		return err
	}
	logger.Infof("download complete. %d bytes downloaded", backupContentLength)
	logger.Infof("uploading %s and %s", sourceUrl, pathUrl)
	input.SendMessage(ctx, ":spin: Download complete. Replace the original item with the new one")
	errsync2 := syncs.NewErrGroup()
	errsync2.Go(func() error {
		logger.Infof("uploading %s to s3. mime type %s", tempBackup.Name(), mimeType)
		tempBackup.Seek(0, io.SeekStart)
		_, err := slackCli.UploadFileV2Context(ctx, slack.UploadFileV2Parameters{
			Channel:         input.SlackChannel,
			ThreadTimestamp: parentTS,
			Title:           path.Base(basePath),
			InitialComment:  fmt.Sprintf("Backup copy of: %s", pathUrl),
			Reader:          tempBackup,
			Filename:        path.Base(basePath),
			FileSize:        int(backupContentLength),
		})
		return err
	})
	errsync2.Go(func() error {
		tempFile.Seek(0, io.SeekStart)
		logger.Infof("uploading %s to s3. mime type %s", tempFile.Name(), mimeType)
		_, err := s3.PutObject(ctx, &s4.PutObjectInput{
			Bucket:      aws.String(bucketConfig.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(mimeType),
			ACL:         types2.ObjectCannedACLPublicRead,
			Body:        tempFile,
		})
		return err
	})
	err = errsync2.Wait()
	input.SendMessage(ctx, ":greentick1: File has been replaced")
	return nil
}
