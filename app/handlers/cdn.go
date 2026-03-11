package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/desygner/desygner-toolchain/app/models"
	facades2 "github.com/desygner/desygner-toolchain/packages/aws/facades"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/str"
	url2 "net/url"
	"regexp"
	"strings"
	"time"
)

type CDN struct {
}

func (receiver *CDN) ResolveBucket(ctx context.Context, url string) (res *models.CDNMap, err error) {
	parsed, err := url2.Parse(url)
	if err != nil {
		return
	}
	cfn := facades2.Aws().Cloudfront()
	distID, err := receiver.getDistribution(ctx, cfn, parsed.Host)
	if err != nil {
		return
	}
	bucket, err := receiver.getBucket(ctx, cfn, distID, parsed.Path)
	res = &models.CDNMap{
		DistributionID: distID,
		Bucket:         bucket,
	}
	return
}

func (receiver *CDN) ResolveDistribution(ctx context.Context, url string) (res *models.CDNMap, err error) {
	parsed, err := url2.Parse(url)
	if err != nil {
		return
	}
	cfn := facades2.Aws().Cloudfront()
	distID, err := receiver.getDistribution(ctx, cfn, parsed.Host)
	if err != nil {
		return
	}
	res = &models.CDNMap{
		DistributionID: distID,
	}
	return
}

func (receiver *CDN) getDistribution(ctx context.Context, cfn *cloudfront.Client, host string) (distID string, err error) {
	cacher := facades.Cache().WithContext(ctx)
	awsCli := facades2.Aws()
	ttl := 24 * time.Hour
	key := fmt.Sprintf("cdn.%s.distribution.%s", awsCli.GetProfile(), host)
	logger := facades.Log().WithContext(ctx)
	dist, err := cacher.Remember(key, ttl, func() (res any, err error) {
		output, err := receiver.loadDistributions(ctx, cfn)
		if err != nil {
			return
		}
		for _, item := range output.DistributionList.Items {
			logger.Infof("matching %s to distribution %s", host, aws.ToString(item.Id))
			logger.Infof("aliases: %s", item.Aliases.Items)
			if len(item.Aliases.Items) < 1 {
				continue
			}
			if !str.Of(host).StartsWith(item.Aliases.Items...) {
				continue
			}
			return aws.ToString(item.Id), nil
		}
		return "", fmt.Errorf("distribution not found")
	})
	if err != nil {
		return
	}
	distID = dist.(string)
	return
}

func (receiver *CDN) loadDistributions(ctx context.Context, cfn *cloudfront.Client) (output *cloudfront.ListDistributionsOutput, err error) {
	cacher := facades.Cache().WithContext(ctx)
	logger := facades.Log().WithContext(ctx).With(map[string]any{"profile": facades2.Aws().GetProfile()})
	key := fmt.Sprintf("cdn.%s.distributions", facades2.Aws().GetProfile())
	logger.Info("loading distributions")
	outputRes, err := cacher.Remember(key, time.Duration(facades.Config().GetInt("aws.ttl"))*time.Second, func() (output any, err error) {
		paginator := cloudfront.NewListDistributionsPaginator(cfn, &cloudfront.ListDistributionsInput{})
		var outRes *cloudfront.ListDistributionsOutput

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("listing distributions: %w", err)
			}
			if outRes == nil {
				outRes = page
				continue
			}
			outRes.DistributionList.Items = append(outRes.DistributionList.Items, page.DistributionList.Items...)
		}
		return outRes, nil
	})
	if err != nil {
		return
	}
	output = outputRes.(*cloudfront.ListDistributionsOutput)
	return
}

func (receiver *CDN) getBucket(ctx context.Context, cfn *cloudfront.Client, distID string, path string) (string, error) {
	var (
		domainName string
		found      bool
	)
	cacher := facades.Cache().WithContext(ctx)
	awsCli := facades2.Aws()
	ttl := 24 * time.Hour
	key := fmt.Sprintf("cdn.%s.bucket.%s.%s", awsCli.GetProfile(), distID, path)
	if cacher.Has(key) {
		bucketID := cacher.GetString(key)
		return bucketID, nil
	}
	out, err := cfn.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
		Id: aws.String(distID),
	})
	if err != nil {
		return "", fmt.Errorf("getting distribution config: %w", err)
	}
	origins := out.DistributionConfig.Origins

	for _, bhv := range out.DistributionConfig.CacheBehaviors.Items {
		pattern := aws.ToString(bhv.PathPattern)
		pattern = strings.ReplaceAll(pattern, "*", "__WILDCARD__")
		pattern = regexp.QuoteMeta(pattern)
		pattern = strings.ReplaceAll(pattern, "__WILDCARD__", ".*")
		pattern = fmt.Sprintf("^%s$", pattern)
		compiler := regexp.MustCompile(pattern)
		logPattern := facades.Log().With(map[string]any{
			"pattern": pattern,
			"path":    path,
		})
		logPattern.Info(
			"matching path pattern",
		)
		if compiler.MatchString(path) {
			logPattern.With(map[string]any{"originID": aws.ToString(bhv.TargetOriginId)}).Info("matching found. lookup the origin")
			found = true
			for _, origin := range origins.Items {
				logPattern.With(map[string]any{
					"target": aws.ToString(origin.Id),
					"source": aws.ToString(bhv.TargetOriginId),
					"path":   aws.ToString(origin.OriginPath),
				}).Info("checking origin")
				if aws.ToString(origin.Id) == aws.ToString(bhv.TargetOriginId) &&
					(strings.Contains(aws.ToString(bhv.TargetOriginId), "amazonaws.com") || strings.Contains(aws.ToString(bhv.TargetOriginId), "bucket")) {
					domainName = aws.ToString(origin.DomainName) + aws.ToString(origin.OriginPath)
					break
				}
			}
			break
		}
	}
	if domainName == "" && found {
		return "", fmt.Errorf("origin not found")
	}
	if domainName == "" {
		domainName = aws.ToString(out.DistributionConfig.DefaultCacheBehavior.TargetOriginId)
	}
	facades.Log().With(map[string]any{"domain": domainName}).Info("domain name")
	s3BucketPattern := regexp.MustCompile(`^(?P<bucket>.*)\.s3\.([^.]+)\.amazonaws\.com$`)
	bucketIndex := s3BucketPattern.SubexpIndex("bucket")
	matches := s3BucketPattern.FindStringSubmatch(domainName)
	if len(matches) == 0 {
		return strings.ReplaceAll(domainName, "S3-", ""), nil
	}
	bucket := matches[bucketIndex]
	_ = cacher.Put(key, bucket, ttl)
	return bucket, nil
}
