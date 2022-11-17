package sdk

import (
	commonsdk "github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
	osssdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/http"
)

type ProxyBucket struct {
	*osssdk.Bucket

	secretName      string
	akExpireHandler service.AKExpireHandler
}

func (bucket ProxyBucket) judgeTempAKExpire(err error) bool {
	return bucket.akExpireHandler.JudgeAKExpire(err)
}

func (bucket ProxyBucket) PutObject(objectKey string, reader io.Reader, options ...osssdk.Option) error {
	err := bucket.Bucket.PutObject(objectKey, reader, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutObject(objectKey, reader, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) PutObjectFromFile(objectKey, filePath string, options ...osssdk.Option) error {
	err := bucket.Bucket.PutObjectFromFile(objectKey, filePath, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutObjectFromFile(objectKey, filePath, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) DoPutObject(request *osssdk.PutObjectRequest, options []osssdk.Option) (*osssdk.Response, error) {
	res, err := bucket.Bucket.DoPutObject(request, options)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DoPutObject(request, options)
		}
	}
	return res, err
}

func (bucket ProxyBucket) GetObject(objectKey string, options ...osssdk.Option) (io.ReadCloser, error) {
	res, err := bucket.Bucket.GetObject(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObject(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) GetObjectToFile(objectKey, filePath string, options ...osssdk.Option) error {
	err := bucket.Bucket.GetObjectToFile(objectKey, filePath, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.GetObjectToFile(objectKey, filePath, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) DoGetObject(request *osssdk.GetObjectRequest, options []osssdk.Option) (*osssdk.GetObjectResult, error) {
	res, err := bucket.Bucket.DoGetObject(request, options)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DoGetObject(request, options)
		}
	}
	return res, err
}

func (bucket ProxyBucket) CopyObject(srcObjectKey, destObjectKey string, options ...osssdk.Option) (osssdk.CopyObjectResult, error) {
	res, err := bucket.Bucket.CopyObject(srcObjectKey, destObjectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.CopyObject(srcObjectKey, destObjectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) CopyObjectTo(destBucketName, destObjectKey, srcObjectKey string, options ...osssdk.Option) (osssdk.CopyObjectResult, error) {
	res, err := bucket.Bucket.CopyObjectTo(destBucketName, destObjectKey, srcObjectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.CopyObjectTo(destBucketName, destObjectKey, srcObjectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) CopyObjectFrom(srcBucketName, srcObjectKey, destObjectKey string, options ...osssdk.Option) (osssdk.CopyObjectResult, error) {
	res, err := bucket.Bucket.CopyObjectFrom(srcBucketName, srcObjectKey, destObjectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.CopyObjectFrom(srcBucketName, srcObjectKey, destObjectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) AppendObject(objectKey string, reader io.Reader, appendPosition int64, options ...osssdk.Option) (int64, error) {
	res, err := bucket.Bucket.AppendObject(objectKey, reader, appendPosition, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.AppendObject(objectKey, reader, appendPosition, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) DoAppendObject(request *osssdk.AppendObjectRequest, options []osssdk.Option) (*osssdk.AppendObjectResult, error) {
	res, err := bucket.Bucket.DoAppendObject(request, options)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DoAppendObject(request, options)
		}
	}
	return res, err
}

func (bucket ProxyBucket) DeleteObject(objectKey string, options ...osssdk.Option) error {
	err := bucket.Bucket.DeleteObject(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.DeleteObject(objectKey, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) DeleteObjects(objectKeys []string, options ...osssdk.Option) (osssdk.DeleteObjectsResult, error) {
	res, err := bucket.Bucket.DeleteObjects(objectKeys, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DeleteObjects(objectKeys, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) DeleteObjectVersions(objectVersions []osssdk.DeleteObject, options ...osssdk.Option) (osssdk.DeleteObjectVersionsResult, error) {
	res, err := bucket.Bucket.DeleteObjectVersions(objectVersions, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DeleteObjectVersions(objectVersions, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) IsObjectExist(objectKey string, options ...osssdk.Option) (bool, error) {
	res, err := bucket.Bucket.IsObjectExist(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.IsObjectExist(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) ListObjects(options ...osssdk.Option) (osssdk.ListObjectsResult, error) {
	res, err := bucket.Bucket.ListObjects(options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.ListObjects(options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) ListObjectsV2(options ...osssdk.Option) (osssdk.ListObjectsResultV2, error) {
	res, err := bucket.Bucket.ListObjectsV2(options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.ListObjectsV2(options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) ListObjectVersions(options ...osssdk.Option) (osssdk.ListObjectVersionsResult, error) {
	res, err := bucket.Bucket.ListObjectVersions(options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.ListObjectVersions(options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) SetObjectMeta(objectKey string, options ...osssdk.Option) error {
	err := bucket.Bucket.SetObjectMeta(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.SetObjectMeta(objectKey, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) GetObjectDetailedMeta(objectKey string, options ...osssdk.Option) (http.Header, error) {
	res, err := bucket.Bucket.GetObjectDetailedMeta(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObjectDetailedMeta(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) GetObjectMeta(objectKey string, options ...osssdk.Option) (http.Header, error) {
	res, err := bucket.Bucket.GetObjectMeta(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObjectMeta(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) SetObjectACL(objectKey string, objectACL osssdk.ACLType, options ...osssdk.Option) error {
	err := bucket.Bucket.SetObjectACL(objectKey, objectACL, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.SetObjectACL(objectKey, objectACL, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) GetObjectACL(objectKey string, options ...osssdk.Option) (osssdk.GetObjectACLResult, error) {
	res, err := bucket.Bucket.GetObjectACL(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObjectACL(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) PutSymlink(symObjectKey string, targetObjectKey string, options ...osssdk.Option) error {
	err := bucket.Bucket.PutSymlink(symObjectKey, targetObjectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutSymlink(symObjectKey, targetObjectKey, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) GetSymlink(objectKey string, options ...osssdk.Option) (http.Header, error) {
	res, err := bucket.Bucket.GetSymlink(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetSymlink(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) RestoreObject(objectKey string, options ...osssdk.Option) error {
	err := bucket.Bucket.RestoreObject(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.RestoreObject(objectKey, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) RestoreObjectDetail(objectKey string, restoreConfig osssdk.RestoreConfiguration, options ...osssdk.Option) error {
	err := bucket.Bucket.RestoreObjectDetail(objectKey, restoreConfig, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.RestoreObjectDetail(objectKey, restoreConfig, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) RestoreObjectXML(objectKey, configXML string, options ...osssdk.Option) error {
	err := bucket.Bucket.RestoreObjectXML(objectKey, configXML, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.RestoreObjectXML(objectKey, configXML, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) PutObjectWithURL(signedURL string, reader io.Reader, options ...osssdk.Option) error {
	err := bucket.Bucket.PutObjectWithURL(signedURL, reader, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutObjectWithURL(signedURL, reader, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) PutObjectFromFileWithURL(signedURL, filePath string, options ...osssdk.Option) error {
	err := bucket.Bucket.PutObjectFromFileWithURL(signedURL, filePath, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutObjectFromFileWithURL(signedURL, filePath, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) DoPutObjectWithURL(signedURL string, reader io.Reader, options []osssdk.Option) (*osssdk.Response, error) {
	res, err := bucket.Bucket.DoPutObjectWithURL(signedURL, reader, options)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DoPutObjectWithURL(signedURL, reader, options)
		}
	}
	return res, err
}

func (bucket ProxyBucket) GetObjectWithURL(signedURL string, options ...osssdk.Option) (io.ReadCloser, error) {
	res, err := bucket.Bucket.GetObjectWithURL(signedURL, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObjectWithURL(signedURL, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) GetObjectToFileWithURL(signedURL, filePath string, options ...osssdk.Option) error {
	err := bucket.Bucket.GetObjectToFileWithURL(signedURL, filePath, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.GetObjectToFileWithURL(signedURL, filePath, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) DoGetObjectWithURL(signedURL string, options []osssdk.Option) (*osssdk.GetObjectResult, error) {
	res, err := bucket.Bucket.DoGetObjectWithURL(signedURL, options)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.DoGetObjectWithURL(signedURL, options)
		}
	}
	return res, err
}

func (bucket ProxyBucket) ProcessObject(objectKey string, process string, options ...osssdk.Option) (osssdk.ProcessObjectResult, error) {
	res, err := bucket.Bucket.ProcessObject(objectKey, process, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.ProcessObject(objectKey, process, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) PutObjectTagging(objectKey string, tagging osssdk.Tagging, options ...osssdk.Option) error {
	err := bucket.Bucket.PutObjectTagging(objectKey, tagging, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.PutObjectTagging(objectKey, tagging, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) GetObjectTagging(objectKey string, options ...osssdk.Option) (osssdk.GetObjectTaggingResult, error) {
	res, err := bucket.Bucket.GetObjectTagging(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.GetObjectTagging(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) DeleteObjectTagging(objectKey string, options ...osssdk.Option) error {
	err := bucket.Bucket.DeleteObjectTagging(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return err
			}
			return bucket.Bucket.DeleteObjectTagging(objectKey, options...)
		}
	}
	return err
}

func (bucket ProxyBucket) OptionsMethod(objectKey string, options ...osssdk.Option) (http.Header, error) {
	res, err := bucket.Bucket.OptionsMethod(objectKey, options...)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.OptionsMethod(objectKey, options...)
		}
	}
	return res, err
}

func (bucket ProxyBucket) Do(method, objectName string, params map[string]interface{}, options []osssdk.Option,
	data io.Reader, listener osssdk.ProgressListener) (*osssdk.Response, error) {
	res, err := bucket.Bucket.Do(method, objectName, params, options, data, listener)
	if err != nil {
		if bucket.judgeTempAKExpire(err) {
			err = commonsdk.RefreshSecretInfo(bucket.secretName)
			if err != nil {
				return res, err
			}
			return bucket.Bucket.Do(method, objectName, params, options, data, listener)
		}
	}
	return res, err
}
