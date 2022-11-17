package sdk

import (
	osssdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
	commonsdk "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
)

type ProxyOssClient struct {
	*osssdk.Client

	secretName      string
	akExpireHandler service.AKExpireHandler
}

func (client *ProxyOssClient) Shutdown() {
	err := closeOssClient(client, client.secretName)
	if err != nil {
		logger.GetCommonLogger(constants.LoggerName).Errorf("action:Shutdown", err)
	}
}

func (client ProxyOssClient) judgeTempAKExpire(err error) bool {
	return client.akExpireHandler.JudgeAKExpire(err)
}

func (client ProxyOssClient) Bucket(bucketName string) (*ProxyBucket, error) {
	bucket, err := client.Client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return &ProxyBucket{
		bucket,
		client.secretName,
		client.akExpireHandler,
	}, nil
}

func (client ProxyOssClient) CreateBucket(bucketName string, options ...osssdk.Option) error {
	err := client.Client.CreateBucket(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.CreateBucket(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) CreateBucketXml(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.CreateBucketXml(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.CreateBucketXml(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) ListBuckets(options ...osssdk.Option) (osssdk.ListBucketsResult, error) {
	res, err := client.Client.ListBuckets(options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.ListBuckets(options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) ListCloudBoxes(options ...osssdk.Option) (osssdk.ListCloudBoxResult, error) {
	res, err := client.Client.ListCloudBoxes(options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.ListCloudBoxes(options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) IsBucketExist(bucketName string) (bool, error) {
	res, err := client.Client.IsBucketExist(bucketName)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.IsBucketExist(bucketName)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucket(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucket(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucket(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketLocation(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketLocation(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketLocation(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketACL(bucketName string, bucketACL osssdk.ACLType, options ...osssdk.Option) error {
	err := client.Client.SetBucketACL(bucketName, bucketACL, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketACL(bucketName, bucketACL, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketACL(bucketName string, options ...osssdk.Option) (osssdk.GetBucketACLResult, error) {
	res, err := client.Client.GetBucketACL(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketACL(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketLifecycle(bucketName string, rules []osssdk.LifecycleRule, options ...osssdk.Option) error {
	err := client.Client.SetBucketLifecycle(bucketName, rules, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketLifecycle(bucketName, rules, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketLifecycleXml(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.SetBucketLifecycleXml(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketLifecycleXml(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketLifecycle(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketLifecycle(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketLifecycle(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketLifecycle(bucketName string, options ...osssdk.Option) (osssdk.GetBucketLifecycleResult, error) {
	res, err := client.Client.GetBucketLifecycle(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketLifecycle(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketLifecycleXml(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketLifecycleXml(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketLifecycleXml(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketReferer(bucketName string, referers []string, allowEmptyReferer bool, options ...osssdk.Option) error {
	err := client.Client.SetBucketReferer(bucketName, referers, allowEmptyReferer, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketReferer(bucketName, referers, allowEmptyReferer, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketReferer(bucketName string, options ...osssdk.Option) (osssdk.GetBucketRefererResult, error) {
	res, err := client.Client.GetBucketReferer(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketReferer(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketLogging(bucketName, targetBucket, targetPrefix string, isEnable bool, options ...osssdk.Option) error {
	err := client.Client.SetBucketLogging(bucketName, targetBucket, targetPrefix, isEnable, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketLogging(bucketName, targetBucket, targetPrefix, isEnable, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketLogging(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketLogging(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketLogging(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketLogging(bucketName string, options ...osssdk.Option) (osssdk.GetBucketLoggingResult, error) {
	res, err := client.Client.GetBucketLogging(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketLogging(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketWebsite(bucketName, indexDocument, errorDocument string, options ...osssdk.Option) error {
	err := client.Client.SetBucketWebsite(bucketName, indexDocument, errorDocument, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketWebsite(bucketName, indexDocument, errorDocument, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketWebsiteDetail(bucketName string, wxml osssdk.WebsiteXML, options ...osssdk.Option) error {
	err := client.Client.SetBucketWebsiteDetail(bucketName, wxml, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketWebsiteDetail(bucketName, wxml, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketWebsiteXml(bucketName string, webXml string, options ...osssdk.Option) error {
	err := client.Client.SetBucketWebsiteXml(bucketName, webXml, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketWebsiteXml(bucketName, webXml, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketWebsite(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketWebsite(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketWebsite(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) OpenMetaQuery(bucketName string, options ...osssdk.Option) error {
	err := client.Client.OpenMetaQuery(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.OpenMetaQuery(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetMetaQueryStatus(bucketName string, options ...osssdk.Option) (osssdk.GetMetaQueryStatusResult, error) {
	res, err := client.Client.GetMetaQueryStatus(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetMetaQueryStatus(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DoMetaQuery(bucketName string, metaQuery osssdk.MetaQuery, options ...osssdk.Option) (osssdk.DoMetaQueryResult, error) {
	res, err := client.Client.DoMetaQuery(bucketName, metaQuery, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.DoMetaQuery(bucketName, metaQuery, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DoMetaQueryXml(bucketName string, metaQueryXml string, options ...osssdk.Option) (osssdk.DoMetaQueryResult, error) {
	res, err := client.Client.DoMetaQueryXml(bucketName, metaQueryXml, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.DoMetaQueryXml(bucketName, metaQueryXml, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) CloseMetaQuery(bucketName string, options ...osssdk.Option) error {
	err := client.Client.CloseMetaQuery(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.CloseMetaQuery(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketWebsite(bucketName string, options ...osssdk.Option) (osssdk.GetBucketWebsiteResult, error) {
	res, err := client.Client.GetBucketWebsite(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketWebsite(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketWebsiteXml(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketWebsiteXml(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketWebsiteXml(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketCORS(bucketName string, corsRules []osssdk.CORSRule, options ...osssdk.Option) error {
	err := client.Client.SetBucketCORS(bucketName, corsRules, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketCORS(bucketName, corsRules, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketCORSXml(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.SetBucketCORSXml(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketCORSXml(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketCORS(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketCORS(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketCORS(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketCORS(bucketName string, options ...osssdk.Option) (osssdk.GetBucketCORSResult, error) {
	res, err := client.Client.GetBucketCORS(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketCORS(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketCORSXml(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketCORSXml(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketCORSXml(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketInfo(bucketName string, options ...osssdk.Option) (osssdk.GetBucketInfoResult, error) {
	res, err := client.Client.GetBucketInfo(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketInfo(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketVersioning(bucketName string, versioningConfig osssdk.VersioningConfig, options ...osssdk.Option) error {
	err := client.Client.SetBucketVersioning(bucketName, versioningConfig, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketVersioning(bucketName, versioningConfig, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketVersioning(bucketName string, options ...osssdk.Option) (osssdk.GetBucketVersioningResult, error) {
	res, err := client.Client.GetBucketVersioning(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketVersioning(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketEncryption(bucketName string, encryptionRule osssdk.ServerEncryptionRule, options ...osssdk.Option) error {
	err := client.Client.SetBucketEncryption(bucketName, encryptionRule, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketEncryption(bucketName, encryptionRule, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketEncryption(bucketName string, options ...osssdk.Option) (osssdk.GetBucketEncryptionResult, error) {
	res, err := client.Client.GetBucketEncryption(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketEncryption(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketEncryption(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketEncryption(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketEncryption(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketTagging(bucketName string, tagging osssdk.Tagging, options ...osssdk.Option) error {
	err := client.Client.SetBucketTagging(bucketName, tagging, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketTagging(bucketName, tagging, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketTagging(bucketName string, options ...osssdk.Option) (osssdk.GetBucketTaggingResult, error) {
	res, err := client.Client.GetBucketTagging(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketTagging(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketTagging(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketTagging(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketTagging(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketStat(bucketName string, options ...osssdk.Option) (osssdk.GetBucketStatResult, error) {
	res, err := client.Client.GetBucketStat(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketStat(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketPolicy(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketPolicy(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketPolicy(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketPolicy(bucketName string, policy string, options ...osssdk.Option) error {
	err := client.Client.SetBucketPolicy(bucketName, policy, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketPolicy(bucketName, policy, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketPolicy(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketPolicy(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketPolicy(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketRequestPayment(bucketName string, paymentConfig osssdk.RequestPaymentConfiguration, options ...osssdk.Option) error {
	err := client.Client.SetBucketRequestPayment(bucketName, paymentConfig, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketRequestPayment(bucketName, paymentConfig, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketRequestPayment(bucketName string, options ...osssdk.Option) (osssdk.RequestPaymentConfiguration, error) {
	res, err := client.Client.GetBucketRequestPayment(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketRequestPayment(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetUserQoSInfo(options ...osssdk.Option) (osssdk.UserQoSConfiguration, error) {
	res, err := client.Client.GetUserQoSInfo(options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetUserQoSInfo(options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketQoSInfo(bucketName string, qosConf osssdk.BucketQoSConfiguration, options ...osssdk.Option) error {
	err := client.Client.SetBucketQoSInfo(bucketName, qosConf, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketQoSInfo(bucketName, qosConf, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketQosInfo(bucketName string, options ...osssdk.Option) (osssdk.BucketQoSConfiguration, error) {
	res, err := client.Client.GetBucketQosInfo(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketQosInfo(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketQosInfo(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketQosInfo(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketQosInfo(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketInventory(bucketName string, inventoryConfig osssdk.InventoryConfiguration, options ...osssdk.Option) error {
	err := client.Client.SetBucketInventory(bucketName, inventoryConfig, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketInventory(bucketName, inventoryConfig, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketInventoryXml(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.SetBucketInventoryXml(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketInventoryXml(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketInventory(bucketName string, strInventoryId string, options ...osssdk.Option) (osssdk.InventoryConfiguration, error) {
	res, err := client.Client.GetBucketInventory(bucketName, strInventoryId, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketInventory(bucketName, strInventoryId, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketInventoryXml(bucketName string, strInventoryId string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketInventoryXml(bucketName, strInventoryId, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketInventoryXml(bucketName, strInventoryId, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) ListBucketInventory(bucketName, continuationToken string, options ...osssdk.Option) (osssdk.ListInventoryConfigurationsResult, error) {
	res, err := client.Client.ListBucketInventory(bucketName, continuationToken, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.ListBucketInventory(bucketName, continuationToken, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) ListBucketInventoryXml(bucketName, continuationToken string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.ListBucketInventoryXml(bucketName, continuationToken, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.ListBucketInventoryXml(bucketName, continuationToken, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketInventory(bucketName, strInventoryId string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketInventory(bucketName, strInventoryId, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketInventory(bucketName, strInventoryId, options...)
		}
	}
	return err
}

func (client ProxyOssClient) SetBucketAsyncTask(bucketName string, asynConf osssdk.AsyncFetchTaskConfiguration, options ...osssdk.Option) (osssdk.AsyncFetchTaskResult, error) {
	res, err := client.Client.SetBucketAsyncTask(bucketName, asynConf, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.SetBucketAsyncTask(bucketName, asynConf, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketAsyncTask(bucketName string, taskID string, options ...osssdk.Option) (osssdk.AsynFetchTaskInfo, error) {
	res, err := client.Client.GetBucketAsyncTask(bucketName, taskID, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketAsyncTask(bucketName, taskID, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) InitiateBucketWorm(bucketName string, retentionDays int, options ...osssdk.Option) (string, error) {
	res, err := client.Client.InitiateBucketWorm(bucketName, retentionDays, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.InitiateBucketWorm(bucketName, retentionDays, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) AbortBucketWorm(bucketName string, options ...osssdk.Option) error {
	err := client.Client.AbortBucketWorm(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.AbortBucketWorm(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) CompleteBucketWorm(bucketName string, wormID string, options ...osssdk.Option) error {
	err := client.Client.CompleteBucketWorm(bucketName, wormID, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.CompleteBucketWorm(bucketName, wormID, options...)
		}
	}
	return err
}

func (client ProxyOssClient) ExtendBucketWorm(bucketName string, retentionDays int, wormID string, options ...osssdk.Option) error {
	err := client.Client.ExtendBucketWorm(bucketName, retentionDays, wormID, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.ExtendBucketWorm(bucketName, retentionDays, wormID, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketWorm(bucketName string, options ...osssdk.Option) (osssdk.WormConfiguration, error) {
	res, err := client.Client.GetBucketWorm(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketWorm(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) SetBucketTransferAcc(bucketName string, accConf osssdk.TransferAccConfiguration, options ...osssdk.Option) error {
	err := client.Client.SetBucketTransferAcc(bucketName, accConf, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.SetBucketTransferAcc(bucketName, accConf, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketTransferAcc(bucketName string, options ...osssdk.Option) (osssdk.TransferAccConfiguration, error) {
	res, err := client.Client.GetBucketTransferAcc(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketTransferAcc(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketTransferAcc(bucketName string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketTransferAcc(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketTransferAcc(bucketName, options...)
		}
	}
	return err
}

func (client ProxyOssClient) PutBucketReplication(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.PutBucketReplication(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.PutBucketReplication(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketReplication(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketReplication(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketReplication(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) DeleteBucketReplication(bucketName string, ruleId string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketReplication(bucketName, ruleId, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketReplication(bucketName, ruleId, options...)
		}
	}
	return err
}

func (client ProxyOssClient) GetBucketReplicationLocation(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketReplicationLocation(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketReplicationLocation(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketReplicationProgress(bucketName string, ruleId string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketReplicationProgress(bucketName, ruleId, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketReplicationProgress(bucketName, ruleId, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketCname(bucketName string, options ...osssdk.Option) (string, error) {
	res, err := client.Client.GetBucketCname(bucketName, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketCname(bucketName, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) CreateBucketCnameToken(bucketName string, cname string, options ...osssdk.Option) (osssdk.CreateBucketCnameTokenResult, error) {
	res, err := client.Client.CreateBucketCnameToken(bucketName, cname, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.CreateBucketCnameToken(bucketName, cname, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) GetBucketCnameToken(bucketName string, cname string, options ...osssdk.Option) (osssdk.GetBucketCnameTokenResult, error) {
	res, err := client.Client.GetBucketCnameToken(bucketName, cname, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return res, err
			}
			return client.Client.GetBucketCnameToken(bucketName, cname, options...)
		}
	}
	return res, err
}

func (client ProxyOssClient) PutBucketCnameXml(bucketName string, xmlBody string, options ...osssdk.Option) error {
	err := client.Client.PutBucketCnameXml(bucketName, xmlBody, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.PutBucketCnameXml(bucketName, xmlBody, options...)
		}
	}
	return err
}

func (client ProxyOssClient) PutBucketCname(bucketName string, cname string, options ...osssdk.Option) error {
	err := client.Client.PutBucketCname(bucketName, cname, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.PutBucketCname(bucketName, cname, options...)
		}
	}
	return err
}

func (client ProxyOssClient) DeleteBucketCname(bucketName string, cname string, options ...osssdk.Option) error {
	err := client.Client.DeleteBucketCname(bucketName, cname, options...)
	if err != nil {
		if client.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(client.secretName)
			if err != nil {
				return err
			}
			return client.Client.DeleteBucketCname(bucketName, cname, options...)
		}
	}
	return err
}
